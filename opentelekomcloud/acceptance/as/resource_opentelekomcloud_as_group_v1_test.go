package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/groups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASV1Group_basic(t *testing.T) {
	var asGroup groups.Group
	resourceName := "opentelekomcloud_as_group_v1.as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.LoadBalancer, Count: 1},
				{Q: quotas.LbListener, Count: 1},
				{Q: quotas.LbPool, Count: 1},
				{Q: quotas.ASGroup, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1GroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "700"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-create"),
					resource.TestCheckResourceAttr(resourceName, "delete_publicip", "false"),
					resource.TestCheckResourceAttr(resourceName, "delete_instances", "no"),
				),
			},
			{
				Config: testAccASV1GroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "500"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-update"),
					resource.TestCheckResourceAttr(resourceName, "delete_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_instances", "yes"),
				),
			},
		},
	})
}

func TestAccASV1Group_RemoveWithSetMinNumber(t *testing.T) {
	var asGroup groups.Group
	resourceName := "opentelekomcloud_as_group_v1.as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.ASGroup, Count: 1},
				{Q: quotas.ASConfiguration, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1GroupRemoveWithSetMinNumber,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "delete_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "scaling_group_name", "as_group"),
				),
			},
		},
	})
}

func TestAccASV1Group_WithoutSecurityGroups(t *testing.T) {
	var asGroup groups.Group

	resourceName := "opentelekomcloud_as_group_v1.as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.ASGroup, Count: 1},
				{Q: quotas.ASConfiguration, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1GroupWithoutSGs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "0"),
				),
			},
		},
	})
}

func testAccCheckASV1GroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_group_v1" {
			continue
		}

		_, err := groups.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("AS group still exists")
		}
	}

	return nil
}

func testAccCheckASV1GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
		}

		found, err := groups.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("AS group not found")
		}
		group = found

		return nil
	}
}

var testAccASV1GroupBasic = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s


resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
}

resource "opentelekomcloud_as_configuration_v1" "as_config" {
  scaling_configuration_name = "as_config"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name = "%s"
  }
}

resource "opentelekomcloud_as_group_v1" "as_group" {
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  delete_publicip          = false
  delete_instances         = "no"

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  lbaas_listeners {
    pool_id       = opentelekomcloud_lb_pool_v2.pool_1.id
    protocol_port = opentelekomcloud_lb_listener_v2.listener_1.protocol_port
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  health_periodic_audit_grace_period = 700

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)

var testAccASV1GroupUpdate = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
}

resource "opentelekomcloud_as_configuration_v1" "as_config" {
  scaling_configuration_name = "as_config"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name = "%s"
  }
}

resource "opentelekomcloud_as_group_v1" "as_group" {
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  lbaas_listeners {
    pool_id       = opentelekomcloud_lb_pool_v2.pool_1.id
    protocol_port = opentelekomcloud_lb_listener_v2.listener_1.protocol_port
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  health_periodic_audit_grace_period = 500

  tags = {
    muh = "value-update"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)

var testAccASV1GroupRemoveWithSetMinNumber = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_as_configuration_v1" "as_config" {
  scaling_configuration_name = "as_config"
  instance_config {
    image    = data.opentelekomcloud_images_image_v2.latest_image.id
    key_name = "%s"
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }

    metadata = {
      environment  = "otc-test"
      generator    = "terraform"
      puppetmaster = "pseudo-puppet"
      role         = "pseudo-role"
      autoscaling  = "proxy_ASG"
    }
  }
}

resource "opentelekomcloud_as_group_v1" "as_group" {
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  available_zones          = ["%s"]
  desire_instance_number   = 3
  min_instance_number      = 1
  max_instance_number      = 10
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }

  lifecycle {
    ignore_changes = [
      instances
    ]
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME, env.OS_AVAILABILITY_ZONE)

var testAccASV1GroupWithoutSGs = fmt.Sprintf(`
// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_as_configuration_v1" "as_config" {
  scaling_configuration_name = "as_config"
  instance_config {
    image    = data.opentelekomcloud_images_image_v2.latest_image.id
    key_name = "%s"
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }

    metadata = {
      environment  = "otc-test"
      generator    = "terraform"
      puppetmaster = "pseudo-puppet"
      role         = "pseudo-role"
      autoscaling  = "proxy_ASG"
    }
  }
}

resource "opentelekomcloud_as_group_v1" "as_group" {
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  available_zones          = ["%s"]
  desire_instance_number   = 0
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME, env.OS_AVAILABILITY_ZONE)
