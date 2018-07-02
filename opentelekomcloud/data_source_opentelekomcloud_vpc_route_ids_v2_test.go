package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOTCVpcRouteIdsV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCRouteIdV2DataSource_vpcroute,
			},
			resource.TestStep{
				Config: testAccOTCRouteIdV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccOTCRouteIdV2DataSourceID("data.opentelekomcloud_vpc_route_ids_v2.route_ids"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_ids_v2.route_ids", "ids.#", "1"),
				),
			},
		},
	})
}
func testAccOTCRouteIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc route data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc Route data source ID not set")
		}

		return nil
	}
}

const testAccOTCRouteIdV2DataSource_vpcroute = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
name = "vpc_test"
cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
		name = "vpc_test1"
        cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
		name = "opentelekomcloud_peering"
		vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
		peer_vpc_id = "${opentelekomcloud_vpc_v1.vpc_2.id}"
}

resource "opentelekomcloud_vpc_route_v2" "route_1" {
   type = "peering"
  nexthop = "${opentelekomcloud_vpc_peering_connection_v2.peering_1.id}"
  destination = "192.168.0.0/16"
  vpc_id ="${opentelekomcloud_vpc_v1.vpc_1.id}"
}
`

var testAccOTCRouteIdV2DataSource_basic = fmt.Sprintf(`
%s
data "opentelekomcloud_vpc_route_ids_v2" "route_ids" {
  vpc_id = "${opentelekomcloud_vpc_route_v2.route_1.vpc_id}"
}
`, testAccOTCRouteIdV2DataSource_vpcroute)
