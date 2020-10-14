package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVpnServiceV2DataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpnServiceV2ConfigById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_id", "name", "vpn_service_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_id", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccVpnServiceV2DataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpnServiceV2ConfigByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_name", "name", "vpn_service_2"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_name", "admin_state_up", "false"),
				),
			},
		},
	})
}

var testAccDataSourceVpnServiceV2ConfigById = fmt.Sprintf(`
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_vpnaas_service_v2" "service_1" {
  name           = "vpn_service_1"
  router_id      = opentelekomcloud_networking_router_v2.router_1.id
  admin_state_up = "true"
}

data "opentelekomcloud_vpnaas_service_v2" "by_id" {
  router_id = opentelekomcloud_vpnaas_service_v2.service_1.router_id
  admin_state_up = "true"
}
`)

var testAccDataSourceVpnServiceV2ConfigByName = fmt.Sprintf(`
resource "opentelekomcloud_networking_router_v2" "router_2" {
  name           = "router_2"
  admin_state_up = "true"
}

resource "opentelekomcloud_vpnaas_service_v2" "service_2" {
  name           = "vpn_service_2"
  router_id      = opentelekomcloud_networking_router_v2.router_2.id
  admin_state_up = "false"
}

data "opentelekomcloud_vpnaas_service_v2" "by_name" {
  name = opentelekomcloud_vpnaas_service_v2.service_2.name
}
`)
