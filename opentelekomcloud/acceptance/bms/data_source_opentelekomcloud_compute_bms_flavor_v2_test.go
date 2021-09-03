package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataFlavorName = "data.opentelekomcloud_compute_bms_flavors_v2.flavor"

func TestAccBMSV2FlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccBmsFlavorPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSV2FlavorDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSV2FlavorDataSourceID(dataFlavorName),
					resource.TestCheckResourceAttr(dataFlavorName, "name", env.OS_BMS_FLAVOR_NAME),
				),
			},
		},
	})
}

func testAccCheckBMSV2FlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("flavor data source ID not set")
		}

		return nil
	}
}

var testAccBMSV2FlavorDataSourceBasic = fmt.Sprintf(`
data "opentelekomcloud_compute_bms_flavors_v2" "flavor" {
  name = "%s"
}
`, env.OS_BMS_FLAVOR_NAME)
