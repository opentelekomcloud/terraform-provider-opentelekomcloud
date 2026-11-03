package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const resourceRocketMQConsumptionVerifyV2Name = "opentelekomcloud_dms_rocketmq_consumption_verify_v2.test"

// The API rejects verification for a consumer group without online consumers, and terraform can't start a consumer.
func TestAccDmsRocketMQConsumptionVerifyV2_groupOffline(t *testing.T) {
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDmsRocketMQConsumptionVerifyV2GroupOffline(name),
				ExpectError: regexp.MustCompile(`consumer group .+ not online`),
			},
		},
	})
}

func TestAccDmsRocketMQConsumptionVerifyV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckDmsRocketMQConsumer(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQConsumptionVerifyV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRocketMQConsumptionVerifyV2Name, "instance_id",
						env.OS_DMS_ROCKETMQ_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceRocketMQConsumptionVerifyV2Name, "group",
						env.OS_DMS_ROCKETMQ_GROUP),
					resource.TestCheckResourceAttr(resourceRocketMQConsumptionVerifyV2Name, "client_id",
						env.OS_DMS_ROCKETMQ_CLIENT_ID),
					resource.TestCheckResourceAttr(resourceRocketMQConsumptionVerifyV2Name, "msg_id_list.#", "1"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumptionVerifyV2Name, "resend_results.#", "1"),
					resource.TestCheckResourceAttrPair(resourceRocketMQConsumptionVerifyV2Name, "resend_results.0.msg_id",
						resourceRocketMQMessageSendV2Name, "msg_id"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumptionVerifyV2Name, "resend_results.0.error_code", "DMS.00000000"),
				),
			},
		},
	})
}

func testAccDmsRocketMQConsumptionVerifyV2GroupOffline(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_consumer_group_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name        = "%[2]s"
}

resource "opentelekomcloud_dms_rocketmq_consumption_verify_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  group       = opentelekomcloud_dms_rocketmq_consumer_group_v2.test.name
  topic       = opentelekomcloud_dms_rocketmq_topic_v2.test.name
  client_id   = "offline-client"
  msg_id_list = [opentelekomcloud_dms_rocketmq_message_send_v2.test.msg_id]
}
`, testAccDmsRocketMQMessageSendV2Basic(name), name)
}

func testAccDmsRocketMQConsumptionVerifyV2Basic() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_dms_rocketmq_message_send_v2" "test" {
  instance_id = "%[1]s"
  topic       = "%[3]s"
  body        = "terraform acceptance test"
}

resource "opentelekomcloud_dms_rocketmq_consumption_verify_v2" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topic       = "%[3]s"
  client_id   = "%[4]s"
  msg_id_list = [opentelekomcloud_dms_rocketmq_message_send_v2.test.msg_id]
}
`, env.OS_DMS_ROCKETMQ_INSTANCE_ID, env.OS_DMS_ROCKETMQ_GROUP, env.OS_DMS_ROCKETMQ_TOPIC, env.OS_DMS_ROCKETMQ_CLIENT_ID)
}
