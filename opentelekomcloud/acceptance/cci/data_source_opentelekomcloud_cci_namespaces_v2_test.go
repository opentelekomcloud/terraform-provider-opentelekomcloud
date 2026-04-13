package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCINamespacesV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("cci-ns-%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cci_namespaces_v2.by_name"
	dataSourceAll := "data.opentelekomcloud_cci_namespaces_v2.all"

	dc := common.InitDataSourceCheck(dataSourceName)
	dcAll := common.InitDataSourceCheck(dataSourceAll)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCINamespacesV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.status", "Active"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.uid"),
					dcAll.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceAll, "namespaces.#"),
				),
			},
		},
	})
}

func testAccCCINamespacesV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_cci_namespace_v2" "test" {
  name = "%s"
}

data "opentelekomcloud_cci_namespaces_v2" "by_name" {
  name = opentelekomcloud_cci_namespace_v2.test.name
}

data "opentelekomcloud_cci_namespaces_v2" "all" {
  depends_on = [opentelekomcloud_cci_namespace_v2.test]
}
`, rName)
}
