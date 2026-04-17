package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceGaussdbMysqlEngineVersions_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_taurusdb_mysql_engine_versions_v3.test"
	dc := common.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlEngineVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.name"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlEngineVersions_basic() string {
	return `
data "opentelekomcloud_taurusdb_mysql_engine_versions_v3" "test" {
  database_name = "gaussdb-mysql"
}
`
}
