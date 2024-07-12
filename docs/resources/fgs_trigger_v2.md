---
subcategory: "FunctionGraph"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fgs_trigger_v2"
sidebar_current: "docs-opentelekomcloud-resource-fgs-trigger-v2"
description: |-
  Manages an FGS Trigger resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for FGS you can get at
[documentation portal](https://docs.otc.t-systems.com/function-graph/api-ref/apis/index.html)

# opentelekomcloud_fgs_trigger_v2

Manages a V2 function graph trigger resource within OpenTelekomCloud.

## Example Usage

### Create the Timing Triggers with rate and cron schedule types

```hcl
variable "function_urn" {}
variable "trigger_name" {}

// Timing trigger (with rate schedule type)
resource "opentelekomcloud_fgs_trigger_v2" "test" {
  function_urn = var.function_urn
  type         = "TIMER"
  event_data = jsonencode({
    "name" : format("%s_rate", var.trigger_name),
    "schedule_type" : "Rate",
    "user_event" : "Created by terraform script",
    "schedule" : "3m"
  })
}

// Timing trigger (with cron schedule type)
resource "opentelekomcloud_fgs_trigger_v2" "timer_cron" {
  function_urn = var.function_urn
  type         = "TIMER"
  event_data = jsonencode({
    "name" : format("%s_cron", var.trigger_name),
    "schedule_type" : "Cron",
    "user_event" : "Created by terraform script",
    "schedule" : "@every 1h30m"
  })
}
```

## Argument Reference

The following arguments are supported:

* `function_urn` - (Required, String, ForceNew) Specifies the function URN to which the function trigger belongs.

* `type` - (Required, String, ForceNew) Specifies the type of the function trigger.
  The valid values are **TIMER**, **APIG**, **CTS**, **DDS**, **DEDICATEDGATEWAY**, etc.

  -> For more available values, please refer to the [documentation table 3](https://docs.otc.t-systems.com/function-graph/api-ref/apis/function_triggers/creating_a_trigger.html#functiongraph-06-0122).

* `event_data` - (Required, String) Specifies the detailed configuration of the function trigger event.
  For various types of trigger parameter configurations, please refer to the
  [documentation](https://docs.otc.t-systems.com/function-graph/api-ref/apis/function_triggers/creating_a_trigger.html#id4).

  -> Please refer to the [documentation](https://docs.otc.t-systems.com/function-graph/api-ref/apis/function_triggers/updating_a_trigger.html#functiongraph-06-0124-request-updateriggereventdata)
     for updatable fields.

* `status` - (Optional, String) Specifies the status of the function trigger.
  The valid values are **ACTIVE** and **DISABLED**.
  For `DDS` and `Kafka` triggers the default value is **DISABLED**, for other triggers= the default value is **ACTIVE**.

  -> Currently, only some triggers support setting the **DISABLED** value, such as `TIMER`, `DDS`, `DMS`, `KAFKA` and
     `LTS`. For more details, please refer to the [documentation](https://docs.otc.t-systems.com/function-graph/api-ref/apis/function_triggers/creating_a_trigger.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - resource ID in UUID format.

* `created_at` - The creation time of the function trigger.

* `updated_at` - The latest update time of the function trigger.

* `region` - The region where the function trigger is located.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

Function trigger can be imported using the `function_urn`, `type` and `id`, separated by the slashes (/), e.g.

```bash
$ terraform import opentelekomcloud_fgs_trigger_v2.test <function_urn>/<type>/<id>
```
