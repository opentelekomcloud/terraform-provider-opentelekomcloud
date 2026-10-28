package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceRocketMQMessageOffsetResetV2Name = "opentelekomcloud_dms_rocketmq_message_offset_reset_v2.test"

func TestAccDmsRocketMQMessageOffsetResetV2_basic(t *testing.T) {
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQMessageOffsetResetV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceRocketMQMessageOffsetResetV2Name, "instance_id",
						resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(resourceRocketMQMessageOffsetResetV2Name, "group", name),
					resource.TestCheckResourceAttr(resourceRocketMQMessageOffsetResetV2Name, "topic", name),
					resource.TestCheckResourceAttr(resourceRocketMQMessageOffsetResetV2Name, "timestamp", "1662652800000"),
					resource.TestCheckResourceAttrSet(resourceRocketMQMessageOffsetResetV2Name, "queues.#"),
				),
			},
		},
	})
}

func testAccDmsRocketMQMessageOffsetResetV2Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_topic_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = "%[2]s"
  queue_num   = 4
}

resource "opentelekomcloud_dms_rocketmq_message_offset_reset_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  group       = opentelekomcloud_dms_rocketmq_consumer_group_v2.test.name
  topic       = opentelekomcloud_dms_rocketmq_topic_v2.test.name
  timestamp   = "1662652800000"
}
`, testAccDmsRocketMQConsumerGroupV2Basic(name), name)
}
