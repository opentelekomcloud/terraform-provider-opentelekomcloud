package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcRouteV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteV2DataSourceID("data.opentelekomcloud_vpc_route_v2.by_id"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.by_id", "type", "peering"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.by_id", "destination", "192.168.0.0/16"),
					testAccCheckRouteV2DataSourceID("data.opentelekomcloud_vpc_route_v2.by_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.by_vpc_id", "type", "peering"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_route_v2.by_vpc_id", "destination", "192.168.0.0/16"),
				),
			},
		},
	})
}

func testAccCheckRouteV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find vpc route connection data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc route connection data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceRouteV2Config = `
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
		vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
		peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_v2" "route_1" {
   type = "peering"
  nexthop = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
  destination = "192.168.0.0/16"
  vpc_id =opentelekomcloud_vpc_v1.vpc_1.id
}

data "opentelekomcloud_vpc_route_v2" "by_id" {
		id = opentelekomcloud_vpc_route_v2.route_1.id
}

data "opentelekomcloud_vpc_route_v2" "by_vpc_id" {
		vpc_id = opentelekomcloud_vpc_route_v2.route_1.vpc_id
}
`
