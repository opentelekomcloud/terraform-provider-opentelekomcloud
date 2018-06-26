package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOTCVpcSubnetIdsV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCSubnetIdV2DataSource_vpcsubnet,
			},
			resource.TestStep{
				Config: testAccOTCSubnetIdV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccOTCSubnetIdV2DataSourceID("data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids", "ids.#", "1"),
				),
			},
		},
	})
}
func testAccOTCSubnetIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc Subnet data source ID not set")
		}

		return nil
	}
}

const testAccOTCSubnetIdV2DataSource_vpcsubnet = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "test_vpc"
	cidr= "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
}
`

var testAccOTCSubnetIdV2DataSource_basic = fmt.Sprintf(`
%s
data "opentelekomcloud_vpc_subnet_ids_v1" "subnet_ids" {
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
}
`, testAccOTCSubnetIdV2DataSource_vpcsubnet)
