package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_basic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_router_v2.router_1", "name", "router_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_update_external_gw(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_update_external_gw_1,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update_external_gw_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_router_v2.router_1", "external_gateway", env.OS_EXTGW_ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_timeout(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_timeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists("opentelekomcloud_networking_router_v2.router_1", &router),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_router_v2" {
			continue
		}

		_, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router still exists")
		}
	}

	return nil
}

const testAccNetworkingV2Router_basic = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
	name = "router_1"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_update = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
	name = "router_2"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_update_external_gw_1 = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
	distributed = "false"
}
`

var testAccNetworkingV2Router_update_external_gw_2 = fmt.Sprintf(`
resource "opentelekomcloud_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
	distributed = "false"
	external_gateway = "%s"
}
`, env.OS_EXTGW_ID)

const testAccNetworkingV2Router_timeout = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
	name = "router_1"
	admin_state_up = "true"
	distributed = "false"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
