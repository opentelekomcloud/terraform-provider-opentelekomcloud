---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_alarm_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-ces-alarm-rule-v2"
description: |-
  Manages a CES Alarm Rule v2 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES alarm rule you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_v2/alarm_rules/index.html)

# opentelekomcloud_ces_alarm_rule_v2

Manages a CES Alarm Rule v2 resource within OpenTelekomCloud.

~>
  Alarm rule `namespaces` and `dimensions` are available on our [github link](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/tree/devel/opentelekomcloud/services/ces/interconnected_services.md) or [official documentation](https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html).

## Example Usage

### Basic alarm rule for multiple ECS instances

```hcl
variable "instance_id_1" {}
variable "instance_id_2" {}
variable "topic_urn" {}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name                 = "alarm-rule-test"
  namespace            = "SYS.ECS"
  type                 = "MULTI_INSTANCE"
  notification_enabled = true
  alarm_enabled        = true

  resources {
    dimensions {
      name  = "instance_id"
      value = var.instance_id_1
    }
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = var.instance_id_2
    }
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 1200
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    level               = 4
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 3600
    filter              = "average"
    comparison_operator = ">="
    value               = 20
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    level               = 4
  }

  alarm_actions {
    type              = "notification"
    notification_list = [var.topic_urn]
  }
}
```

### Alarm rule for all instances

```hcl
variable "topic_urn" {}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name                 = "alarm-rule-all-instance"
  namespace            = "AGT.ECS"
  type                 = "ALL_INSTANCE"
  notification_enabled = true
  alarm_enabled        = true

  resources {
    dimensions {
      name = "instance_id"
    }

    dimensions {
      name = "mount_point"
    }
  }

  policies {
    metric_name         = "disk_usedPercent"
    period              = 1
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    count               = 1
    suppress_duration   = 0
    level               = 2
  }

  alarm_actions {
    type              = "notification"
    notification_list = [var.topic_urn]
  }
}
```

### Alarm rule for system event monitoring

```hcl
resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name                 = "alarm-rule-sys-event"
  namespace            = "SYS.ECS"
  type                 = "EVENT.SYS"
  notification_enabled = false
  alarm_enabled        = true

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
}
```

### Alarm rule for system event monitoring with notification

```hcl
variable "topic_urn" {}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name                 = "alarm-rule-sys-event"
  namespace            = "SYS.ECS"
  type                 = "EVENT.SYS"
  notification_enabled = true
  alarm_enabled        = true

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

  alarm_actions {
    type              = "notification"
    notification_list = [var.topic_urn]
  }
}
```

### Alarm rule with OK actions

```hcl
variable "instance_id" {}
variable "topic_urn" {}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name                 = "alarm-rule-with-ok-actions"
  namespace            = "SYS.ECS"
  type                 = "MULTI_INSTANCE"
  notification_enabled = true
  alarm_enabled        = true

  resources {
    dimensions {
      name  = "instance_id"
      value = var.instance_id
    }
  }

  policies {
    metric_name         = "network_outgoing_bytes_rate_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    level               = 2
  }

  alarm_actions {
    type              = "notification"
    notification_list = [var.topic_urn]
  }

  ok_actions {
    type              = "notification"
    notification_list = [var.topic_urn]
  }
}
```

### Alarm rule using alarm template

```hcl
variable "instance_id" {}
variable "topic_urn" {}
variable "alarm_template_id" {}

resource "opentelekomcloud_ces_alarm_rule_v2" "test" {
  name                 = "alarm-rule-with-template"
  namespace            = "SYS.ECS"
  type                 = "MULTI_INSTANCE"
  alarm_template_id    = var.alarm_template_id
  notification_enabled = true
  alarm_enabled        = true

  resources {
    dimensions {
      name  = "instance_id"
      value = var.instance_id
    }
  }

  alarm_actions {
    type              = "notification"
    notification_list = [var.topic_urn]
  }

  ok_actions {
    type              = "notification"
    notification_list = [var.topic_urn]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of an alarm rule. The value can be a string of `1` to `128`
  characters that can consist of letters, digits, underscores (_), and hyphens (-).
  Changing this creates a new resource.

* `namespace` - (Required, String, ForceNew) Specifies the namespace in `service.item` format. `service` and `item`
  each must be a string that starts with a letter and contains only letters, digits, and underscores (_).
  Changing this creates a new resource.
  For details, see [Services Interconnected with Cloud Eye](https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html).

* `policies` - (Optional, Set) Specifies the alarm policies. The [policies](#policies_struct) structure is
  documented below. Exactly one of `policies` or `alarm_template_id` must be specified.
  When using `alarm_template_id`, the policies are inherited from the template and cannot be modified directly.

* `description` - (Optional, String, ForceNew) The value can be a string of `0` to `256` characters.
  Changing this creates a new resource.

* `type` - (Optional, String, ForceNew) Specifies the alarm type. The value can be:
  + **EVENT.SYS**: The alarm rule is created for system events.
  + **EVENT.CUSTOM**: The alarm rule is created for custom events.
  + **MULTI_INSTANCE**: The alarm rule is created for multiple instances.
  + **ALL_INSTANCE**: The alarm rule is created for all instances.

  Defaults to **MULTI_INSTANCE**. Changing this creates a new resource.

* `resources` - (Optional, Set) Specifies the list of resources to monitor.
  The [resources](#resources_struct) structure is documented below.

* `resource_group_id` - (Optional, String, ForceNew) Specifies the resource group ID.
  Changing this creates a new resource.

* `alarm_template_id` - (Optional, String, ForceNew) Specifies the ID of the alarm template.
  When using an alarm template, the policies are inherited from the template.
  Exactly one of `alarm_template_id` or `policies` must be specified.
  Changing this creates a new resource.

* `alarm_actions` - (Optional, List, ForceNew) Specifies the action triggered by an alarm.
  The [alarm_actions](#actions_struct) structure is documented below.
  Changing this creates a new resource.

* `ok_actions` - (Optional, List, ForceNew) Specifies the action triggered by the clearing of an alarm.
  The [ok_actions](#actions_struct) structure is documented below.
  Changing this creates a new resource.

* `notification_enabled` - (Optional, Bool, ForceNew) Specifies whether to enable the action to be triggered by an alarm.
  The default value is `false`. Changing this creates a new resource.

* `alarm_enabled` - (Optional, Bool) Specifies whether to enable the alarm rule.
  The default value is `true`.

* `notification_begin_time` - (Optional, String, ForceNew) Specifies the alarm notification start time,
  for example: **05:30**. Changing this creates a new resource.

* `notification_end_time` - (Optional, String, ForceNew) Specifies the alarm notification stop time,
  for example: **22:10**. Changing this creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the alarm rule.
  Changing this creates a new resource.

-> **Note** If `notification_enabled` is set to `true`, either `alarm_actions` or `ok_actions` cannot be empty.
If `alarm_actions` and `ok_actions` coexist, their corresponding `notification_list` must be of the same value.

<a name="policies_struct"></a>
The `policies` block supports:

* `metric_name` - (Required, String) Specifies the metric name. The value can be a string of `1` to `64` characters
  that must start with a letter and contain only letters, digits, and underscores (_).
  For details, see [Services Interconnected with Cloud Eye](https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html).

* `period` - (Required, Int) Specifies the alarm checking period in seconds. The value can be `0`, `1`, `300`,
  `1200`, `3600`, `14400`, and `86400`.

  -> If `period` is set to `1`, the raw metric data is used to determine whether to generate an alarm.
  When the value of `type` is **EVENT.SYS** or **EVENT.CUSTOM**, `period` can be set to `0`.

* `filter` - (Required, String) Specifies the data rollup method. The value can be `max`, `min`, `average`,
  `sum`, and `variance`.

* `comparison_operator` - (Required, String) Specifies the comparison condition of alarm thresholds.
  The value can be `>`, `=`, `<`, `>=`, `<=`, or `!=`.

* `value` - (Optional, Float) Specifies the alarm threshold. The value ranges from `0` to
  `Number.MAX_VALUE` (1.7976931348623157e+108).

* `count` - (Required, Int) Specifies the number of consecutive occurrence times.
  The value ranges from `1` to `5`.

* `unit` - (Optional, String) Specifies the data unit. The value can be a string of `0` to `32` characters.

* `suppress_duration` - (Optional, Int) Specifies the interval for triggering an alarm if the alarm persists.
  Possible values are as follows:
  + **0**: Cloud Eye triggers the alarm only once.
  + **300**: Cloud Eye triggers the alarm every 5 minutes.
  + **600**: Cloud Eye triggers the alarm every 10 minutes.
  + **900**: Cloud Eye triggers the alarm every 15 minutes.
  + **1800**: Cloud Eye triggers the alarm every 30 minutes.
  + **3600**: Cloud Eye triggers the alarm every hour.
  + **10800**: Cloud Eye triggers the alarm every 3 hours.
  + **21600**: Cloud Eye triggers the alarm every 6 hours.
  + **43200**: Cloud Eye triggers the alarm every 12 hours.
  + **86400**: Cloud Eye triggers the alarm every day.

  The default value is `0`.

* `level` - (Optional, Int) Specifies the alarm severity. The value can be `1`, `2`, `3`, or `4`, which indicates
  *critical*, *major*, *minor*, and *informational*, respectively.
  The default value is `2`.

<a name="resources_struct"></a>
The `resources` block supports:

* `dimensions` - (Optional, List) Specifies the list of metric dimensions.
  The [dimensions](#dimensions_struct) structure is documented below.
  This block can be omitted when the alarm scope applies to all resources, for example
  project-wide `EVENT.SYS` or `EVENT.CUSTOM` alarms.

<a name="dimensions_struct"></a>
The `dimensions` block supports:

* `name` - (Required, String) Specifies the dimension name. The value can be a string of `1` to `32` characters
  that must start with a letter and contain only letters, digits, underscores (_), and hyphens (-).

* `value` - (Optional, String) Specifies the dimension value. The value can be a string of `1` to `64` characters
  that must start with a letter or a number and contain only letters, digits, underscores (_), and hyphens (-).
  This field can be left empty when the alarm scope applies to all resources.

<a name="actions_struct"></a>
The `alarm_actions` block supports:

* `type` - (Required, String) Specifies the type of action triggered by an alarm. The value can be:
  + **notification**: indicates that a notification will be sent to the user.
  + **autoscaling**: indicates that a scaling action will be triggered.

* `notification_list` - (Required, List) Specifies the list of objects to be notified if the alarm status changes.
  The maximum length is `5`. If `type` is set to **notification**, the value of `notification_list` cannot be empty.
  If `type` is set to **autoscaling**, the value of `notification_list` must be `[]` and the value of `namespace`
  must be **SYS.AS**.

  -> To enable the autoscaling alarm rules take effect, you must bind scaling policies.

The `ok_actions` block supports:

* `type` - (Required, String) Specifies the type of action triggered by the clearing of an alarm.
  The value can be:
  + **notification**: indicates that a notification will be sent to the user.
  + **autoscaling**: indicates that a scaling action will be triggered.

* `notification_list` - (Required, List) Specifies the list of objects to be notified if the alarm status changes.
  The maximum length is `5`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `alarm_id` - The alarm rule ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

CES alarm rules v2 can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_ces_alarm_rule_v2.alarm_rule al1619578509719Ga0X1RGWv
```
