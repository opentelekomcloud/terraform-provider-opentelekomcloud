package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASV1Policy_basic(t *testing.T) {
	supportedRegions := []string{"eu-de", "eu-nl", "eu-ch2"}
	var asPolicy policies.Policy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckServiceAvailability(t, testServiceV1, supportedRegions)
			qts := quotas.MultipleQuotas{
				{Q: quotas.ASGroup, Count: 1},
				{Q: quotas.ASConfiguration, Count: 1},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1PolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1PolicyExists("opentelekomcloud_as_policy_v1.as_policy", &asPolicy),
				),
			},
		},
	})
}

func testAccCheckASV1PolicyDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_policy_v1" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("AS policy still exists")
		}
	}

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
		client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
		}

		found, err := policies.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		policy = found

		return nil
	}
}

var testAccASV1PolicyBasic = fmt.Sprintf(`
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
  delete_instances         = "yes"
  delete_publicip          = true
  networks {
    id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
  security_groups {
    id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  }
  vpc_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
}

resource "opentelekomcloud_as_policy_v1" "as_policy"{
  scaling_policy_name = "as_policy"
  scaling_group_id    = opentelekomcloud_as_group_v1.as_group.id
  scaling_policy_type = "SCHEDULED"
  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time = "2022-12-22T12:00Z"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceImage, common.DataSourceSubnet, env.OS_KEYPAIR_NAME)
