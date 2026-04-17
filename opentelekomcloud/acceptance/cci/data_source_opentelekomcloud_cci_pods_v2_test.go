package cci

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCIPodsV2DataSource_basic(t *testing.T) {
	cciImage := os.Getenv("OS_CCI_IMAGE")
	if cciImage == "" {
		t.Skip("OS_CCI_IMAGE is not set, skipping CCI pods data source test")
	}

	rName := fmt.Sprintf("cci-pod-%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cci_pods_v2.by_name"
	dataSourceAll := "data.opentelekomcloud_cci_pods_v2.all"

	dc := common.InitDataSourceCheck(dataSourceName)
	dcAll := common.InitDataSourceCheck(dataSourceAll)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPodsV2DataSource_basic(rName, cciImage),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "namespace", rName),
					resource.TestCheckResourceAttr(dataSourceName, "pods.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.namespace", rName),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.image", cciImage),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.name", "c1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pods.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pods.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pods.0.resource_version"),
					dcAll.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceAll, "pods.#"),
				),
			},
		},
	})
}

func testAccCCIPodsV2DataSource_basic(rName, image string) string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_cci_pods_v2" "by_name" {
  depends_on = [opentelekomcloud_cci_pod_v2.test]
  namespace  = opentelekomcloud_cci_namespace_v2.test.name
  name       = opentelekomcloud_cci_pod_v2.test.name
}

data "opentelekomcloud_cci_pods_v2" "all" {
  depends_on = [opentelekomcloud_cci_pod_v2.test]
  namespace  = opentelekomcloud_cci_namespace_v2.test.name
}
`, testAccV2Pod_basic(rName, image))
}
