package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceGeminiDBQuotas_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_gemini_quotas_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiDBQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBQuotasDataSourceID(dataSource),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBQuotasDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GeminiDB quotas data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GeminiDB quotas data source ID not set")
		}

		return nil
	}
}

func testDataSourceGeminiDBQuotas_basic() string {
	return `
data "opentelekomcloud_gemini_quotas_v3" "test" {}
`
}
