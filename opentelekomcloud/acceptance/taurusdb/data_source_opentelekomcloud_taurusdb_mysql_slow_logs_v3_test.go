package acceptance

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataSourceTaurusDBSlowLogs = "data.opentelekomcloud_taurusdb_mysql_slow_logs_v3.test"

func TestAccTaurusDBMysqlSlowLogsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBMysqlSlowLogsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBMysqlSlowLogsDataSourceID(dataSourceTaurusDBSlowLogs),
					resource.TestCheckResourceAttrSet(dataSourceTaurusDBSlowLogs, "slow_log_list.#"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBMysqlSlowLogsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB MySQL slow logs data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the TaurusDB MySQL slow logs data source ID not set")
		}

		return nil
	}
}

func testAccTaurusDBMysqlSlowLogsDataSource_basic() string {
	endDate := time.Now().Format("2006-01-02T15:04:05-0700")
	startDate := time.Now().Add(-1 * time.Hour).Format("2006-01-02T15:04:05-0700")

	return fmt.Sprintf(`
data "opentelekomcloud_taurusdb_mysql_slow_logs_v3" "test" {
  instance_id = "%s"
  node_id     = "%s"
  start_date  = "%s"
  end_date    = "%s"
}
`, os.Getenv("OS_TAURUSDB_INSTANCE_ID"), os.Getenv("OS_TAURUSDB_NODE_ID"), startDate, endDate)
}
