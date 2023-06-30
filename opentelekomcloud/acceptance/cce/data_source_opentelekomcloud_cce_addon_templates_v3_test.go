package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCEAddonTemplatesV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_cce_addon_templates_v3.templates"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonTemplatesV3DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonTemplatesV3DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addons.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.addon_version", "1.2.20"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.image_version", ""),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.swr_addr", "100.125.7.25:20202"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.swr_user", "cce-addons"),
				),
			},
			{
				Config: testAccCCEAddonTemplatesV3DataSourceListVersions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonTemplatesV3DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addons.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.addon_version", "1.3.7"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.swr_addr", "100.125.7.25:20202"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.swr_user", "hwofficial"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.platform", "linux-amd64"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.1.addon_version", "1.4.5"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.1.swr_addr", "100.125.7.25:20202"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.1.swr_user", "cce-addons"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.1.platform", "linux-amd64"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.2.addon_version", "1.7.1"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.2.swr_addr", "100.125.7.25:20202"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.2.swr_user", "cce-addons"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.2.platform", "linux-amd64"),
				),
			},
			{
				Config: testAccCCEAddonTemplatesV3DataSourceListVersionsReleaseCandidate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonTemplatesV3DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addons.#", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.addon_version", "1.0.10"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.obs_url", "obs.eu-de.otc.t-systems.com"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.euleros_version", "2.2.5"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.swr_addr", "100.125.7.25:20202"),
					resource.TestCheckResourceAttr(dataSourceName, "addons.0.swr_user", "hwofficial"),
				),
			},
		},
	})
}

func testAccCheckCCEAddonTemplatesV3DataSourceID(n string) resource.TestCheckFunc {
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

const testAccCCEAddonTemplatesV3DataSourceBasic = `
data "opentelekomcloud_cce_addon_templates_v3" "templates" {
  cluster_version = "1.25.2"
  addon_name      = "gpu-beta"
}
`

const testAccCCEAddonTemplatesV3DataSourceListVersions = `
data "opentelekomcloud_cce_addon_templates_v3" "templates" {
  cluster_version = "1.21.3"
  addon_name      = "volcano"
}
`

const testAccCCEAddonTemplatesV3DataSourceListVersionsReleaseCandidate = `
data "opentelekomcloud_cce_addon_templates_v3" "templates" {
  cluster_version = "1.11.6"
  addon_name      = "storage-driver"
}
`
