package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/streams"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccLogTankTopicV2_basic(t *testing.T) {
	var topic streams.LogStream
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLogTankTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankTopicV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankTopicV2Exists(
						"opentelekomcloud_logtank_topic_v2.testacc_topic", &topic),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_topic_v2.testacc_topic", "topic_name", "testacc_topic"),
				),
			},
		},
	})
}

func testAccCheckLogTankTopicV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	ltsclient, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_logtank_topic_v2" {
			continue
		}

		groupId := rs.Primary.Attributes["group_id"]
		allStreams, err := streams.ListLogStream(ltsclient, groupId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return nil
			} else {
				return err
			}
		}
		for _, stream := range allStreams {
			if stream.LogStreamId == rs.Primary.ID {
				return fmt.Errorf("log topic (%s) still exists", rs.Primary.ID)
			}
		}

		if _, ok := err.(golangsdk.ErrDefault400); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckLogTankTopicV2Exists(n string, topic *streams.LogStream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		ltsclient, err := config.LtsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
		}

		group_id := rs.Primary.Attributes["group_id"]

		allStreams, err := streams.ListLogStream(ltsclient, group_id)
		if err != nil {
			return err
		}

		for _, stream := range allStreams {
			if stream.LogStreamId == rs.Primary.ID {
				*topic = stream
				break
			}
		}

		return nil
	}
}

const testAccLogTankTopicV2_basic = `
resource "opentelekomcloud_logtank_group_v2" "testacc_group" {
  group_name  = "testacc_group"
  ttl_in_days = 7
}
resource "opentelekomcloud_logtank_topic_v2" "testacc_topic" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "testacc_topic"
}
`
