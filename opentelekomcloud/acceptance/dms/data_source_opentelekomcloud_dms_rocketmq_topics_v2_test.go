package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDmsRocketMQTopicsV2DataSource_basic(t *testing.T) {
	all := "data.opentelekomcloud_dms_rocketmq_topics_v2.all"
	byName := "data.opentelekomcloud_dms_rocketmq_topics_v2.by_name"
	byType := "data.opentelekomcloud_dms_rocketmq_topics_v2.by_message_type"
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))
	dc := common.InitDataSourceCheck(all)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQTopicsV2DataSourceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("all_contains_topics", "true"),
					resource.TestCheckResourceAttr(byName, "topics.#", "1"),
					resource.TestCheckResourceAttr(byName, "topics.0.name", name),
					resource.TestCheckResourceAttr(byName, "topics.0.message_type", "FIFO"),
					resource.TestCheckResourceAttr(byName, "topics.0.permission", "all"),
					resource.TestCheckResourceAttr(byName, "topics.0.total_read_queue_num", "4"),
					resource.TestCheckResourceAttr(byName, "topics.0.total_write_queue_num", "4"),
					resource.TestCheckResourceAttr(byType, "topics.#", "1"),
					resource.TestCheckResourceAttr(byType, "topics.0.name", name+"-2"),
				),
			},
		},
	})
}

func testAccDmsRocketMQTopicsV2DataSourceBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_topic_v2" "test_2" {
  instance_id  = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name         = "%[2]s-2"
  queue_num    = 2
  message_type = "DELAY"
}

data "opentelekomcloud_dms_rocketmq_topics_v2" "all" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id

  depends_on = [
    opentelekomcloud_dms_rocketmq_topic_v2.test,
    opentelekomcloud_dms_rocketmq_topic_v2.test_2,
  ]
}

data "opentelekomcloud_dms_rocketmq_topics_v2" "by_name" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = opentelekomcloud_dms_rocketmq_topic_v2.test.name

  depends_on = [opentelekomcloud_dms_rocketmq_topic_v2.test_2]
}

data "opentelekomcloud_dms_rocketmq_topics_v2" "by_message_type" {
  instance_id  = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  message_type = "DELAY"

  depends_on = [
    opentelekomcloud_dms_rocketmq_topic_v2.test,
    opentelekomcloud_dms_rocketmq_topic_v2.test_2,
  ]
}

output "all_contains_topics" {
  value = length(setintersection(
    data.opentelekomcloud_dms_rocketmq_topics_v2.all.topics[*].name, ["%[2]s", "%[2]s-2"]
  )) == 2
}
`, testAccDmsRocketMQTopicV2Basic(name, "topic test"), name)
}
