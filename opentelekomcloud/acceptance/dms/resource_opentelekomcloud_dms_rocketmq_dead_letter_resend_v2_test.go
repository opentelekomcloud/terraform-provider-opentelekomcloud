package acceptance

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const resourceRocketMQDeadLetterResendV2Name = "opentelekomcloud_dms_rocketmq_dead_letter_resend_v2.test"

// The API rejects resending for a consumer group without online consumers, and terraform can't start a consumer.
func TestAccDmsRocketMQDeadLetterResendV2_groupOffline(t *testing.T) {
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDmsRocketMQDeadLetterResendV2GroupOffline(name),
				ExpectError: regexp.MustCompile(`consumer group .+ not online`),
			},
		},
	})
}

func TestAccDmsRocketMQDeadLetterResendV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckDmsRocketMQConsumer(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQDeadLetterResendV2Message(),
			},
			{
				// The consumer has to fail the message, so it lands in the dead letter topic.
				PreConfig: func() { time.Sleep(time.Minute) },
				Config:    testAccDmsRocketMQDeadLetterResendV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRocketMQDeadLetterResendV2Name, "instance_id",
						env.OS_DMS_ROCKETMQ_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceRocketMQDeadLetterResendV2Name, "topic",
						"%DLQ%"+env.OS_DMS_ROCKETMQ_GROUP),
					resource.TestCheckResourceAttr(resourceRocketMQDeadLetterResendV2Name, "msg_id_list.#", "1"),
					resource.TestCheckResourceAttr(resourceRocketMQDeadLetterResendV2Name, "resend_results.#", "1"),
					resource.TestCheckResourceAttrPair(resourceRocketMQDeadLetterResendV2Name, "resend_results.0.msg_id",
						resourceRocketMQMessageSendV2Name, "msg_id"),
					resource.TestCheckResourceAttr(resourceRocketMQDeadLetterResendV2Name, "resend_results.0.error_code", "DMS.00000000"),
				),
			},
		},
	})
}

func testAccDmsRocketMQDeadLetterResendV2GroupOffline(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_consumer_group_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = "%[2]s"
}

resource "opentelekomcloud_dms_rocketmq_dead_letter_resend_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  topic       = "%%DLQ%%${opentelekomcloud_dms_rocketmq_consumer_group_v2.test.name}"
  msg_id_list = [opentelekomcloud_dms_rocketmq_message_send_v2.test.msg_id]
}
`, testAccDmsRocketMQMessageSendV2Basic(name), name)
}

// The consumer of OS_DMS_ROCKETMQ_GROUP must fail messages with the "dlq" tag
// and accept messages of the %DLQ%<group> topic, which are resent directly to it.
// The group must have retry_max_time = 1, so the failed message lands in the dead letter topic within the sleep.
func testAccDmsRocketMQDeadLetterResendV2Message() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_rocketmq_message_send_v2" "test" {
  instance_id = "%s"
  topic       = "%s"
  body        = "terraform acceptance test"

  property_list {
    name  = "TAGS"
    value = "dlq"
  }
}
`, env.OS_DMS_ROCKETMQ_INSTANCE_ID, env.OS_DMS_ROCKETMQ_TOPIC)
}

func testAccDmsRocketMQDeadLetterResendV2Basic() string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_dead_letter_resend_v2" "test" {
  instance_id = "%s"
  topic       = "%%DLQ%%%s"
  msg_id_list = [opentelekomcloud_dms_rocketmq_message_send_v2.test.msg_id]
}
`, testAccDmsRocketMQDeadLetterResendV2Message(), env.OS_DMS_ROCKETMQ_INSTANCE_ID, env.OS_DMS_ROCKETMQ_GROUP)
}
