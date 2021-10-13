package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/alarmnotifications"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/waf"
)

const resourceAlarmNotificationName = "opentelekomcloud_waf_alarm_notification_v1.notification_1"

func TestAccWafAlarmNotificationV1_basic(t *testing.T) {
	var notification alarmnotifications.AlarmNotification

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafAlarmNotificationV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafAlarmNotificationV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafAlarmNotificationV1Exists(resourceAlarmNotificationName, &notification),
					resource.TestCheckResourceAttr(resourceAlarmNotificationName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceAlarmNotificationName, "send_frequency", "30"),
					resource.TestCheckResourceAttr(resourceAlarmNotificationName, "times", "200"),
					resource.TestCheckResourceAttrSet(resourceAlarmNotificationName, "topic_urn"),
				),
			},
			{
				Config: testAccWafAlarmNotificationV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafAlarmNotificationV1Exists(resourceAlarmNotificationName, &notification),
					resource.TestCheckResourceAttr(resourceAlarmNotificationName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceAlarmNotificationName, "send_frequency", "15"),
					resource.TestCheckResourceAttr(resourceAlarmNotificationName, "times", "100"),
					resource.TestCheckResourceAttrSet(resourceAlarmNotificationName, "topic_urn"),
				),
			},
		},
	})
}

func testAccCheckWafAlarmNotificationV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.WafV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(waf.WafClientError, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_waf_alarm_notification_v1" {
			continue
		}

		alarmNotification, err := alarmnotifications.List(client).Extract()
		if err != nil {
			return fmt.Errorf("error receiving alarm notification: %w", err)
		}
		if alarmNotification.Enabled {
			return fmt.Errorf("alarm notification still enabled")
		}
	}

	return nil
}

func testAccCheckWafAlarmNotificationV1Exists(n string, notification *alarmnotifications.AlarmNotification) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.WafV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(waf.WafClientError, err)
		}

		found, err := alarmnotifications.List(client).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("alarm notification not found")
		}

		*notification = *found

		return nil
	}
}

const testAccWafAlarmNotificationV1Basic = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name = "topic_alarm"
}

resource "opentelekomcloud_waf_alarm_notification_v1" "notification_1" {
  enabled        = true
  topic_urn      = opentelekomcloud_smn_topic_v2.topic_1.id
  send_frequency = 30
  times          = 200
  threats        = ["cc", "cmdi"]
}
`

const testAccWafAlarmNotificationV1Update = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name = "topic_alarm"
}

resource "opentelekomcloud_waf_alarm_notification_v1" "notification_1" {
  enabled        = true
  topic_urn      = opentelekomcloud_smn_topic_v2.topic_1.id
  send_frequency = 15
  times          = 100
  threats        = ["all"]
}
`
