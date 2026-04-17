package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenTelekomCloudTaurusDBMysqlFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudTaurusDBMysqlFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBMysqlFlavorsDataSourceID("data.opentelekomcloud_taurusdb_mysql_flavors_v3.flavor"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBMysqlFlavorsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB MySQL data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the TaurusDB MySQL data source ID not set")
		}

		return nil
	}
}

const testAccOpenTelekomCloudTaurusDBMysqlFlavorsDataSource_basic = `
data "opentelekomcloud_taurusdb_mysql_flavors_v3" "flavor" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = "multi"
}
`
