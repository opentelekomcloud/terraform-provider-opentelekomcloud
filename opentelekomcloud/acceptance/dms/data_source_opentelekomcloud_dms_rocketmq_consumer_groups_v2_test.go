package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsRocketMQConsumerGroupsV2DataSource_basic(t *testing.T) {
	all := "data.opentelekomcloud_dms_rocketmq_consumer_groups_v2.all"
	byName := "data.opentelekomcloud_dms_rocketmq_consumer_groups_v2.by_name"
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))
	dc := common.InitDataSourceCheck(all)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQConsumerGroupsV2DataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("all_contains_groups", "true"),
					resource.TestCheckResourceAttr(byName, "groups.#", "1"),
					resource.TestCheckResourceAttr(byName, "groups.0.name", name),
					resource.TestCheckResourceAttr(byName, "groups.0.description", "group test"),
					resource.TestCheckResourceAttr(byName, "groups.0.enabled", "true"),
					resource.TestCheckResourceAttr(byName, "groups.0.broadcast", "false"),
					resource.TestCheckResourceAttr(byName, "groups.0.retry_max_time", "16"),
					resource.TestCheckResourceAttr(byName, "groups.0.consume_orderly", "true"),
					resource.TestCheckResourceAttr(byName, "groups.0.group_online", "false"),
					resource.TestCheckResourceAttrSet(byName, "groups.0.created_at"),
				),
			},
		},
	})
}

func testAccDmsRocketMQConsumerGroupsV2DataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_consumer_group_v2" "test_2" {
  instance_id    = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name           = "%[2]s-2"
  retry_max_time = 3
}

data "opentelekomcloud_dms_rocketmq_consumer_groups_v2" "all" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id

  depends_on = [
    opentelekomcloud_dms_rocketmq_consumer_group_v2.test,
    opentelekomcloud_dms_rocketmq_consumer_group_v2.test_2,
  ]
}

data "opentelekomcloud_dms_rocketmq_consumer_groups_v2" "by_name" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = opentelekomcloud_dms_rocketmq_consumer_group_v2.test.name

  depends_on = [opentelekomcloud_dms_rocketmq_consumer_group_v2.test_2]
}

output "all_contains_groups" {
  value = length(setintersection(
    data.opentelekomcloud_dms_rocketmq_consumer_groups_v2.all.groups[*].name, ["%[2]s", "%[2]s-2"]
  )) == 2
}
`, testAccDmsRocketMQConsumerGroupV2Basic(name), name)
}
