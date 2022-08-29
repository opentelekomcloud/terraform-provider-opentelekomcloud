package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataFlavors = "data.opentelekomcloud_lb_flavors_v3.flavors_names"

func TestAccELBV3DataSourceFlavors_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3DataSourceFlavorsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBV3DataSourceFlavors(dataFlavors),
				),
			},
		},
	})
}

func testAccCheckELBV3DataSourceFlavors(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find flavors data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("flavors data source ID not set ")
		}

		return nil
	}
}

var testAccElbV3DataSourceFlavorsBasic = fmt.Sprintf(`
data "opentelekomcloud_lb_flavors_v3" "flavors_names" {}
`)
