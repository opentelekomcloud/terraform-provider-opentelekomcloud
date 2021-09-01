package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestRDSVersionsV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testRDSVersionsV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV3DataSourceID("data.opentelekomcloud_rds_versions_v3.sqls_versions"),
					testAccCheckRdsFlavorV3DataSourceID("data.opentelekomcloud_rds_versions_v3.mysql_versions"),
					testAccCheckRdsFlavorV3DataSourceID("data.opentelekomcloud_rds_versions_v3.psql_versions"),
				),
			},
		},
	})
}

const testRDSVersionsV3Basic = `
data "opentelekomcloud_rds_versions_v3" "sqls_versions" {
  database_name = "sqlserver"
}

data "opentelekomcloud_rds_versions_v3" "mysql_versions" {
  database_name = "mysql"
}

data "opentelekomcloud_rds_versions_v3" "psql_versions" {
  database_name = "postgresql"
}
`
