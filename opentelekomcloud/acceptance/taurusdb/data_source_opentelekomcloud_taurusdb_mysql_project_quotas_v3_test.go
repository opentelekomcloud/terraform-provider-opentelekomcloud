package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccTaurusDBMysqlProjectQuotasDataSource_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_taurusdb_mysql_project_quotas_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBMysqlProjectQuotasDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBMysqlProjectQuotasDataSourceID(dataSource),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.quota"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBMysqlProjectQuotasDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB MySQL project quotas data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the TaurusDB MySQL project quotas data source ID not set")
		}

		return nil
	}
}

const testAccTaurusDBMysqlProjectQuotasDataSource_basic = `
data "opentelekomcloud_taurusdb_mysql_project_quotas_v3" "test" {}

locals {
  type = "instance"
}

data "opentelekomcloud_taurusdb_mysql_project_quotas_v3" "type_filter" {
  type = "instance"
}

output "type_filter_is_useful" {
  value = length(data.opentelekomcloud_taurusdb_mysql_project_quotas_v3.type_filter.quotas) > 0 && alltrue(
    [for v in data.opentelekomcloud_taurusdb_mysql_project_quotas_v3.type_filter.quotas[*].resources : length(v) > 0 && alltrue(
      [for vv in v : vv.type == local.type]
    )]
  )
}
`
