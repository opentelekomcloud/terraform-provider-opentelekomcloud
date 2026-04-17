package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/oneclickalarms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getOneClickAlarmV2Func(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CesV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CES v2 client: %s", err)
	}
	alarmList, err := oneclickalarms.List(c)
	if err != nil {
		return nil, err
	}
	for _, a := range alarmList {
		if a.OneClickAlarmId == state.Primary.ID {
			return a, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func TestCESOneClickAlarmV2_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("ces-oca-%s", acctest.RandString(5))
	rName := "opentelekomcloud_ces_one_click_alarm_v2.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getOneClickAlarmV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESOneClickAlarmV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "one_click_alarm_id", "OBSSystemOneClickAlarm"),
					resource.TestCheckResourceAttr(rName, "dimension_names.0.metric.0", "bucket_name"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "internal_alarm_id"),
					resource.TestCheckResourceAttrSet(rName, "namespace"),
					resource.TestCheckResourceAttrSet(rName, "description"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				Config: testCESOneClickAlarmV2Update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "one_click_alarm_id", "OBSSystemOneClickAlarm"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_notifications.#", "1"),
					resource.TestCheckResourceAttr(rName, "alarm_notifications.0.type", "notification"),
					resource.TestCheckResourceAttrSet(rName, "internal_alarm_id"),
					resource.TestCheckResourceAttrSet(rName, "namespace"),
					resource.TestCheckResourceAttrSet(rName, "description"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"one_click_alarm_id", "dimension_names", "notification_enabled",
					"alarm_notifications", "ok_notifications", "notification_begin_time", "notification_end_time",
				},
			},
		},
	})
}

func TestCESOneClickAlarmV2_withNotifications(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("ces-oca-%s", acctest.RandString(5))
	rName := "opentelekomcloud_ces_one_click_alarm_v2.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getOneClickAlarmV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESOneClickAlarmV2WithNotifications(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "one_click_alarm_id", "OBSSystemOneClickAlarm"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_notifications.#", "1"),
					resource.TestCheckResourceAttr(rName, "alarm_notifications.0.type", "notification"),
					resource.TestCheckResourceAttr(rName, "ok_notifications.#", "1"),
					resource.TestCheckResourceAttr(rName, "ok_notifications.0.type", "notification"),
					resource.TestCheckResourceAttr(rName, "notification_begin_time", "00:00"),
					resource.TestCheckResourceAttr(rName, "notification_end_time", "23:59"),
					resource.TestCheckResourceAttrSet(rName, "internal_alarm_id"),
					resource.TestCheckResourceAttrSet(rName, "namespace"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"one_click_alarm_id", "dimension_names", "notification_enabled",
					"alarm_notifications", "ok_notifications", "notification_begin_time", "notification_end_time",
				},
			},
		},
	})
}

func TestCESOneClickAlarmV2_withEventDimensions(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("ces-oca-%s", acctest.RandString(5))
	rName := "opentelekomcloud_ces_one_click_alarm_v2.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getOneClickAlarmV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESOneClickAlarmV2WithEventDimensions(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "one_click_alarm_id", "ECSSystemOneClickAlarm"),
					resource.TestCheckResourceAttr(rName, "dimension_names.0.metric.0", "instance_id"),
					resource.TestCheckResourceAttr(rName, "dimension_names.0.event", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_notifications.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "internal_alarm_id"),
					resource.TestCheckResourceAttrSet(rName, "namespace"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"one_click_alarm_id", "dimension_names", "notification_enabled",
					"alarm_notifications", "ok_notifications", "notification_begin_time", "notification_end_time",
				},
			},
		},
	})
}

func testCESOneClickAlarmV2TopicBase(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_topic_v2" "test" {
  name         = "%s"
  display_name = "SMN topic for one-click alarm test"
}
`, name)
}

func testCESOneClickAlarmV2Basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  one_click_alarm_id = "OBSSystemOneClickAlarm"

  dimension_names {
    metric = ["bucket_name"]
  }

  notification_enabled = false
}
`, testCESOneClickAlarmV2TopicBase(name))
}

func testCESOneClickAlarmV2Update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  one_click_alarm_id = "OBSSystemOneClickAlarm"

  dimension_names {
    metric = ["bucket_name"]
  }

  notification_enabled = true

  alarm_notifications {
    type = "notification"

    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }

  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
`, testCESOneClickAlarmV2TopicBase(name))
}

func testCESOneClickAlarmV2WithNotifications(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  one_click_alarm_id = "OBSSystemOneClickAlarm"

  dimension_names {
    metric = ["bucket_name"]
  }

  notification_enabled = true

  alarm_notifications {
    type = "notification"

    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }

  ok_notifications {
    type = "notification"

    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }

  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
`, testCESOneClickAlarmV2TopicBase(name))
}

func testCESOneClickAlarmV2WithEventDimensions(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  one_click_alarm_id = "ECSSystemOneClickAlarm"

  dimension_names {
    metric = ["instance_id"]
    event  = true
  }

  notification_enabled = true

  alarm_notifications {
    type = "notification"

    notification_list = [
      opentelekomcloud_smn_topic_v2.test.topic_urn
    ]
  }

  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
`, testCESOneClickAlarmV2TopicBase(name))
}
