package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceTaurusdbMysqlConfigurations_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_taurusdb_mysql_configurations_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusdbMysqlConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusdbMysqlConfigurationsDataSourceID(dataSource),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.user_defined"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.updated_at"),
				),
			},
		},
	})
}

func testAccCheckTaurusdbMysqlConfigurationsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB MySQL configurations data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the TaurusDB MySQL configurations data source ID not set")
		}

		return nil
	}
}

func testDataSourceTaurusdbMysqlConfigurations_basic() string {
	return `
data "opentelekomcloud_taurusdb_mysql_configurations_v3" "test" {}
`
}
