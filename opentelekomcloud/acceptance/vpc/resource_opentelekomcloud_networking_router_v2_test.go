package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2RouterBasic(t *testing.T) {
	var router routers.Router
	resourceName := "opentelekomcloud_networking_router_v2.router_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists(resourceName, &router),
				),
			},
			{
				Config: testAccNetworkingV2RouterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "router_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterUpdateExternalGw(t *testing.T) {
	var router routers.Router
	resourceName := "opentelekomcloud_networking_router_v2.router_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists(resourceName, &router),
				),
			},
			{
				Config: testAccNetworkingV2RouterUpdateExternalGw,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "external_gateway", env.OS_EXTGW_ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_timeout(t *testing.T) {
	var router routers.Router
	resourceName := "opentelekomcloud_networking_router_v2.router_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterTimeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists(resourceName, &router),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_router_v2" {
			continue
		}

		_, err := routers.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("router still exists")
		}
	}

	return nil
}

const testAccNetworkingV2RouterBasic = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = true
  distributed    = false
}
`

const testAccNetworkingV2RouterUpdate = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_2"
  admin_state_up = true
  distributed    = false
}
`

var testAccNetworkingV2RouterUpdateExternalGw = fmt.Sprintf(`
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "router_1"
  admin_state_up   = true
  distributed      = false
  external_gateway = "%s"
}
`, env.OS_EXTGW_ID)

const testAccNetworkingV2RouterTimeout = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = true
  distributed    = false

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
