package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v2/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASPolicyV2_basic(t *testing.T) {
	var asPolicy policies.Policy
	resourceName := "opentelekomcloud_as_policy_v2.as_policy"

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
		CheckDestroy:      testAccCheckASV2PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV2PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.operation", "ADD"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.percentage", "15"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", "as_policy"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "300"),
				),
			},
			{
				Config: testAccASPolicyV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV2PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.percentage", "30"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", "as_policy_update"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "100"),
				),
			},
		},
	})
}

func TestAccASPolicyV2_withSize(t *testing.T) {
	var asPolicy policies.Policy
	resourceName := "opentelekomcloud_as_policy_v2.as_policy"

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
		CheckDestroy:      testAccCheckASV2PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASPolicyV2ActionWithSize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV2PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.operation", "ADD"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.percentage", "0"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", "as_policy"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "100"),
				),
			},
		},
	})
}

func TestAccASPolicyV2_conflictsAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccASPolicyV2ActionConflict,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("select one from `percentage` or `size`+"),
			},
		},
	})
}

func testAccCheckASV2PolicyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.AutoscalingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_policy_v2" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS policyV2 still exists")
		}
	}

	return nil
}

func testAccCheckASV2PolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		asClient, err := config.AutoscalingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV2 client: %w", err)
		}

		found, err := policies.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		policy = &found

		return nil
	}
}

var testAccASPolicyV2Basic = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_as_configuration_v1" "as_config"{
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

resource "opentelekomcloud_as_group_v1" "as_group"{
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  delete_publicip  = true
  delete_instances = "yes"
}

resource "opentelekomcloud_as_policy_v2" "as_policy"{
  scaling_policy_name   = "as_policy"
  scaling_policy_type   = "RECURRENCE"
  scaling_resource_id   = opentelekomcloud_as_group_v1.as_group.id
  scaling_resource_type = "SCALING_GROUP"

  scaling_policy_action {
    operation  = "ADD"
    percentage = 15
  }
  scheduled_policy {
    launch_time      = "10:30"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    end_time         = "2040-12-31T10:30Z"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)

var testAccASPolicyV2Update = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_as_configuration_v1" "as_config"{
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

resource "opentelekomcloud_as_group_v1" "as_group"{
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  delete_publicip  = true
  delete_instances = "yes"
}

resource "opentelekomcloud_as_policy_v2" "as_policy"{
  scaling_policy_name   = "as_policy_update"
  scaling_policy_type   = "RECURRENCE"
  scaling_resource_id   = opentelekomcloud_as_group_v1.as_group.id
  scaling_resource_type = "SCALING_GROUP"

  scaling_policy_action {
    operation  = "ADD"
    percentage = 30
  }
  scheduled_policy {
    launch_time      = "10:30"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    end_time         = "2040-12-31T10:30Z"
  }
  cool_down_time = 100
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)

var testAccASPolicyV2ActionWithSize = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_as_configuration_v1" "as_config"{
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

resource "opentelekomcloud_as_group_v1" "as_group"{
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  delete_publicip  = true
  delete_instances = "yes"
}

resource "opentelekomcloud_as_policy_v2" "as_policy"{
  scaling_policy_name   = "as_policy"
  scaling_policy_type   = "RECURRENCE"
  scaling_resource_id   = opentelekomcloud_as_group_v1.as_group.id
  scaling_resource_type = "SCALING_GROUP"

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }
  scheduled_policy {
    launch_time      = "10:30"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    end_time         = "2040-12-31T10:30Z"
  }
  cool_down_time = 100
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)

var testAccASPolicyV2ActionConflict = fmt.Sprintf(`
// default SecGroup data-source
%s

// default Image data-source
%s

// default Subnet data-source
%s

resource "opentelekomcloud_as_configuration_v1" "as_config"{
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

resource "opentelekomcloud_as_group_v1" "as_group"{
  scaling_group_name       = "as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id

  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id

  delete_publicip  = true
  delete_instances = "yes"
}

resource "opentelekomcloud_as_policy_v2" "as_policy"{
  scaling_policy_name   = "as_policy"
  scaling_policy_type   = "RECURRENCE"
  scaling_resource_id   = opentelekomcloud_as_group_v1.as_group.id
  scaling_resource_type = "SCALING_GROUP"

  scaling_policy_action {
    operation  = "ADD"
    size       = 1
	percentage = 15
  }
  scheduled_policy {
    launch_time      = "10:30"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    end_time         = "2040-12-31T10:30Z"
  }
  cool_down_time = 100
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)
