package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCEAddonTemplateV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_cce_addon_template_v3.template"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonTemplateMappingV3DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonTemplateV3DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_versions", "v1.(15|17|19|21).*"),
					resource.TestCheckResourceAttr(dataSourceName, "image_version", ""),
					resource.TestCheckResourceAttr(dataSourceName, "swr_addr", "100.125.7.25:20202"),
					resource.TestCheckResourceAttr(dataSourceName, "swr_user", "hwofficial"),
				),
			},
		},
	})
}

func testAccCheckCCEAddonTemplateV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find templates data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster data source ID not set ")
		}

		return nil
	}
}

const testAccCCEAddonTemplateMappingV3DataSourceBasic = `
data "opentelekomcloud_cce_addon_template_v3" "template" {
  addon_version = "1.2.9"
  addon_name    = "gpu-beta"
}
`
