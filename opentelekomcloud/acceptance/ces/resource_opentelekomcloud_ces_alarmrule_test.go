package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/alarms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	ecs "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceAlarmRuleName = "opentelekomcloud_ces_alarmrule.alarmrule_1"

func TestCESAlarmRule_basic(t *testing.T) {
	var ar alarms.MetricAlarms

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testCESAlarmRuleExists(resourceAlarmRuleName, &ar),
				),
			},
			{
				Config: testCESAlarmRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceAlarmRuleName, "alarm_enabled", "false"),
				),
			},
		},
	})
}

func TestCESAlarmRule_systemEvents(t *testing.T) {
	var ar alarms.MetricAlarms

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleSystemEvents,
				Check: resource.ComposeTestCheckFunc(
					testCESAlarmRuleExists(resourceAlarmRuleName, &ar),
				),
			},
		},
	})
}

func TestAccCESAlarmRules_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleBasic,
			},

			{
				ResourceName:      resourceAlarmRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCheckCESV1AlarmValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testCESAlarmRuleValidation,
				ExpectError: regexp.MustCompile("Error: either `alarm_actions` or `ok_actions` should be specified.+"),
			},
		},
	})
}

func TestCESAlarmRule_slashes(t *testing.T) {
	var ar alarms.MetricAlarms
	resourceName := "opentelekomcloud_ces_alarmrule.alarmrule_s"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := ecs.QuotasForFlavor(env.OsFlavorID)
			qts = append(qts,
				&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.Volume, Count: 1},
				&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4},
			)
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRuleSlashes,
				Check: resource.ComposeTestCheckFunc(
					testCESAlarmRuleExists(resourceName, &ar),
				),
			},
		},
	})
}

func testCESAlarmRuleDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.CesV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CESv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_ces_alarmrule" {
			continue
		}

		id := rs.Primary.ID
		_, err := alarms.ShowAlarm(client, id)
		if err == nil {
			return fmt.Errorf("alarm rule still exists")
		}
	}

	return nil
}

func testCESAlarmRuleExists(n string, ar *alarms.MetricAlarms) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CesV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud CESv1 client: %w", err)
		}

		id := rs.Primary.ID
		found, err := alarms.ShowAlarm(client, id)
		if err != nil {
			return err
		}

		*ar = found[0]

		return nil
	}
}

var testCESAlarmRuleBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "instance_1"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_ces_alarmrule" "alarmrule_1" {
  alarm_name = "alarm_rule1"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }
  alarm_action_enabled = false

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.topic_1.topic_urn
    ]
  }
}
`, common.DataSourceSubnet)

var testCESAlarmRuleUpdate = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "instance_1"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_ces_alarmrule" "alarmrule_1" {
  alarm_name = "alarm_rule1"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }
  alarm_action_enabled = false
  alarm_enabled        = false

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.topic_1.topic_urn
    ]
  }
}
`, common.DataSourceSubnet)

var testCESAlarmRuleValidation = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "instance_1"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_ces_alarmrule" "alarmrule_1" {
  alarm_name = "alarm_rule1"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }
  alarm_action_enabled = true
}
`, common.DataSourceSubnet)

var testCESAlarmRuleSlashes = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name        = "instance_1"
  image_name  = "Standard_Debian_11_latest"
  flavor_name = "s3.large.2"

  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_ces_alarmrule" "alarmrule_s" {
  alarm_name = "alarm_rule_s"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "/mnt/share_disk_usedPercent"
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.vm_1.id
    }
  }
  condition {
    period              = 1
    filter              = "sum"
    comparison_operator = ">="
    value               = 90
    unit                = "%%"
    count               = 3
  }
  alarm_action_enabled = false
}
`, common.DataSourceSubnet)

const testCESAlarmRuleSystemEvents = `
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_2"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_ces_alarmrule" "alarmrule_1" {
  alarm_name = "alarm_rule1"
  alarm_type = "EVENT.SYS"

  metric {
    namespace   = "SYS.CBR"
    metric_name = "backupFailed"
    dimensions {
      name  = "instance_id"
      value = "test-id"
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
    alarm_frequency     = 300
  }
  alarm_action_enabled = false

  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.topic_1.topic_urn
    ]
  }
}
`
