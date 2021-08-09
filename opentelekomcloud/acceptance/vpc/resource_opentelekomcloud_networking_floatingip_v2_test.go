package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2FloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2FloatingIPExists("opentelekomcloud_networking_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func TestAccNetworkingV2FloatingIP_timeout(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPTimeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2FloatingIPExists("opentelekomcloud_networking_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2FloatingIPDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud floating IP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_floatingip_v2" {
			continue
		}

		_, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("floatingIP still exists")
		}
	}

	return nil
}

const testAccNetworkingV2FloatingIPBasic = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}
`

const testAccNetworkingV2FloatingIPTimeout = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
