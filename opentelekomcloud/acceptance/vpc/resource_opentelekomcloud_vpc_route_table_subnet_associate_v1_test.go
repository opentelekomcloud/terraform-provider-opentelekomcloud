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

const resourceVPCRouteTableSubnetAssociateName = "opentelekomcloud_vpc_route_table_subnet_associate_v1.assoc"

func TestAccVpcRouteTableSubnetAssociateV1_basic(t *testing.T) {
	name := tools.RandomString("rtba-", 5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteTableSubnetAssociateV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableSubnetAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceVPCRouteTableSubnetAssociateName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceVPCRouteTableSubnetAssociateName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceVPCRouteTableSubnetAssociateName, "vpc_id"),
				),
			},
			{
				ResourceName:      resourceVPCRouteTableSubnetAssociateName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcRouteTableSubnetAssociateImportStateID(resourceVPCRouteTableSubnetAssociateName),
			},
		},
	})
}

func testAccCheckRouteTableSubnetAssociateV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_route_table_subnet_associate_v1" {
			continue
		}

		rtID := rs.Primary.Attributes["route_table_id"]
		subnetID := rs.Primary.Attributes["subnet_id"]

		routeTable, err := routetables.Get(client, rtID)
		if err != nil {
			continue
		}

		for _, subnet := range routeTable.Subnets {
			if subnet.ID == subnetID {
				return fmt.Errorf("subnet %s still associated with route table %s", subnetID, rtID)
			}
		}
	}

	return nil
}

func testAccVpcRouteTableSubnetAssociateImportStateID(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		rtID := rs.Primary.Attributes["route_table_id"]
		subnetID := rs.Primary.Attributes["subnet_id"]
		return rtID + "/" + subnetID, nil
	}
}

func testAccVpcRouteTableSubnetAssociate_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table" {
  name   = "%[1]s"
  vpc_id = opentelekomcloud_vpc_v1.vpc.id

  lifecycle {
    ignore_changes = [subnets]
  }
}

resource "opentelekomcloud_vpc_route_table_subnet_associate_v1" "assoc" {
  route_table_id = opentelekomcloud_vpc_route_table_v1.table.id
  subnet_id      = opentelekomcloud_vpc_subnet_v1.subnet.id
}
`, name)
}
