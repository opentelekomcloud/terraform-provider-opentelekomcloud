---
subcategory: "Cloud Eye (CES)"
---

Up-to-date reference of API arguments for CES alarm rule you can get at
`https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/alarm_rule_managements`.

# opentelekomcloud_ces_alarmrule

Manages a V1 CES Alarm Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
variable server_id {}
variable smn_topic_id {}

resource "opentelekomcloud_ces_alarmrule" "alarm_rule" {
  alarm_name = "alarm_rule"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = var.server_id
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

  alarm_actions {
    type              = "notification"
    notification_list = [var.smn_topic_id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_name` - (Required) Specifies the name of an alarm rule. The value can
  be a string of `1` to `128` characters that can consist of numbers, lowercase letters,
  uppercase letters, underscores (_), or hyphens (-).

* `alarm_type` - (Optional) Specifies the alarm rule type.
  `EVENT.SYS`: The alarm rule is created for system events.
  `EVENT.CUSTOM`: The alarm rule is created for custom events.

* `alarm_description` - (Optional) Alarm description. The value can be a string of `0` to `256` characters.

* `alarm_level` - (Optional) Specifies the alarm severity. The value can be `1`, `2`, `3` or `4`,
  which indicates `critical`, `major`, `minor`, and `informational`. The default value is `2`.

* `metric` - (Required) Specifies the alarm metrics. The structure is described below.

* `condition` - (Required) Specifies the alarm triggering condition. The structure
  is described below.

* `alarm_actions` - (Optional) Specifies the actions list triggered by an alarm. The
  structure is described below.

* `ok_actions` - (Optional) Specifies the actions list triggered by the clearing of
  an alarm. The structure is described below.

* `alarm_enabled` - (Optional) Specifies whether to enable the alarm. The default
  value is `true`.

* `alarm_action_enabled` - (Optional) Specifies whether to enable the action
  to be triggered by an alarm. The default value is `true`.

-> If `alarm_action_enabled` is set to `true`, at least one of the following
  parameters `alarm_actions` or `ok_actions` cannot be empty.
  If `alarm_actions` and `ok_actions` coexist, their corresponding
  `notification_list` must be of the same value.

The `metric` block supports:

* `namespace` - (Required) Specifies the namespace in `service.item` format. `service.item`
  can be a string of `3` to `32` characters that must start with a letter and can
  consists of uppercase letters, lowercase letters, numbers, or underscores (_).

* `metric_name` - (Required) Specifies the metric name. The value can be a string
  of `1` to `64` characters that must start with a letter and can consists of uppercase
  letters, lowercase letters, numbers, underscores (_) or slashes (/).

* `dimensions` - (Required) Specifies the list of metric dimensions. Currently,
  the maximum length of the dimension list that are supported is `3`. The structure
  is described below.

The `dimensions` block supports:

* `name` - (Required) Specifies the dimension name. The value can be a string
  of `1` to `32` characters that must start with a letter and can consists of uppercase
  letters, lowercase letters, numbers, underscores (_), or hyphens (-).

* `value` - (Required) Specifies the dimension value. The value can be a string
  of `1` to `64` characters that must start with a letter or a number and can consists
  of uppercase letters, lowercase letters, numbers, underscores (_), or hyphens (-).

The `condition` block supports:

* `period` - (Required) Specifies the alarm checking period in seconds. The
  value can be `1`, `300`, `1200`, `3600`, `14400`, and `86400`.

-> If `period` is set to `1`, the raw metric data is used to determine
  whether to generate an alarm.

* `filter` - (Required) Specifies the data rollup methods. The value can be
  `max`, `min`, `average`, `sum`, and `variance`.

* `comparison_operator` - (Required) Specifies the comparison condition of alarm
  thresholds. The value can be `>`, `=`, `<`, `>=`, or `<=`.

* `value` - (Required) Specifies the alarm threshold. The value ranges from
  `0` to `Number.MAX_VALUE` of `1.7976931348623157e+108`.

* `unit` - (Optional) Specifies the data unit.

* `count` - (Required) Specifies the number of consecutive occurrence times.
  The value ranges from `1` to `5`.

* `alarm_frequency` - (Optional) Specifies frequency for alarm triggering. If argument is not provided alarm will be triggered once.
  `300`: Cloud Eye triggers the alarm every 5 minutes.
  `600`: Cloud Eye triggers the alarm every 10 minutes.
  `900`: Cloud Eye triggers the alarm every 15 minutes.
  `1800`: Cloud Eye triggers the alarm every 30 minutes.
  `3600`: Cloud Eye triggers the alarm every hour.
  `10800`: Cloud Eye triggers the alarm every 3 hours.
  `21600`: Cloud Eye triggers the alarm every 6 hours.
  `43200`: Cloud Eye triggers the alarm every 12 hours.
  `86400`: Cloud Eye triggers the alarm every day.

the `alarm_actions` block supports:

* `type` - (Optional) Specifies the type of action triggered by an alarm. The
  value can be notification or autoscaling.
  * `notification`: indicates that a notification will be sent to the user.
  * `autoscaling`: indicates that a scaling action will be triggered.

* `notification_list` - (Required) Specifies the topic urn list of the target
  notification objects. The maximum length is `5`. The topic urn list can be
  obtained from simple message notification (SMN) and in the following format:
  `urn:smn:([a-z]|[a-z]|[0-9]|\-){1,32}:([a-z]|[a-z]|[0-9]){32}:([a-z]|[a-z]|[0-9]|\-|\_){1,256}`.
  If `type` is set to `notification`, the value of `notification_list` cannot be
  empty. If `type` is set to `autoscaling`, the value of `notification_list` must
  be `[]`.

-> To enable the AS alarm rules take effect, you must bind scaling
  policies. For details, see the [AutoScaling API Reference](https://docs.otc.t-systems.com/en-us/api/as/en-us_topic_0045219159.html).

The `ok_actions` block supports:

* `type` - (Optional) specifies the type of action triggered by an alarm. the
  value is notification.
  * `notification`: indicates that a notification will be sent to the user.
  * `autoscaling`: indicates that a scaling action will be triggered.

* `notification_list` - (Optional) Indicates the list of objects to be notified
  if the alarm status changes. The maximum length is `5`.

## Attributes Reference

The following attributes are exported:

* `alarm_name` - See Argument Reference above.

* `alarm_description` - See Argument Reference above.

* `alarm_level` - See Argument Reference above.

* `metric` - See Argument Reference above.

* `condition` - See Argument Reference above.

* `alarm_actions` - See Argument Reference above.

* `ok_actions` - See Argument Reference above.

* `alarm_enabled` - See Argument Reference above.

* `alarm_action_enabled` - See Argument Reference above.

* `id` - Specifies the alarm rule ID.

* `update_time` - Specifies the time when the alarm status changed. The value
  is a UNIX timestamp and the unit is ms.

* `alarm_state` - Specifies the alarm status. The value can be:
  * `ok`: The alarm status is normal;
  * `alarm`: An alarm is generated;
  * `insufficient_data`: The required data is insufficient;

## Import

CES alarms can be imported using alarm rule `id`, e.g.

```sh
terraform import opentelekomcloud_ces_alarmrule.alarmrule c1881895-cdcb-4d23-96cb-032e6a3ee667
```
