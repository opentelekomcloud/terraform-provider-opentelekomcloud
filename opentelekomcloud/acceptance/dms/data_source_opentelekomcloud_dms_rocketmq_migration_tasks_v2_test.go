package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsRocketMQMigrationTasksV2DataSource_basic(t *testing.T) {
	all := "data.opentelekomcloud_dms_rocketmq_migration_tasks_v2.all"
	byID := "data.opentelekomcloud_dms_rocketmq_migration_tasks_v2.by_id"
	byName := "data.opentelekomcloud_dms_rocketmq_migration_tasks_v2.by_name"
	byType := "data.opentelekomcloud_dms_rocketmq_migration_tasks_v2.by_type"
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))
	dc := common.InitDataSourceCheck(all)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQMigrationTasksV2DataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "tasks.#", "1"),
					resource.TestCheckResourceAttr(byID, "tasks.#", "1"),
					resource.TestCheckResourceAttr(byID, "tasks.0.name", name),
					resource.TestCheckResourceAttr(byID, "tasks.0.type", "rocketmq"),
					resource.TestCheckResourceAttr(byID, "tasks.0.status", "finished"),
					resource.TestCheckResourceAttrSet(byID, "tasks.0.start_date"),
					resource.TestCheckResourceAttr(byName, "tasks.#", "1"),
					resource.TestCheckResourceAttr(byType, "tasks.#", "0"),
				),
			},
		},
	})
}

func testAccDmsRocketMQMigrationTasksV2DataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_dms_rocketmq_migration_tasks_v2" "all" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id

  depends_on = [opentelekomcloud_dms_rocketmq_migration_task_v2.test]
}

data "opentelekomcloud_dms_rocketmq_migration_tasks_v2" "by_id" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  task_id     = split("/", opentelekomcloud_dms_rocketmq_migration_task_v2.test.id)[1]
}

data "opentelekomcloud_dms_rocketmq_migration_tasks_v2" "by_name" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = opentelekomcloud_dms_rocketmq_migration_task_v2.test.name
}

data "opentelekomcloud_dms_rocketmq_migration_tasks_v2" "by_type" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  type        = "rabbitToRocket"

  depends_on = [opentelekomcloud_dms_rocketmq_migration_task_v2.test]
}
`, testAccDmsRocketMQMigrationTaskV2Basic(name))
}
