package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceCesAlarmRulesV2_basic(t *testing.T) {
	var (
		dataSource = "data.opentelekomcloud_ces_alarm_rules_v2.test"
		dc         = common.InitDataSourceCheck(dataSource)
		name       = fmt.Sprintf("ces-rule-%s", acctest.RandString(5))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCesAlarmRulesV2Basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.0.metric_name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.0.period"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.0.filter"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.0.comparison_operator"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.0.count"),

					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_filter_by_id_useful", "true"),
					resource.TestCheckOutput("is_filter_by_name_useful", "true"),
					resource.TestCheckOutput("is_filter_by_namespace_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCesAlarmRulesV2Basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name      = "%[1]s"
  namespace = "SYS.ECS"
  type      = "EVENT.SYS"

  resources {
    dimensions {
      name  = "resource_id"
      value = "all_instance"
    }
  }

  policies {
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    suppress_duration   = 0
    level               = 2
  }

  notification_enabled = false
  alarm_enabled        = true
}

locals {
  id        = opentelekomcloud_ces_alarm_rule_v2.test.alarm_id
  name      = opentelekomcloud_ces_alarm_rule_v2.test.name
  namespace = opentelekomcloud_ces_alarm_rule_v2.test.namespace
}

data "opentelekomcloud_ces_alarm_rules_v2" "test" {
  depends_on = [opentelekomcloud_ces_alarm_rule_v2.test]
}

output "is_default_filter_useful" {
  value = length(data.opentelekomcloud_ces_alarm_rules_v2.test.alarms) > 0
}

data "opentelekomcloud_ces_alarm_rules_v2" "filter_by_id" {
  alarm_id = local.id

  depends_on = [opentelekomcloud_ces_alarm_rule_v2.test]
}

output "is_filter_by_id_useful" {
  value = length(data.opentelekomcloud_ces_alarm_rules_v2.filter_by_id.alarms) > 0 && alltrue(
    [for alarm in data.opentelekomcloud_ces_alarm_rules_v2.filter_by_id.alarms[*] : alarm.alarm_id == local.id]
  )
}

data "opentelekomcloud_ces_alarm_rules_v2" "filter_by_name" {
  name = local.name

  depends_on = [opentelekomcloud_ces_alarm_rule_v2.test]
}

output "is_filter_by_name_useful" {
  value = length(data.opentelekomcloud_ces_alarm_rules_v2.filter_by_name.alarms) > 0 && alltrue(
    [for alarm in data.opentelekomcloud_ces_alarm_rules_v2.filter_by_name.alarms[*] : alarm.name == local.name]
  )
}

data "opentelekomcloud_ces_alarm_rules_v2" "filter_by_namespace" {
  namespace = local.namespace

  depends_on = [opentelekomcloud_ces_alarm_rule_v2.test]
}

output "is_filter_by_namespace_useful" {
  value = length(data.opentelekomcloud_ces_alarm_rules_v2.filter_by_namespace.alarms) > 0 && alltrue(
    [for alarm in data.opentelekomcloud_ces_alarm_rules_v2.filter_by_namespace.alarms[*] : alarm.namespace == local.namespace]
  )
}
`, name)
}
