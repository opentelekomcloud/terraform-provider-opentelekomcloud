package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceTaurusDBMysqlBackups_basic(t *testing.T) {
	dataSource := "data.opentelekomcloud_taurusdb_mysql_backups_v3.test"
	rName := common.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBMysqlBackups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.take_up_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),

					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("backup_id_filter_is_useful", "true"),
					resource.TestCheckOutput("backup_type_filter_is_useful", "true"),
					resource.TestCheckOutput("begin_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBMysqlBackups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "opentelekomcloud_taurusdb_mysql_backups_v3" "test" {
  depends_on = [opentelekomcloud_taurusdb_mysql_backup_v3.test]
}

locals {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
}

data "opentelekomcloud_taurusdb_mysql_backups_v3" "instance_id_filter" {
  depends_on  = [opentelekomcloud_taurusdb_mysql_backup_v3.test]
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
}

output "instance_id_filter_is_useful" {
  value = length(data.opentelekomcloud_taurusdb_mysql_backups_v3.instance_id_filter.backups) > 0 && alltrue(
    [for v in data.opentelekomcloud_taurusdb_mysql_backups_v3.instance_id_filter.backups[*].instance_id : v == local.instance_id]
  )
}

locals {
  backup_id = opentelekomcloud_taurusdb_mysql_backup_v3.test.id
}

data "opentelekomcloud_taurusdb_mysql_backups_v3" "backup_id_filter" {
  depends_on = [opentelekomcloud_taurusdb_mysql_backup_v3.test]
  backup_id  = opentelekomcloud_taurusdb_mysql_backup_v3.test.id
}

output "backup_id_filter_is_useful" {
  value = length(data.opentelekomcloud_taurusdb_mysql_backups_v3.backup_id_filter.backups) > 0 && alltrue(
    [for v in data.opentelekomcloud_taurusdb_mysql_backups_v3.backup_id_filter.backups[*].id : v == local.backup_id]
  )
}

locals {
  backup_type = "manual"
}

data "opentelekomcloud_taurusdb_mysql_backups_v3" "backup_type_filter" {
  depends_on  = [opentelekomcloud_taurusdb_mysql_backup_v3.test]
  backup_type = "manual"
}

output "backup_type_filter_is_useful" {
  value = length(data.opentelekomcloud_taurusdb_mysql_backups_v3.backup_type_filter.backups) > 0 && alltrue(
    [for v in data.opentelekomcloud_taurusdb_mysql_backups_v3.backup_type_filter.backups[*].type : v == local.backup_type]
  )
}

locals {
  begin_time = opentelekomcloud_taurusdb_mysql_backup_v3.test.begin_time
}

data "opentelekomcloud_taurusdb_mysql_backups_v3" "begin_time_filter" {
  depends_on = [opentelekomcloud_taurusdb_mysql_backup_v3.test]
  begin_time = opentelekomcloud_taurusdb_mysql_backup_v3.test.begin_time
}

output "begin_time_filter_is_useful" {
  value = length(data.opentelekomcloud_taurusdb_mysql_backups_v3.begin_time_filter.backups) > 0 && alltrue(
    [for v in data.opentelekomcloud_taurusdb_mysql_backups_v3.begin_time_filter.backups[*].begin_time : v == local.begin_time]
  )
}
`, testBackup_basic(name))
}
