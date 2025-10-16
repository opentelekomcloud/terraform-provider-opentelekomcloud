package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenTelekomCloudTaurusDBMysqlConfigurationDataSource_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_taurusdb_mysql_configuration_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudTaurusDBMysqlConfigurationDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBMysqlConfigurationDataSourceID(dataSource),
					resource.TestCheckResourceAttr(
						dataSource,
						"name",
						"Default-TaurusDB V2.0"),
					resource.TestCheckResourceAttrSet(
						dataSource,
						"description"),
					resource.TestCheckResourceAttrSet(
						dataSource,
						"datastore_version"),
					resource.TestCheckResourceAttrSet(
						dataSource,
						"datastore_name"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBMysqlConfigurationDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB MySQL configuration data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the TaurusDB MySQL configuration data source ID not set")
		}

		return nil
	}
}

const testAccOpenTelekomCloudTaurusDBMysqlConfigurationDataSource_basic = `
data "opentelekomcloud_taurusdb_mysql_configuration_v3" "test" {
  name = "Default-TaurusDB V2.0"
}
`
