package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceRocketMQMessageSendV2Name = "opentelekomcloud_dms_rocketmq_message_send_v2.test"

func TestAccDmsRocketMQMessageSendV2_basic(t *testing.T) {
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQMessageSendV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceRocketMQMessageSendV2Name, "instance_id",
						resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(resourceRocketMQMessageSendV2Name, "topic", name),
					resource.TestCheckResourceAttr(resourceRocketMQMessageSendV2Name, "body", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceRocketMQMessageSendV2Name, "property_list.#", "2"),
					resource.TestCheckResourceAttrPair(resourceRocketMQMessageSendV2Name, "msg_id",
						resourceRocketMQMessageSendV2Name, "id"),
					resource.TestCheckResourceAttrSet(resourceRocketMQMessageSendV2Name, "queue_id"),
					resource.TestCheckResourceAttrSet(resourceRocketMQMessageSendV2Name, "queue_offset"),
					resource.TestCheckResourceAttrSet(resourceRocketMQMessageSendV2Name, "broker_name"),
				),
			},
		},
	})
}

func testAccDmsRocketMQMessageSendV2Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_topic_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = "%[2]s"
  queue_num   = 4
}

resource "opentelekomcloud_dms_rocketmq_message_send_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  topic       = opentelekomcloud_dms_rocketmq_topic_v2.test.name
  body        = "terraform acceptance test"

  property_list {
    name  = "KEYS"
    value = "owner"
  }
  property_list {
    name  = "TAGS"
    value = "terraform"
  }
}
`, testAccDmsRocketMQInstanceV2Basic(name), name)
}
