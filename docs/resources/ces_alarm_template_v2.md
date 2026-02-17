---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_alarm_template_v2"
sidebar_current: "docs-opentelekomcloud-resource-ces-alarm-template-v2"
description: |-
  Manages a CES Alarm Template resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES alarm template you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_v2/alarm_templates/index.html)

# opentelekomcloud_ces_alarm_template_v2

Manages a CES alarm template resource within OpenTelekomCloud.

## Example Usage

### Create a metric alarm template

```hcl
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "my-alarm-template"
  description = "Alarm template for ECS monitoring"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">="
    value               = 80
    unit                = "%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 3600
  }
}
```

### Create an event alarm template

```hcl
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "my-event-template"
  type        = 2
  description = "Event alarm template"

  policies {
    namespace           = "SYS.ECS"
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    alarm_level         = 2
    suppress_duration   = 0
  }
}
```

### Create an alarm template with multiple policies

```hcl
resource "opentelekomcloud_ces_alarm_template_v2" "test" {
  name        = "multi-policy-template"
  description = "Alarm template with multiple monitoring policies"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "mem_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 85
    unit                = "%"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "disk_util_inband"
    period              = 300
    filter              = "average"
    comparison_operator = ">="
    value               = 90
    unit                = "%"
    count               = 3
    alarm_level         = 1
    suppress_duration   = 600
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the CES alarm template.
  An alarm template name starts with a letter, consists of `1` to `128` characters,
  and can contain only letters, digits, underscores (_), hyphens (-), parentheses, and periods (.).

* `policies` - (Required, Set) Specifies the policy list of the CES alarm template.
  A maximum of `50` policies are supported.
  The [policies](#CesAlarmTemplate_Policy) structure is documented below.

* `type` - (Optional, Int, ForceNew) Specifies the type of the CES alarm template.
  Defaults to `0`. Changing this parameter will create a new resource.
  The valid values are as follows:
  + **0**: metric alarm template.
  + **2**: event alarm template.

* `description` - (Optional, String) Specifies the description of the CES alarm template.
  The description can contain a maximum of `256` characters.

* `delete_associate_alarm` - (Optional, Bool) Specifies whether to delete the alarm rules
  associated with the alarm template when deleting the template. Defaults to **false**.

<a name="CesAlarmTemplate_Policy"></a>
The `policies` block supports:

* `namespace` - (Required, String) Specifies the namespace of the service.
  The value must be in the `service.item` format and can contain `3` to `32` characters.
  For details, see [Services Interconnected with Cloud Eye](https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html).

* `metric_name` - (Required, String) Specifies the alarm metric name.
  The value must start with a letter and can contain `1` to `96` characters,
  including letters, digits, and underscores (_).

* `period` - (Required, Int) Specifies the aggregation period of alarm condition in seconds.
  Value options: **0**, **1**, **300**, **1200**, **3600**, **14400**, **86400**.

  -> If `period` is set to `1`, the raw metric data is used to determine whether to generate an alarm.
  When the value of `type` is **2** (event alarm template), `period` can be set to `0`.

* `filter` - (Required, String) Specifies the data rollup methods.
  Value options: **average**, **variance**, **min**, **max**, **sum**.

* `comparison_operator` - (Required, String) Specifies the comparison conditions for alarm threshold.
  Value options: **>**, **<**, **=**, **>=**, **<=**, **!=**.

* `count` - (Required, Int) Specifies the number of consecutive alarm triggering times.
  + For event alarms, the value ranges from **1** to **180**.
  + For metric alarms, the value can be **1**, **2**, **3**, **4**, **5**, **10**, **15**, **30**, **60**,
    **90**, **120**, **180**.

* `suppress_duration` - (Required, Int) Specifies the alarm suppression cycle in seconds.
  Only one alarm is sent when the alarm suppression period is **0**.
  Value options: **0**, **300**, **600**, **900**, **1800**, **3600**, **10800**, **21600**,
  **43200**, **86400**.

* `value` - (Optional, Float) Specifies the alarm threshold.
  The value ranges from `0` to `Number.MAX_VALUE` (1.7976931348623157e+108).

* `alarm_level` - (Optional, Int) Specifies the alarm level.
  Defaults to `2`. The valid values are as follows:
  + **1**: critical.
  + **2**: major.
  + **3**: minor.
  + **4**: warning.

* `unit` - (Optional, String) Specifies the unit string of the alarm threshold.

* `dimension_name` - (Optional, String) Specifies the resource dimension.
  Leave this parameter blank for an event alarm template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `template_id`).

* `template_id` - The alarm template ID.

## Import

The CES alarm template can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_ces_alarm_template_v2.test <template_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `delete_associate_alarm`.
It is generally recommended running `terraform plan` after importing an alarm template.
You can then decide if changes should be applied to the alarm template, or the resource definition should be updated to
align with the alarm template. Also you can ignore changes as below.

```hcl
resource "opentelekomcloud_ces_alarm_template_v2" "test" {

  lifecycle {
    ignore_changes = [
      delete_associate_alarm,
    ]
  }
}
```
