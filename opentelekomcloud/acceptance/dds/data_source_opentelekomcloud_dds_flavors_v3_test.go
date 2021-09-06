package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataFlavorName = "data.opentelekomcloud_dds_flavors_v3.flavor"

func TestAccDDSFlavorV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSFlavorV3DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSFlavorV3DataSourceID(dataFlavorName),
					resource.TestCheckResourceAttrSet(dataFlavorName, "flavors.#"),
				),
			},
		},
	})
}

func testAccCheckDDSFlavorV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find DDS Flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DDS Flavor data source ID not set")
		}

		return nil
	}
}

var testAccDDSFlavorV3DataSourceBasic = `
data "opentelekomcloud_dds_flavors_v3" "flavor" {
  engine_name = "DDS-Community"
  vcpus       = 8
}
`
