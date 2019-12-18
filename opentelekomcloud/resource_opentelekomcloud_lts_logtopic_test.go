package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/lts/v2/logtopics"
)

func TestAccLogTankTopicV2_basic(t *testing.T) {
	var topic logtopics.LogTopic
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogTankTopicV2Destroy,
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
	config := testAccProvider.Meta().(*Config)
	ltsclient, err := config.ltsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_logtank_topic_v2" {
			continue
		}

		group_id := rs.Primary.Attributes["group_id"]
		_, err = logtopics.Get(ltsclient, group_id, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Log topic (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckLogTankTopicV2Exists(n string, topic *logtopics.LogTopic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ltsclient, err := config.ltsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud LTS client: %s", err)
		}

		var found *logtopics.LogTopic
		group_id := rs.Primary.Attributes["group_id"]

		found, err = logtopics.Get(ltsclient, group_id, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*topic = *found

		return nil
	}
}

const testAccLogTankTopicV2_basic = `
resource "opentelekomcloud_logtank_group_v2" "testacc_group" {
    group_name  = "testacc_group"
}
resource "opentelekomcloud_logtank_topic_v2" "testacc_topic" {
  group_id = "${opentelekomcloud_logtank_group_v2.testacc_group.id}"
  topic_name = "testacc_topic"
}
`
