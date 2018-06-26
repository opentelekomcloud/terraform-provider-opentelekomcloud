package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOTCVpcSubnetV1DataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOTCVpcSubnetV1Config,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOTCVpcSubnetV1Check("data.opentelekomcloud_vpc_subnet_v1.by_id", "opentelekomcloud_subnet", "192.168.0.0/16",
						"192.168.0.1", "eu-de-02"),
					testAccDataSourceOTCVpcSubnetV1Check("data.opentelekomcloud_vpc_subnet_v1.by_name", "opentelekomcloud_subnet", "192.168.0.0/16",
						"192.168.0.1", "eu-de-02"),
					testAccDataSourceOTCVpcSubnetV1Check("data.opentelekomcloud_vpc_subnet_v1.by_vpc_id", "opentelekomcloud_subnet", "192.168.0.0/16",
						"192.168.0.1", "eu-de-02"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_subnet_v1.by_id", "status", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_subnet_v1.by_id", "dhcp_enable", "true"),
				),
			},
		},
	})
}

func testAccDataSourceOTCVpcSubnetV1Check(n, name, cidr, gateway_ip, availability_zone string) resource.TestCheckFunc {
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
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				subnetRs.Primary.Attributes["id"],
			)
		}

		if attr["cidr"] != cidr {
			return fmt.Errorf("bad subnet cidr %s, expected: %s", attr["cidr"], cidr)
		}
		if attr["name"] != name {
			return fmt.Errorf("bad subnet name %s", attr["name"])
		}
		if attr["gateway_ip"] != gateway_ip {
			return fmt.Errorf("bad subnet gateway_ip %s", attr["gateway_ip"])
		}
		if attr["availability_zone"] != availability_zone {
			return fmt.Errorf("bad subnet availability_zone %s", attr["availability_zone"])
		}

		return nil
	}
}

const testAccDataSourceOTCVpcSubnetV1Config = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "test_vpc"
	cidr= "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "eu-de-02"
 }

data "opentelekomcloud_vpc_subnet_v1" "by_id" {
  id = "${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
}

data "opentelekomcloud_vpc_subnet_v1" "by_name" {
	name = "${opentelekomcloud_vpc_subnet_v1.subnet_1.name}"
}

data "opentelekomcloud_vpc_subnet_v1" "by_vpc_id" {
	vpc_id = "${opentelekomcloud_vpc_subnet_v1.subnet_1.vpc_id}"
}
`
