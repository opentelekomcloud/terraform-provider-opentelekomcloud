package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
	"testing"
)

func TestAccVpnServiceV2DataSource_basic(t *testing.T) {
	var service services.Service
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpnServiceV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnServiceV2Exists("opentelekomcloud_vpnaas_service_v2.service_1", &service),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpnaas_service_v1.by_id", "admin_state_up", "false"),
				),
			},
		},
	})
}

var testAccDataSourceVpnServiceV2Config = fmt.Sprintf(`
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_gateway = "%s"
}
resource "opentelekomcloud_vpnaas_service_v2" "service_1" {
  router_id = "${opentelekomcloud_networking_router_v2.router_1.id}"
  admin_state_up = "false"
}

data "opentelekomcloud_vpnaas_service_v1" "by_id" {
  id = "${opentelekomcloud_vpnaas_service_v2.service_1.id}"
}
`, OS_EXTGW_ID)
