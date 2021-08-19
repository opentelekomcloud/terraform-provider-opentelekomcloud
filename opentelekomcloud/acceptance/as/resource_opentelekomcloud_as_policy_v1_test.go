package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASV1Policy_basic(t *testing.T) {
	var asPolicy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccFlavorPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Policy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1PolicyExists("opentelekomcloud_as_policy_v1.hth_as_policy", &asPolicy),
				),
			},
		},
	})
}

func testAccCheckASV1PolicyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	asClient, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_policy_v1" {
			continue
		}

		_, err := policies.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS policy still exists")
		}
	}

	log.Printf("[DEBUG] testCheckASV1PolicyDestroy success!")

	return nil
}

func testAccCheckASV1PolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
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

		found, err := policies.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] test found is: %#v", found)
		policy = &found

		return nil
	}
}

var testASV1Policy_basic = fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name        = "terraform"
  description = "This is a terraform test security group"
}

resource "opentelekomcloud_compute_keypair_v2" "hth_key" {
  name = "hth_key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "opentelekomcloud_as_configuration_v1" "hth_as_config"{
  scaling_configuration_name = "hth_as_config"
  instance_config {
    image = "%s"
    disk {
      size = 40
      volume_type = "SATA"
      disk_type = "SYS"
    }
    key_name = opentelekomcloud_compute_keypair_v2.hth_key.id
  }
}

resource "opentelekomcloud_as_group_v1" "hth_as_group"{
  scaling_group_name = "hth_as_group"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.hth_as_config.id
  networks {
    id = "%s"
  }
  security_groups {
    id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  }
  vpc_id = "%s"
}

resource "opentelekomcloud_as_policy_v1" "hth_as_policy"{
  scaling_policy_name = "terraform"
  scaling_group_id = opentelekomcloud_as_group_v1.hth_as_group.id
  scaling_policy_type = "SCHEDULED"
  scaling_policy_action {
    operation = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time = "2022-12-22T12:00Z"
  }
}
`, env.OS_IMAGE_ID, env.OS_NETWORK_ID, env.OS_VPC_ID)
