package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

const (
	resourceRTVpc1       = "data.opentelekomcloud_vpc_route_tables_v1.vpc_1_rtbs"
	resourceRTVpc2       = "data.opentelekomcloud_vpc_route_tables_v1.vpc_2_rtbs"
	resourceRTBySubnetId = "data.opentelekomcloud_vpc_route_tables_v1.vpc_1_rtb_by_subnet_id"
	resourceRTById       = "data.opentelekomcloud_vpc_route_tables_v1.vpc_1_rtb_by_id"
)

func TestAccVpcRouteTablesV1DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesV1RouteTablesFull(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRTVpc1, "routetables.#", "1"),
					resource.TestCheckResourceAttr(resourceRTVpc2, "routetables.#", "1"),
					resource.TestCheckResourceAttr(resourceRTBySubnetId, "routetables.#", "1"),
					resource.TestCheckResourceAttr(resourceRTById, "routetables.#", "1"),
					resource.TestCheckResourceAttr(resourceRTById, "routetables.0.routes.0.type", "nat"),
					testAccRouteTablesV1DataSourceID(resourceRTById),
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

const testAccRouteTablesV1Base = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "test_route_tables_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_1_subnet_1" {
  name       = "vpc-1-subnet-1"
  cidr       = "192.168.100.0/24"
  gateway_ip = "192.168.100.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_1_subnet_2" {
  name       = "vpc-1-subnet-2"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "test_route_tables_vpc_2"
  cidr = "10.0.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_2_subnet_1" {
  name       = "vpc-2-subnet-1"
  cidr       = "10.0.0.0/24"
  gateway_ip = "10.0.0.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_subnet_v1" "vpc_2_subnet_2" {
  name       = "vpc-2-subnet-2"
  cidr       = "10.0.100.0/24"
  gateway_ip = "10.0.100.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_nat_gateway_v2" "vpc_1_natgw" {
  name                = "vpc_1_natgw"
  internal_network_id = opentelekomcloud_vpc_subnet_v1.vpc_1_subnet_1.network_id
  spec                = 0
  router_id           = opentelekomcloud_vpc_v1.vpc_1.id
}
`

func testAccRouteTablesV1RouteTablesFull() string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_route_table_v1" "vpc_1_table_1" {
  name   = "vpc_1_table_1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
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
  name   = "vpc_2_table_1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_1_rtbs" {
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_2_rtbs" {
  vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_1_rtb_by_subnet_id" {
  subnet_id = opentelekomcloud_vpc_subnet_v1.vpc_1_subnet_1.network_id
}

data "opentelekomcloud_vpc_route_tables_v1" "vpc_1_rtb_by_id" {
  id = opentelekomcloud_vpc_route_table_v1.vpc_1_table_1.id
}
`, testAccRouteTablesV1Base)
}
