package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestAccBandWidthV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataName := "data.opentelekomcloud_vpc_bandwidth_v2.test"

	t.Parallel()
	quotas.BookOne(t, quotas.SharedBandwidth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthDataSourceV2Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandWidthV2DataSourceExists(dataName),
					resource.TestCheckResourceAttr(dataName, "name", rName),
					resource.TestCheckResourceAttr(dataName, "size", "10"),
				),
			},
		},
	})
}

func testAccCheckBandWidthV2DataSourceExists(n string) resource.TestCheckFunc { // nolint:unused
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		bandwidthRs, ok := s.RootModule().Resources["opentelekomcloud_vpc_bandwidth_v2.test"]
		if !ok {
			return fmt.Errorf("can't find opentelekomcloud_vpc_bandwidth_v2.test in state")
		}

		attr := rs.Primary.Attributes
		if attr["id"] != bandwidthRs.Primary.ID {
			return fmt.Errorf("attribute 'id' expected %s; got %s",
				bandwidthRs.Primary.Attributes["id"], attr["id"])
		}

		return nil
	}
}

func testAccBandWidthDataSourceV2Basic(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_bandwidth_v2" "test" {
  name = "%s"
  size = 10
}

data "opentelekomcloud_vpc_bandwidth_v2" "test" {
  name = opentelekomcloud_vpc_bandwidth_v2.test.name
}
`, rName)
}
