package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRocketMQConsumerGroupV2Name = "opentelekomcloud_dms_rocketmq_consumer_group_v2.test"

func getDmsRocketMQConsumerGroupFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %s", err)
	}
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	return groups.Get(client, parts[0], parts[1])
}

func TestAccDmsRocketMQConsumerGroupV2_basic(t *testing.T) {
	var group groups.ConsumerGroup
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(resourceRocketMQConsumerGroupV2Name, &group, getDmsRocketMQConsumerGroupFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQConsumerGroupV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceRocketMQConsumerGroupV2Name, "instance_id",
						resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "name", name),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "broadcast", "false"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "retry_max_time", "16"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "consume_orderly", "true"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "description", "group test"),
				),
			},
			{
				Config: testAccDmsRocketMQConsumerGroupV2Update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "broadcast", "true"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "retry_max_time", "5"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "consume_orderly", "true"),
					resource.TestCheckResourceAttr(resourceRocketMQConsumerGroupV2Name, "description", "group test updated"),
				),
			},
			{
				ResourceName:      resourceRocketMQConsumerGroupV2Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDmsRocketMQConsumerGroupV2Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_consumer_group_v2" "test" {
  instance_id     = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name            = "%s"
  retry_max_time  = 16
  consume_orderly = true
  description     = "group test"
}
`, testAccDmsRocketMQInstanceV2Basic(name), name)
}

func testAccDmsRocketMQConsumerGroupV2Update(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_consumer_group_v2" "test" {
  instance_id     = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name            = "%s"
  enabled         = false
  broadcast       = true
  retry_max_time  = 5
  consume_orderly = true
  description     = "group test updated"
}
`, testAccDmsRocketMQInstanceV2Basic(name), name)
}
