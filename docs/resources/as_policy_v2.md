---
subcategory: "Autoscaling"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_as_policy_v2"
sidebar_current: "docs-opentelekomcloud-resource-as-policy-v2"
description: |-
  Manages a AS Policy v2 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for AS policy you can get at
[documentation portal](https://docs.otc.t-systems.com/auto-scaling/api-ref/apis/as_policies)

# opentelekomcloud_as_policy_v2

Manages a V2 AS Policy resource within OpenTelekomCloud.

## Example Usage

### AS Recurrence Policy

```hcl
resource "opentelekomcloud_as_policy_v2" "policy_1" {
  scaling_policy_name   = "policy_create"
  scaling_policy_type   = "RECURRENCE"
  scaling_resource_id   = var.as_group
  scaling_resource_type = "SCALING_GROUP"

  scaling_policy_action {
    operation  = "ADD"
    percentage = 15
  }
  scheduled_policy {
    launch_time      = "10:30"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    end_time         = "2040-12-31T10:30Z"
  }
}
```

### AS Alarm Policy

```hcl
resource "opentelekomcloud_as_policy_v2" "policy_1" {
  scaling_policy_name   = "policy_create"
  scaling_policy_type   = "ALARM"
  scaling_resource_id   = var.as_group
  scaling_resource_type = "SCALING_GROUP"

  alarm_id = var.alarm_id

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }

  cool_down_time = 900
}
```

## Argument Reference

The following arguments are supported:

* `scaling_policy_name` - (Required) The name of the AS policy. The name can contain letters,
    digits, underscores(_), and hyphens(-),and cannot exceed 64 characters.

* `scaling_resource_id` - (Required) The Scaling resource ID.

* `scaling_resource_type` - (Required) Specifies the scaling resource type. Valid values are:
  * AS group: `SCALING_GROUP`
  * Bandwidth: `BANDWIDTH`

* `scaling_policy_type` - (Required) The AS policy type. The values can be:
  * `ALARM` - Indicates that the scaling action is triggered by an alarm. A value is returned for
    `alarm_id`, and no value is returned for `scheduled_policy`.
  * `SCHEDULED` - Indicates that the scaling action is triggered as scheduled.
    A value is returned for `scheduled_policy`, and no value is returned for `alarm_id`,
    `recurrence_type`, `recurrence_value`, `start_time`, or `end_time`.
  * `RECURRENCE` - Indicates that the scaling action is triggered periodically.
    Values are returned for `scheduled_policy`, `recurrence_type`, `recurrence_value`,
    `start_time`, and `end_time`, and no value is returned for `alarm_id`.

* `alarm_id` - (Optional) Specifies the alarm rule ID. This parameter is mandatory
  when `scaling_policy_type` is set to `ALARM`.

* `scheduled_policy` - (Optional) Specifies the periodic or scheduled AS policy.
  This parameter is mandatory when `scaling_policy_type` is set to `SCHEDULED` or `RECURRENCE`.
  After this parameter is specified, the value of `alarm_id` does not take effect.
  The `scheduled_policy` structure is documented below.

* `scaling_policy_action` - (Optional) The action of the AS policy. The `scaling_policy_action`
    structure is documented below.

* `cool_down_time` - (Optional) Specifies the cooldown period (in seconds).

The `scheduled_policy` block supports:

* `launch_time` - (Required) The time when the scaling action is triggered. If `scaling_policy_type`
  is set to `SCHEDULED`, the time format is `YYYY-MM-DDThh:mmZ`. If `scaling_policy_type` is set to
  `RECURRENCE`, the time format is `hh:mm`.

* `recurrence_type` - (Optional) The periodic triggering type. This argument is mandatory when
  `scaling_policy_type` is set to `RECURRENCE`. The options include `Daily`, `Weekly`, and `Monthly`.

* `recurrence_value` - (Optional) The frequency at which scaling actions are triggered.

-> When `recurrence_type` is set to `Daily`, this parameter does not take effect.

* `start_time` - (Optional) The start time of the scaling action triggered periodically.
  The time format complies with UTC. The current time is used by default. The time
  format is `YYYY-MM-DDThh:mmZ`.

* `end_time` - (Optional) The end time of the scaling action triggered periodically.
  The time format complies with UTC. This argument is mandatory when `scaling_policy_type`
  is set to `RECURRENCE`. The time format is `YYYY-MM-DDThh:mmZ`.

The `scaling_policy_action` block supports:

* `operation` - (Optional) The operation to be performed.

  If `scaling_resource_type` is set to `SCALING_GROUP`, the following operations are supported:
  * `ADD`: indicates adding instances.
  * `REMOVE`/`REDUCE`: indicates removing or reducing instances.
  * `SET`: indicates setting the number of instances to a specified value.

  If `scaling_resource_type` is set to `BANDWIDTH`, the following operations are supported:
  * `ADD`: indicates adding instances.
  * `REDUCE`: indicates reducing instances.
  * `SET`: indicates setting the number of instances to a specified value.

* `size` - (Optional) Specifies the operation size. The value is an integer from `0` to `300`.
  The default value is `1`. This parameter can be set to `0` only when operation is set to `SET`.
  * If `scaling_resource_type` is set to `SCALING_GROUP`, this parameter indicates the number
    of instances. The value is an integer from `0` to `300` and the default value is `1`.
  * If `scaling_resource_type` is set to `BANDWIDTH`, this parameter indicates the bandwidth
    (Mbit/s). The value is an integer from `1` to `300` and the default value is `1`.
  * If `scaling_resource_type` is set to `SCALING_GROUP`, either `size` or `percentage` can be set.

* `percentage` - (Optional) Specifies the percentage of instances to be operated.
  If operation is set to `ADD`, `REMOVE`, or `REDUCE`, the value of this parameter
  is an integer from `1` to `20000`.
  * If operation is set to `SET`, the value is an integer from `0` to `20000`.
  * If `scaling_resource_type` is set to `SCALING_GROUP`, either `size` or `percentage` can be set.
  * If neither `size` nor `percentage` is set, the default value of `size` is `1`.
  * If `scaling_resource_type` is set to `BANDWIDTH`, `percentage` is unavailable.

* `limits` - (Optional) Specifies the operation restrictions.
  * If `scaling_resource_type` is set to `BANDWIDTH` and operation is not `SET`,
  this parameter takes effect and the unit is `Mbit/s`.
  * If operation is set to `ADD`, this parameter indicates the maximum bandwidth allowed.
  * If operation is set to `REDUCE`, this parameter indicates the minimum bandwidth allowed.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Specifies the time when an AS policy was created. The time format complies with UTC.

* `metadata` - Provides additional information. The `metadata` structure is documented below.

The `metadata` block supports:

* `bandwidth_share_type` - Specifies the bandwidth sharing type in the bandwidth scaling policy.

* `eip_id` - Specifies the EIP ID for the bandwidth in the bandwidth scaling policy.

* `eip_address` - Specifies the EIP for the bandwidth in the bandwidth scaling policy.

