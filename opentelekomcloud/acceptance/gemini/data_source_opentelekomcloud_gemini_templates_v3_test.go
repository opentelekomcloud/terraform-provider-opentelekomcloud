package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceGeminiDBTemplates_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_gemini_templates_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeminiDBTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBTemplatesDataSourceID(dataSource),
					resource.TestCheckResourceAttrSet(dataSource, "templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.user_defined"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.mode"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBTemplatesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GeminiDB templates data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GeminiDB templates data source ID not set")
		}

		return nil
	}
}

func testDataSourceGeminiDBTemplates_basic() string {
	return `
data "opentelekomcloud_gemini_templates_v3" "test" {}
`
}
