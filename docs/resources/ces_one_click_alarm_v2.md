---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_one_click_alarm_v2"
sidebar_current: "docs-opentelekomcloud-resource-ces-one-click-alarm-v2"
description: |-
  Manages a CES One-Click Alarm v2 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES one-click monitoring you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_v2/one-click_monitoring/index.html)

# opentelekomcloud_ces_one_click_alarm_v2

Manages a CES One-Click Alarm v2 resource within OpenTelekomCloud.

## Example Usage

### Basic one-click monitoring without notifications

```hcl
resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  one_click_alarm_id = "ECSSystemOneClickAlarm"

  dimension_names {
    metric = ["instance_id"]
  }

  notification_enabled = false
}
```

### One-click monitoring with notifications

```hcl
variable "topic_urn" {}

resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  one_click_alarm_id = "OBSSystemOneClickAlarm"

  dimension_names {
    metric = ["bucket_name"]
  }

  notification_enabled = true

  alarm_notifications {
    type              = "notification"
    notification_list = [var.topic_urn]
  }

  ok_notifications {
    type              = "notification"
    notification_list = [var.topic_urn]
  }

  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
```

### One-click monitoring with event dimensions

```hcl
variable "topic_urn" {}

resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  one_click_alarm_id = "ECSSystemOneClickAlarm"

  dimension_names {
    metric = ["instance_id"]
    event  = true
  }

  notification_enabled = true

  alarm_notifications {
    type              = "notification"
    notification_list = [var.topic_urn]
  }

  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
```

## Argument Reference

The following arguments are supported:

* `one_click_alarm_id` - (Required, String, ForceNew) Specifies the one-click monitoring ID.
  The value can be a string of `1` to `64` alphanumeric characters.
  Changing this creates a new resource.

* `dimension_names` - (Required, List, ForceNew) Specifies the dimension names for metric and event alarm rules.
  The [dimension_names](#dimension_names_struct) structure is documented below.
  Changing this creates a new resource.

* `notification_enabled` - (Required, Bool) Specifies whether to enable alarm notifications.

* `alarm_notifications` - (Optional, List) Specifies the action to be triggered by an alarm.
  + If the value of `notification_enabled` is **false**, this parameter should not be set.
  + If the value of `notification_enabled` is **true**, this parameter is required.

  The [alarm_notifications](#notifications_struct) structure is documented below.

* `ok_notifications` - (Optional, List) Specifies the action to be triggered after an alarm is cleared.
  + If the value of `notification_enabled` is **false**, this parameter should not be set.

  The [ok_notifications](#notifications_struct) structure is documented below.

* `notification_begin_time` - (Optional, String) Specifies the alarm notification start time,
  for example: **00:00**. The format is `HH:MM`.

* `notification_end_time` - (Optional, String) Specifies the alarm notification end time,
  for example: **23:59**. The format is `HH:MM`.

<a name="dimension_names_struct"></a>
The `dimension_names` block supports:

* `metric` - (Optional, List) Specifies the dimension strings for metric alarm rules.
  A maximum of `100` items are supported. Each element can be a string of `1` to `131` characters
  that must start with a letter and contain only letters, digits, underscores (_), hyphens (-), and commas (,).

* `event` - (Optional, Bool) Specifies whether to enable event alarm rules.
  Defaults to `false`.

-> At least one of `metric` or `event` must be configured within `dimension_names`.

<a name="notifications_struct"></a>
The `alarm_notifications` and `ok_notifications` blocks support:

* `type` - (Required, String) Specifies the notification type. The value is **notification**,
  which indicates that a notification will be sent via SMN topic subscriptions.

* `notification_list` - (Required, List) Specifies the list of SMN topic URNs.
  A maximum of `20` items are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The one-click alarm ID.

* `internal_alarm_id` - The auto-generated internal alarm ID returned by the API.

* `namespace` - The namespace of the monitored service.

* `description` - The description of the one-click monitoring configuration.

* `enabled` - Whether one-click monitoring is enabled.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

CES one-click alarms v2 can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_ces_one_click_alarm_v2.test OBSSystemOneClickAlarm
```

Note that the imported state may not be identical to your resource definition, due to the API response not including
some attributes. The missing attributes include: `one_click_alarm_id`, `dimension_names`, `notification_enabled`,
`alarm_notifications`, `ok_notifications`, `notification_begin_time`, and `notification_end_time`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide
if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "opentelekomcloud_ces_one_click_alarm_v2" "test" {
  lifecycle {
    ignore_changes = [
      one_click_alarm_id, dimension_names, notification_enabled,
      alarm_notifications, ok_notifications, notification_begin_time, notification_end_time,
    ]
  }
}
```
