package acceptance

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcV1DataSource_basic(t *testing.T) {
	dataSourceNameByID := "data.opentelekomcloud_vpc_v1.by_id"
	dataSourceNameByCidr := "data.opentelekomcloud_vpc_v1.by_cidr"
	dataSourceNameByName := "data.opentelekomcloud_vpc_v1.by_name"

	cidr := fmt.Sprintf("172.16.%d.0/24", rand.Intn(50))
	name := tools.RandomString("vpc-test-", 3)

	t.Parallel()
	th.AssertNoErr(t, quotas.Router.Acquire())
	defer quotas.Router.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcV1Config(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceVpcV1Check(dataSourceNameByID, name, cidr),
					testAccDataSourceVpcV1Check(dataSourceNameByCidr, name, cidr),
					testAccDataSourceVpcV1Check(dataSourceNameByName, name, cidr),
					resource.TestCheckResourceAttr(dataSourceNameByID, "shared", "false"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcV1Check(n, name, cidr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		vpcRs, ok := s.RootModule().Resources["opentelekomcloud_vpc_v1.vpc_1"]
		if !ok {
			return fmt.Errorf("can't find opentelekomcloud_vpc_v1.vpc_1 in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf("ID is: %s, expected: %s", attr["id"], vpcRs.Primary.Attributes["id"])
		}

		if attr["cidr"] != cidr {
			return fmt.Errorf("bad VPC cidr: %s, expected: %s", attr["cidr"], cidr)
		}
		if attr["name"] != name {
			return fmt.Errorf("bad VPC name: %s, expected: %s", attr["name"], name)
		}

		return nil
	}
}

func testAccDataSourceVpcV1Config(name, cidr string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "%s"
  cidr = "%s"
}

data "opentelekomcloud_vpc_v1" "by_id" {
  id = opentelekomcloud_vpc_v1.vpc_1.id
}

data "opentelekomcloud_vpc_v1" "by_cidr" {
  cidr = opentelekomcloud_vpc_v1.vpc_1.cidr
}

data "opentelekomcloud_vpc_v1" "by_name" {
  name = opentelekomcloud_vpc_v1.vpc_1.name
}
`, name, cidr)
}
