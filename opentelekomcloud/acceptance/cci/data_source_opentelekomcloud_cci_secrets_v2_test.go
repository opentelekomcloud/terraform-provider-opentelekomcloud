package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCCISecretsV2DataSource_basic(t *testing.T) {
	nsName := fmt.Sprintf("cci-ns-%s", acctest.RandString(5))
	rName := fmt.Sprintf("cci-secret-%s", acctest.RandString(5))
	dataSourceName := "data.opentelekomcloud_cci_secrets_v2.by_name"
	dataSourceAll := "data.opentelekomcloud_cci_secrets_v2.all"

	dc := common.InitDataSourceCheck(dataSourceName)
	dcAll := common.InitDataSourceCheck(dataSourceAll)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCISecretsV2DataSource_basic(nsName, rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "namespace", nsName),
					resource.TestCheckResourceAttr(dataSourceName, "secrets.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "secrets.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "secrets.0.namespace", nsName),
					resource.TestCheckResourceAttr(dataSourceName, "secrets.0.type", "Opaque"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secrets.0.resource_version"),
					dcAll.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceAll, "secrets.#"),
				),
			},
		},
	})
}

func testAccCCISecretsV2DataSource_basic(nsName, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_cci_secret_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "%[2]s"
  type      = "Opaque"

  data = {
    "key1" = "dmFsdWUx"
  }
}

data "opentelekomcloud_cci_secrets_v2" "by_name" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = opentelekomcloud_cci_secret_v2.test.name
}

data "opentelekomcloud_cci_secrets_v2" "all" {
  depends_on = [opentelekomcloud_cci_secret_v2.test]
  namespace  = opentelekomcloud_cci_namespace_v2.test.name
}
`, testAccV2Namespace_basic(nsName), rName)
}
