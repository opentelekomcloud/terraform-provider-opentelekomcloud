package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceLtsStructuringCustomTemplatesV3_basic(t *testing.T) {
	const dataSource = "data.opentelekomcloud_lts_structuring_custom_templates_v3.test"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{{
			Config: testDataSourceLtsStructuringCustomTemplatesV3(""),
			Check: resource.ComposeTestCheckFunc(
				dc.CheckResourceExists(),
				resource.TestCheckResourceAttrSet(dataSource, "region"),
			),
		}},
	})
}

func TestAccDataSourceLtsStructuringCustomTemplatesV3_byID(t *testing.T) {
	templateID := os.Getenv("OS_LTS_STRUCTURING_TEMPLATE_ID")
	if templateID == "" {
		t.Skip("OS_LTS_STRUCTURING_TEMPLATE_ID must be set for this acceptance test")
	}

	const dataSource = "data.opentelekomcloud_lts_structuring_custom_templates_v3.test"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{{
			Config: testDataSourceLtsStructuringCustomTemplatesV3(templateID),
			Check: resource.ComposeTestCheckFunc(
				dc.CheckResourceExists(),
				resource.TestCheckResourceAttr(dataSource, "id", templateID),
				resource.TestCheckResourceAttr(dataSource, "templates.#", "1"),
				resource.TestCheckResourceAttr(dataSource, "templates.0.id", templateID),
			),
		}},
	})
}

func testDataSourceLtsStructuringCustomTemplatesV3(templateID string) string {
	if templateID == "" {
		return `data "opentelekomcloud_lts_structuring_custom_templates_v3" "test" {}`
	}
	return fmt.Sprintf(`
data "opentelekomcloud_lts_structuring_custom_templates_v3" "test" {
  id = %q
}
`, templateID)
}
