package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v2/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASPolicyV2_basic(t *testing.T) {
	var asPolicy policies.Policy
	resourceName := "opentelekomcloud_as_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccFlavorPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV2PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV2PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.operation", "ADD"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.percentage", "15"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", "policy_create"),
				),
			},
			{
				Config: testASPolicyV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV2PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.percentage", "30"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", "policy_update"),
				),
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

var testASPolicyV2Basic = fmt.Sprintf(`
data "opentelekomcloud_networking_secgroup_v2" "sg_1" {
  name = "default"
}

resource "opentelekomcloud_as_group_v1" "group_1"{
  scaling_group_name = "group_1"

  networks {
    id = "%s"
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.sg_1.id
  }
  vpc_id = "%s"
}

resource "opentelekomcloud_as_policy_v2" "policy_1"{
  scaling_policy_name   = "policy_create"
  scaling_policy_type   = "RECURRENCE"
  scaling_resource_id   = opentelekomcloud_as_group_v1.group_1.id
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
`, env.OS_NETWORK_ID, env.OS_VPC_ID)

var testASPolicyV2Update = fmt.Sprintf(`
data "opentelekomcloud_networking_secgroup_v2" "sg_1" {
  name = "default"
}

resource "opentelekomcloud_as_group_v1" "group_1"{
  scaling_group_name = "group_1"

  networks {
    id = "%s"
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.sg_1.id
  }
  vpc_id = "%s"
}

resource "opentelekomcloud_as_policy_v2" "policy_1"{
  scaling_policy_name   = "policy_update"
  scaling_policy_type   = "RECURRENCE"
  scaling_resource_id   = opentelekomcloud_as_group_v1.group_1.id
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
  cool_down_time = 0
}
`, env.OS_NETWORK_ID, env.OS_VPC_ID)
