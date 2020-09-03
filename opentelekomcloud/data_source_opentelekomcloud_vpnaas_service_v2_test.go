package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccVpnServiceV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpnServiceV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_id", "name", "vpn_service_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_name", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_id", "description", ""),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v2.by_name", "subnet_id", ""),
				),
			},
		},
	})
}

var testAccDataSourceVpnServiceV2Config = fmt.Sprintf(`
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

data "opentelekomcloud_vpnaas_service_v2" "by_name" {
  name = opentelekomcloud_vpnaas_service_v2.service_1.name
  admin_state_up = "true"
}
`)
