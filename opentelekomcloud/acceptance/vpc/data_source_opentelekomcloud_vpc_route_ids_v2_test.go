package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcRouteIdsV2DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteIdV2DataSourceVpcRoute,
			},
			{
				Config: testAccRouteIdV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccRouteIdV2DataSourceID("data.opentelekomcloud_vpc_route_ids_v2.route_ids"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_ids_v2.route_ids", "ids.#", "1"),
				),
			},
		},
	})
}

func testAccRouteIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find vpc route data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc Route data source ID not set")
		}

		return nil
	}
}

const testAccRouteIdV2DataSourceVpcRoute = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_ds_ids"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_ds_ids1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name        = "opentelekomcloud_peering"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_v2" "route_1" {
  type        = "peering"
  nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering_1.id
  destination = "192.168.0.0/16"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
}
`

var testAccRouteIdV2DataSourceBasic = fmt.Sprintf(`
%s
data "opentelekomcloud_vpc_route_ids_v2" "route_ids" {
  vpc_id = opentelekomcloud_vpc_route_v2.route_1.vpc_id
}
`, testAccRouteIdV2DataSourceVpcRoute)
