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

const dataRocketMQConsumersV2Name = "data.opentelekomcloud_dms_rocketmq_consumers_v2.test"

// The API rejects the consumer list query for a consumer group without online consumers,
// and terraform can't start a consumer.
func TestAccDmsRocketMQConsumersV2DataSource_groupOffline(t *testing.T) {
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDmsRocketMQConsumersV2DataSourceGroupOffline(name),
				ExpectError: regexp.MustCompile(`consumer group .+ not online`),
			},
		},
	})
}

// Requires a consumer connected to the OS_DMS_ROCKETMQ_GROUP group and subscribed to the OS_DMS_ROCKETMQ_TOPIC topic.
func TestAccDmsRocketMQConsumersV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckDmsRocketMQConsumer(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQConsumersV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataRocketMQConsumersV2Name, "online", "true"),
					resource.TestCheckResourceAttr(dataRocketMQConsumersV2Name, "subscription_consistency", "true"),
					resource.TestCheckTypeSetElemNestedAttrs(dataRocketMQConsumersV2Name, "clients.*", map[string]string{
						"client_id":             env.OS_DMS_ROCKETMQ_CLIENT_ID,
						"subscriptions.0.topic": env.OS_DMS_ROCKETMQ_TOPIC,
					}),
					resource.TestCheckResourceAttrSet(dataRocketMQConsumersV2Name, "clients.0.language"),
					resource.TestCheckResourceAttrSet(dataRocketMQConsumersV2Name, "clients.0.version"),
					resource.TestCheckResourceAttrSet(dataRocketMQConsumersV2Name, "clients.0.client_addr"),
				),
			},
		},
	})
}

func testAccDmsRocketMQConsumersV2DataSourceGroupOffline(name string) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_dms_rocketmq_consumers_v2" "test" {
  instance_id = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  group       = opentelekomcloud_dms_rocketmq_consumer_group_v2.test.name
}
`, testAccDmsRocketMQConsumerGroupV2Basic(name))
}

func testAccDmsRocketMQConsumersV2DataSourceBasic() string {
	return fmt.Sprintf(`
data "opentelekomcloud_dms_rocketmq_consumers_v2" "test" {
  instance_id = "%s"
  group       = "%s"
  is_detail   = true
}
`, env.OS_DMS_ROCKETMQ_INSTANCE_ID, env.OS_DMS_ROCKETMQ_GROUP)
}
