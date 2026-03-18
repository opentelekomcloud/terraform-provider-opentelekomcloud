package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/routetables"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVPCRouteTableRouteName = "opentelekomcloud_vpc_route_table_route_v1.route"

func TestAccVpcRouteTableRouteV1_basic(t *testing.T) {
	name := tools.RandomString("rtbr-", 5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteTableRouteV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableRoute_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "destination", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "type", "peering"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "description", "peering route"),
					resource.TestCheckResourceAttrSet(resourceVPCRouteTableRouteName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceVPCRouteTableRouteName, "route_table_name"),
				),
			},
			{
				Config: testAccVpcRouteTableRoute_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "destination", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "type", "peering"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "description", "peering route updated"),
				),
			},
			{
				ResourceName:      resourceVPCRouteTableRouteName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcRouteTableRouteImportStateID(resourceVPCRouteTableRouteName),
			},
		},
	})
}

func TestAccVpcRouteTableRouteV1_withRouteTable(t *testing.T) {
	name := tools.RandomString("rtbr-", 5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteTableRouteV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableRoute_withRouteTable(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "destination", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableRouteName, "type", "peering"),
					resource.TestCheckResourceAttrSet(resourceVPCRouteTableRouteName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceVPCRouteTableRouteName, "route_table_name"),
				),
			},
			{
				ResourceName:      resourceVPCRouteTableRouteName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcRouteTableRouteImportStateID(resourceVPCRouteTableRouteName),
			},
		},
	})
}

func testAccCheckRouteTableRouteV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_route_table_route_v1" {
			continue
		}

		rtID := rs.Primary.Attributes["route_table_id"]
		destination := rs.Primary.Attributes["destination"]

		routeTable, err := routetables.Get(client, rtID)
		if err != nil {
			continue
		}

		for _, route := range routeTable.Routes {
			if route.DestinationCIDR == destination {
				return fmt.Errorf("route %s still exists in route table %s", destination, rtID)
			}
		}
	}

	return nil
}

func testAccVpcRouteTableRouteImportStateID(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		rtID := rs.Primary.Attributes["route_table_id"]
		destination := rs.Primary.Attributes["destination"]
		return rtID + "/" + destination, nil
	}
}

func testAccVpcRouteTableRoute_network(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "%[1]s-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "%[1]s-2"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "%[1]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
`, name)
}

func testAccVpcRouteTableRoute_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_route_table_route_v1" "route" {
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  destination = "172.16.0.0/16"
  type        = "peering"
  nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
  description = "peering route"
}
`, testAccVpcRouteTableRoute_network(name))
}

func testAccVpcRouteTableRoute_update(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_route_table_route_v1" "route" {
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  destination = "172.16.0.0/16"
  type        = "peering"
  nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
  description = "peering route updated"
}
`, testAccVpcRouteTableRoute_network(name))
}

func testAccVpcRouteTableRoute_withRouteTable(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_route_table_v1" "table" {
  name   = "%[2]s"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_route_table_route_v1" "route" {
  vpc_id         = opentelekomcloud_vpc_v1.vpc_1.id
  route_table_id = opentelekomcloud_vpc_route_table_v1.table.id
  destination    = "172.16.0.0/16"
  type           = "peering"
  nexthop        = opentelekomcloud_vpc_peering_connection_v2.peering.id
  description    = "peering route on custom table"
}
`, testAccVpcRouteTableRoute_network(name), name)
}
