---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_alarm_rules_v2"
sidebar_current: "docs-opentelekomcloud-datasource-ces-alarm-rules-v2"
description: |-
  Use this data source to get the list of CES alarm rules within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES alarm rules you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_v2/alarm_rules/index.html)

# opentelekomcloud_ces_alarm_rules_v2

Use this data source to get the list of CES alarm rules within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_ces_alarm_rules_v2" "all" {}
```

### Filter by alarm rule name

```hcl
data "opentelekomcloud_ces_alarm_rules_v2" "by_name" {
  name = "my-alarm-rule"
}
```

### Filter by namespace

```hcl
data "opentelekomcloud_ces_alarm_rules_v2" "by_namespace" {
  namespace = "SYS.ECS"
}
```

## Argument Reference

The following arguments are supported:

* `alarm_id` - (Optional, String) Specifies the alarm rule ID.

* `name` - (Optional, String) Specifies the name of an alarm rule.

* `namespace` - (Optional, String) Specifies the namespace of a service.

* `resource_id` - (Optional, String) Specifies the alarm resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarms` - The alarm rule list.

  The [alarms](#alarms_struct) structure is documented below.

<a name="alarms_struct"></a>
The `alarms` block supports:

* `alarm_id` - The alarm rule ID.

* `name` - The alarm rule name.

* `description` - The alarm rule description.

* `namespace` - The namespace of a service.

* `type` - The alarm rule type.

* `alarm_enabled` - Whether the alarm rule is enabled.

* `notification_enabled` - Whether the action to be triggered by an alarm is enabled.

* `notification_begin_time` - The time when the alarm notification was enabled.

* `notification_end_time` - The time when the alarm notification was disabled.

* `enterprise_project_id` - The enterprise project ID.

* `alarm_template_id` - The ID of an alarm template associated with an alarm rule.

* `policies` - The alarm policy list.

  The [policies](#alarms_policies_struct) structure is documented below.

* `resources` - The resource list.

  The [resources](#alarms_resources_struct) structure is documented below.

* `alarm_actions` - The action triggered by an alarm.

  The [alarm_actions](#notification_struct) structure is documented below.

* `ok_actions` - The action triggered after an alarm is cleared.

  The [ok_actions](#notification_struct) structure is documented below.

<a name="alarms_policies_struct"></a>
The `policies` block supports:

* `metric_name` - The metric name of a resource.

* `period` - The monitoring period of a metric.

* `filter` - The data rollup method.

* `comparison_operator` - The comparison condition of alarm thresholds.

* `value` - The alarm threshold.

* `unit` - The metric unit.

* `count` - The number of consecutive times that the alarm triggering conditions are met.

* `suppress_duration` - The interval for triggering an alarm if the alarm persists.

* `level` - The alarm severity.

<a name="alarms_resources_struct"></a>
The `resources` block supports:

* `dimensions` - The dimension information.

  The [dimensions](#resources_dimensions_struct) structure is documented below.

<a name="resources_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The name of the metric dimension.

* `value` - The value of the metric dimension.

<a name="notification_struct"></a>
The `alarm_actions` or `ok_actions` blocks support:

* `type` - The type of action triggered by an alarm.

* `notification_list` - The list of objects to be notified if the alarm status changes.
