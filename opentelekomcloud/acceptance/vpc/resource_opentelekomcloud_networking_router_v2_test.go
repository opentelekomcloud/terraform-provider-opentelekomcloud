package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRouterName = "opentelekomcloud_networking_router_v2.router_1"

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router
	t.Parallel()
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists(resourceRouterName, &router),
				),
			},
			{
				Config: testAccNetworkingV2RouterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRouterName, "name", "router_2_b"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_update_external_gw(t *testing.T) {
	var router routers.Router
	t.Parallel()
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterExternalGw,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists(resourceRouterName, &router),
				),
			},
			{
				Config: testAccNetworkingV2RouterExternalGwUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceRouterName, "external_gateway"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_timeout(t *testing.T) {
	var router routers.Router
	t.Parallel()
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterTimeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists(resourceRouterName, &router),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterSnat(t *testing.T) {
	var router routers.Router
	t.Parallel()
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterEnableSnat,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2RouterExists(resourceRouterName, &router),
					resource.TestCheckResourceAttr(resourceRouterName, "enable_snat", "true"),
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
  name           = "router_1_b"
  admin_state_up = true
  distributed    = false
}
`

const testAccNetworkingV2RouterUpdate = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_2_b"
  admin_state_up = true
  distributed    = false
}
`

const testAccNetworkingV2RouterExternalGw = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_gw"
  admin_state_up = true
  distributed    = false
}
`

var testAccNetworkingV2RouterExternalGwUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "router_1_gw"
  admin_state_up   = true
  distributed      = false
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
}
`, common.DataSourceExtNetwork)

const testAccNetworkingV2RouterTimeout = `
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1_t"
  admin_state_up = true
  distributed    = false

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

var testAccNetworkingV2RouterEnableSnat = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "router_1_snat"
  admin_state_up   = true
  distributed      = false
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
  enable_snat      = true
}
`, common.DataSourceExtNetwork)
