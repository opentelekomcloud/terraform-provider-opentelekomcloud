package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcSubnetIdsV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetIdV2DataSource_vpcsubnet,
			},
			{
				Config: testAccSubnetIdV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccSubnetIdV2DataSourceID("data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids", "ids.#", "1"),
				),
			},
		},
	})
}
func testAccSubnetIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find vpc subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc Subnet data source ID not set")
		}

		return nil
	}
}

const testAccSubnetIdV2DataSource_vpcsubnet = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "test_vpc"
	cidr= "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
}
`

var testAccSubnetIdV2DataSource_basic = fmt.Sprintf(`
%s
data "opentelekomcloud_vpc_subnet_ids_v1" "subnet_ids" {
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
}
`, testAccSubnetIdV2DataSource_vpcsubnet)
