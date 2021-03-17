package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASV1Group_basic(t *testing.T) {
	var asGroup groups.Group

	resourceName := "opentelekomcloud_as_group_v1.hth_as_group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccFlavorPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Group_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", env.OS_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "700"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testASV1Group_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "500"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccASV1Group_RemoveWithSetMinNumber(t *testing.T) {
	var asGroup groups.Group

	resourceName := "opentelekomcloud_as_group_v1.proxy_group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccFlavorPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Group_removeWithSetMinNumber,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "delete_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "scaling_group_name", "proxy-test-asg"),
				),
			},
		},
	})
}

func testAccCheckASV1GroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	asClient, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_group_v1" {
			continue
		}

		_, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS group still exists")
		}
	}

	log.Printf("[DEBUG] testCheckASV1GroupDestroy success!")

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
		asClient, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud autoscaling client: %s", err)
		}

		found, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("autoscaling Group not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		group = &found

		return nil
	}
}

var testASV1Group_basic = fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = "test-acc"
}

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
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

resource "opentelekomcloud_as_configuration_v1" "hth_as_config"{
  scaling_configuration_name = "hth_as_config"
  instance_config {
    image = "%s"
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name = "%s"
  }
}

resource "opentelekomcloud_as_group_v1" "hth_as_group"{
  scaling_group_name       = "hth_as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.hth_as_config.id
  networks {
    id = "%s"
  }
  security_groups {
    id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  }
  lbaas_listeners {
    pool_id =       opentelekomcloud_lb_pool_v2.pool_1.id
    protocol_port = opentelekomcloud_lb_listener_v2.listener_1.protocol_port
  }
  vpc_id = "%s"

  health_periodic_audit_grace_period = 700

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, env.OS_SUBNET_ID, env.OS_IMAGE_ID, env.OS_KEYPAIR_NAME, env.OS_NETWORK_ID, env.OS_VPC_ID)

var testASV1Group_update = fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = "test-acc"
}

resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
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

resource "opentelekomcloud_as_configuration_v1" "hth_as_config"{
  scaling_configuration_name = "hth_as_config"
  instance_config {
    image = "%s"
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name = "%s"
  }
}

resource "opentelekomcloud_as_group_v1" "hth_as_group"{
  scaling_group_name       = "hth_as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.hth_as_config.id
  networks {
    id = "%s"
  }
  security_groups {
    id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  }
  lbaas_listeners {
    pool_id =       opentelekomcloud_lb_pool_v2.pool_1.id
    protocol_port = opentelekomcloud_lb_listener_v2.listener_1.protocol_port
  }
  vpc_id = "%s"

  health_periodic_audit_grace_period = 500

  tags = {
    muh = "value-update"
  }
}
`, env.OS_SUBNET_ID, env.OS_IMAGE_ID, env.OS_KEYPAIR_NAME, env.OS_NETWORK_ID, env.OS_VPC_ID)

var testASV1Group_removeWithSetMinNumber = fmt.Sprintf(`
resource "opentelekomcloud_compute_secgroup_v2" "secgroup" {
  name        = "acc-test-sg"
  description = "Security group for AS tf test"
}

# Proxy AS configuration
resource "opentelekomcloud_as_configuration_v1" "proxy_config" {
  scaling_configuration_name = "proxy-test-asg"
  instance_config {
    image     = "%s"
    key_name  = "%s"
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }

    metadata  = {
      environment  = "otc-test"
      generator    = "terraform"
      puppetmaster = "pseudo-puppet"
      role         = "pseudo-role"
      autoscaling  = "proxy_ASG"
    }
  }
}

resource "opentelekomcloud_as_group_v1" "proxy_group" {
  scaling_group_name       = "proxy-test-asg"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.proxy_config.id
  available_zones          = ["%s"]
  desire_instance_number   = 3
  min_instance_number      = 1
  max_instance_number      = 10
  vpc_id                   = "%s"
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = "%s"
  }
  security_groups {
    id = opentelekomcloud_compute_secgroup_v2.secgroup.id
  }

  lifecycle {
    ignore_changes = [
      instances
    ]
  }
}
`, env.OS_IMAGE_ID, env.OS_KEYPAIR_NAME, env.OS_AVAILABILITY_ZONE, env.OS_VPC_ID, env.OS_NETWORK_ID)
