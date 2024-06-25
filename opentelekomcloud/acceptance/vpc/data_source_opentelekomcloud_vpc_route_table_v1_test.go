package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcRouteTableDataSource_basic(t *testing.T) {
	rName := tools.RandomString("rtb-", 5)
	dataSourceName := "data.opentelekomcloud_vpc_route_table_v1.drtb"
	dc := common.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteTable_base(rName),
			},
			{
				Config: testAccDataSourceRouteTable_default(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "default", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.#", "1"),
				),
			},
			{
				Config: testAccDataSourceRouteTable_custom(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "default", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceRouteTable_base(rName string) string {
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
`, rName)
}

func testAccDataSourceRouteTable_default(rName string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_vpc_route_table_v1" "drtb" {
  vpc_id = opentelekomcloud_vpc_v1.vpc.id
}
`, testAccDataSourceRouteTable_base(rName))
}

func testAccDataSourceRouteTable_custom(rName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_vpc_route_table_v1" "rtb" {
  name        = "%[2]s"
  vpc_id      = opentelekomcloud_vpc_v1.vpc.id
  description = "created by terraform"
}

data "opentelekomcloud_vpc_route_table_v1" "drtb" {
  vpc_id = opentelekomcloud_vpc_v1.vpc.id
  name   = "%[2]s"

  depends_on = [opentelekomcloud_vpc_route_table_v1.rtb]
}
`, testAccDataSourceRouteTable_base(rName), rName)
}
