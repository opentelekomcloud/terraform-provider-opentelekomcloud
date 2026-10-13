package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/topics"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceRocketMQTopicV2Name = "opentelekomcloud_dms_rocketmq_topic_v2.test"

func getDmsRocketMQTopicFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %s", err)
	}
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	return topics.Get(client, parts[0], parts[1])
}

func TestAccDmsRocketMQTopicV2_basic(t *testing.T) {
	var topic topics.Topic
	name := fmt.Sprintf("rocketmq-acc-%s", acctest.RandString(5))

	rc := common.InitResourceCheck(resourceRocketMQTopicV2Name, &topic, getDmsRocketMQTopicFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQTopicV2Basic(name, "topic test"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceRocketMQTopicV2Name, "instance_id",
						resourceRocketMQInstanceV2Name, "id"),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "name", name),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "queue_num", "4"),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "message_type", "FIFO"),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "description", "topic test"),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "permission", "all"),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "total_read_queue_num", "4"),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "total_write_queue_num", "4"),
				),
			},
			{
				Config: testAccDmsRocketMQTopicV2Basic(name, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRocketMQTopicV2Name, "description", ""),
				),
			},
			{
				ResourceName:            resourceRocketMQTopicV2Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"queue_num", "description"},
			},
		},
	})
}

func testAccDmsRocketMQTopicV2Basic(name, description string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_dms_rocketmq_topic_v2" "test" {
  instance_id  = opentelekomcloud_dms_rocketmq_instance_v2.test.id
  name         = "%s"
  queue_num    = 4
  message_type = "FIFO"
  description  = "%s"
}
`, testAccDmsRocketMQInstanceV2Basic(name), name, description)
}
