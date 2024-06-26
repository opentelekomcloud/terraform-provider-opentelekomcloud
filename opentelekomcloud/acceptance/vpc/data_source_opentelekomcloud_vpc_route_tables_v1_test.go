package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestAccVpcRouteTablesV1DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesV1RouteTablesNoResources,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_tables_v1.all_rtbs", "routetables.#", "0"),
					testAccRouteTablesV1DataSourceID("data.opentelekomcloud_vpc_route_tables_v1.all_rtbs"),
				),
			},
			{
				Config: testAccRouteTablesV1RouteTablesFull,
				Check: resource.ComposeTestCheckFunc(
					// In total there should be 4 route tables: 2 defaults + 2 created from resources
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_tables_v1.all_rtbs", "routetables.#", "4"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_tables_v1.vpc_1_rtbs", "routetables.#", "2"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_tables_v1.vpc_2_rtbs", "routetables.#", "2"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_tables_v1.vpc_1_rtb_by_subnet_id", "routetables.#", "1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_route_tables_v1.vpc_1_rtb_by_id", "routetables.#", "1"),
					testAccRouteTablesV1DataSourceID("data.opentelekomcloud_vpc_route_tables_v1.all_rtbs"),
				),
			},
		},
	})
}

func testAccRouteTablesV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find route tables data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("route tables data source ID not set")
		}

		return nil
	}
}

const testAccRouteTablesV1RouteTablesNoResources = `
data "opentelekomcloud_vpc_route_tables_v1" "empty_rtbs" {
}
`

const testAccRouteTablesV1RouteTablesFull = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "test_route_tables_vpc_1"
  cidr = "192.168.0.0/24"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "test_route_tables_vpc_2"
  cidr = "10.0.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_1_subnet_1" {
  name       = "vpc-1-subnet-1"
  cidr       = "192.168.0.0/28"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_1_subnet_2" {
  name       = "vpc-1-subnet-2"
  cidr       = "192.168.0.128/28"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_2_subnet_1" {
  name       = "vpc-2-subnet-1"
  cidr       = "10.0.0.0/24"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_2_subnet_2" {
  name       = "vpc-2-subnet-2"
  cidr       = "10.0.100.0/24"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_nat_gateway_v2" "vpc_1_natgw" {
  name                = "vpc_1_natgw"
  internal_network_id = opentelekomcloud_vpc_subnet_v1.vpc_1_subnet_1.network_id
  spec      = 0
  router_id = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_route_table_v1" "vpc_1_table_1" {
  name        = "vpc_1_table_1"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  subnets = [
    opentelekomcloud_vpc_subnet_v1.vpc_1_subnet_1.network_id,
  ]
  route {
    destination = "0.0.0.0/0"
    type        = "nat"
    nexthop     = opentelekomcloud_nat_gateway_v2.vpc_1_natgw.id
    description = "vpc_1_natgw"
  }
}

resource "opentelekomcloud_vpc_route_table_v1" "vpc_2_table_1" {
  name        = "vpc_2_table_1"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_2.id
}

data "opentelekomcloud_vpc_route_tables_v1" "all_rtbs" {
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_1_rtbs" {
	vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_2_rtbs" {
	vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_1_rtb_by_subnet_id" {
	subnet_id = opentelekomcloud_vpc_route_table_v1.vpc_1_subnet_1.network_id
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_1_rtb_by_id" {
	id = opentelekomcloud_vpc_route_table_v1.vpc_1_table_1.id
}
`
