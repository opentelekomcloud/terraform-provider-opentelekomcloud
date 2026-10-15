package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/migration"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRocketMQMigrationTaskV2Name = "opentelekomcloud_dms_rocketmq_migration_task_v2.test"

func getDmsRocketMQMigrationTaskFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %s", err)
	}
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	return migration.Get(client, parts[0], parts[1], migration.GetOpts{})
}

func TestAccDmsRocketMQMigrationTaskV2_basic(t *testing.T) {
	var task migration.TaskDetails
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(resourceRocketMQMigrationTaskV2Name, &task, getDmsRocketMQMigrationTaskFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQMigrationTaskV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceRocketMQMigrationTaskV2Name, "instance_id",
						resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(resourceRocketMQMigrationTaskV2Name, "name", name),
					resource.TestCheckResourceAttr(resourceRocketMQMigrationTaskV2Name, "type", "rocketmq"),
					resource.TestCheckResourceAttr(resourceRocketMQMigrationTaskV2Name, "status", "finished"),
					resource.TestCheckResourceAttrSet(resourceRocketMQMigrationTaskV2Name, "start_date"),
				),
			},
			{
				ResourceName:            resourceRocketMQMigrationTaskV2Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite", "topic_configs", "subscription_groups"},
			},
		},
	})
}

func testAccDmsRocketMQMigrationTaskV2Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_migration_task_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = "%[2]s"
  type        = "rocketmq"
  overwrite   = true

  topic_configs {
    topic_name        = "topic-%[2]s"
    perm              = 6
    read_queue_nums   = 3
    write_queue_nums  = 3
    topic_filter_type = "SINGLE_TAG"
  }

  subscription_groups {
    group_name      = "group-%[2]s"
    retry_max_times = 2
  }
}
`, testAccDmsRocketMQInstanceV2Basic(name), name)
}
