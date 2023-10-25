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

const resourceVPCRouteTableName = "opentelekomcloud_vpc_route_table_v1.table_1"

func TestAccVpcRouteTableV1_basic(t *testing.T) {
	var rtb routetables.RouteTable
	rtbName := tools.RandomString("rtb-", 5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteTableV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTable_basic(rtbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableV1Exists(resourceVPCRouteTableName, &rtb),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.#", "0"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "subnets.#", "0"),
				),
			},
			{
				Config: testAccVpcRouteTable_peering(rtbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "description", "created by terraform with routes"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.#", "1"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.0.destination", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.0.type", "peering"),
				),
			},
			{
				Config: testAccVpcRouteTable_associate(rtbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "description", "created by terraform with subnets"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.#", "1"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.0.destination", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.0.type", "peering"),
				),
			},
			{
				ResourceName:      resourceVPCRouteTableName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVpcRouteTableV1_multipleRoutes(t *testing.T) {
	var rtb routetables.RouteTable
	rtbName := tools.RandomString("rtb-", 5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteTableV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTable_multiRoutes(rtbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableV1Exists(resourceVPCRouteTableName, &rtb),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "name", rtbName),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "description", "created by terraform with multi routes"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "route.#", "6"),
					resource.TestCheckResourceAttr(resourceVPCRouteTableName, "subnets.#", "0"),
				),
			},
		},
	})
}

func testAccCheckRouteTableV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_route_table_v1" {
			continue
		}

		_, err := routetables.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("route table still exists")
		}
	}

	return nil
}

func testAccCheckRouteTableV1Exists(n string, route *routetables.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %s", err)
		}

		found, err := routetables.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("route table not found")
		}

		*route = *found

		return nil
	}
}

func testAccVpcRouteTable_network(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "%[1]s-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1-1" {
  name       = "%[1]s-1-1"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1-2" {
  name       = "%[1]s-1-2"
  cidr       = "192.168.10.0/24"
  gateway_ip = "192.168.10.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "%[1]s-2"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_2-1" {
  name       = "%[1]s-2-1"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_2.id
}
`, name)
}

func testAccVpcRouteTable_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_route_table_v1" "table_1" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  description = "created by terraform"
}
`, testAccVpcRouteTable_network(name), name)
}

func testAccVpcRouteTable_peering(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table_1" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  description = "created by terraform with routes"

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering rule"
  }
}
`, testAccVpcRouteTable_network(name), name)
}

func testAccVpcRouteTable_associate(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table_1" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  description = "created by terraform with subnets"

  subnets     = [
    opentelekomcloud_vpc_subnet_v1.subnet_1-1.id,
    opentelekomcloud_vpc_subnet_v1.subnet_1-2.id,
  ]

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering rule"
  }
}
`, testAccVpcRouteTable_network(name), name)
}

func testAccVpcRouteTable_multiRoutes(name string) string {
	return fmt.Sprintf(`
%s

%[3]s

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor   = "s2.medium.1"
  vpc_id   = opentelekomcloud_vpc_v1.vpc_1.id

  nics {
    network_id = opentelekomcloud_vpc_subnet_v1.subnet_1-1.network_id
  }

  data_disks {
    size = 10
    type = "SSD"
  }

  password                    = "Password@123"
  availability_zone           = "%[4]s"
  auto_recovery               = true
  delete_disks_on_termination = true

}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table_1" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  description = "created by terraform with multi routes"

  route {
    destination = "172.16.1.0/24"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering one rule"
  }
  route {
    destination = "172.16.2.0/24"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering two rule"
  }
  route {
    destination = "172.16.3.0/24"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering three rule"
  }
  route {
    destination = "172.16.4.0/24"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering four rule"
  }
  route {
    destination = "172.16.5.0/24"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering five rule"
  }
  route {
    destination = "172.16.6.0/24"
    type        = "ecs"
    nexthop     = opentelekomcloud_ecs_instance_v1.instance_1.id
    description = "ecs six rule"
  }
}
`, testAccVpcRouteTable_network(name), name, common.DataSourceImage, env.OS_AVAILABILITY_ZONE)
}
