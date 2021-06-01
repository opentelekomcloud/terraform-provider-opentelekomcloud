package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcSubnetV1DataSource_basic(t *testing.T) {
	dataSourceNameByID := "data.opentelekomcloud_vpc_subnet_v1.by_id"
	dataSourceNameByCIDR := "data.opentelekomcloud_vpc_subnet_v1.by_cidr"
	dataSourceNameByName := "data.opentelekomcloud_vpc_subnet_v1.by_name"
	dataSourceNameByVPC := "data.opentelekomcloud_vpc_subnet_v1.by_vpc_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { common.TestAccPreCheck(t) },
		Providers: common.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcSubnetV1Config,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceVpcSubnetV1Check(dataSourceNameByID, "test_subnet", "10.0.0.0/24",
						"10.0.0.1", "eu-de-02"),
					testAccDataSourceVpcSubnetV1Check(dataSourceNameByCIDR, "test_subnet", "10.0.0.0/24",
						"10.0.0.1", "eu-de-02"),
					testAccDataSourceVpcSubnetV1Check(dataSourceNameByName, "test_subnet", "10.0.0.0/24",
						"10.0.0.1", "eu-de-02"),
					testAccDataSourceVpcSubnetV1Check(dataSourceNameByVPC, "test_subnet", "10.0.0.0/24",
						"10.0.0.1", "eu-de-02"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "dhcp_enable", "true"),
				),
			},
		},
	})
}

func testAccDataSourceVpcSubnetV1Check(n, name, cidr, gatewayIP, availabilityZone string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		subnetRs, ok := s.RootModule().Resources["opentelekomcloud_vpc_subnet_v1.subnet_1"]
		if !ok {
			return fmt.Errorf("can't find opentelekomcloud_vpc_subnet_v1.subnet_1 in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != subnetRs.Primary.Attributes["id"] {
			return fmt.Errorf("bad id %s", attr["id"])
		}

		if attr["cidr"] != cidr {
			return fmt.Errorf("bad subnet cidr %s, expected: %s", attr["cidr"], cidr)
		}
		if attr["name"] != name {
			return fmt.Errorf("bad subnet name %s", attr["name"])
		}
		if attr["gateway_ip"] != gatewayIP {
			return fmt.Errorf("bad subnet gateway_ip %s", attr["gateway_ip"])
		}
		if attr["availability_zone"] != availabilityZone {
			return fmt.Errorf("bad subnet availability_zone %s", attr["availability_zone"])
		}

		return nil
	}
}

const testAccDataSourceVpcSubnetV1Config = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "test_vpc"
  cidr= "10.0.0.0/24"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "test_subnet"
  cidr              = "10.0.0.0/24"
  gateway_ip        = "10.0.0.1"
  vpc_id            = opentelekomcloud_vpc_v1.vpc_1.id
  availability_zone = "eu-de-02"
}

data "opentelekomcloud_vpc_subnet_v1" "by_id" {
  id = opentelekomcloud_vpc_subnet_v1.subnet_1.id
}

data "opentelekomcloud_vpc_subnet_v1" "by_cidr" {
  cidr = opentelekomcloud_vpc_subnet_v1.subnet_1.cidr
}

data "opentelekomcloud_vpc_subnet_v1" "by_name" {
  name = opentelekomcloud_vpc_subnet_v1.subnet_1.name
}

data "opentelekomcloud_vpc_subnet_v1" "by_vpc_id" {
  vpc_id = opentelekomcloud_vpc_subnet_v1.subnet_1.vpc_id
}
`
