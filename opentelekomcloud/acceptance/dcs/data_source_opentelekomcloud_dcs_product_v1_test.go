package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataProductName = "data.opentelekomcloud_dcs_product_v1.product1"

func TestAccDcsProductV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsProductV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsProductV1DataSourceID(dataProductName),
					resource.TestCheckResourceAttr(dataProductName, "spec_code", "dcs.single_node"),
				),
			},
		},
	})
}

func testAccCheckDcsProductV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find DCS product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("dcs product data source ID not set")
		}

		return nil
	}
}

const testAccDcsProductV1DataSourceBasic = `
data "opentelekomcloud_dcs_product_v1" "product1" {
  spec_code = "dcs.single_node"
}
`
