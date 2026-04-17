---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_event_report_v1"
sidebar_current: "docs-opentelekomcloud-resource-ces-event-report-v1"
description: |-
  Manages a v1 CES Event Report resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES event report v1 you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/event_monitoring/index.html)

# opentelekomcloud_ces_event_report_v1

Manages a V1 CES Event Report (used for reporting custom events) resource within OpenTelekomCloud.

## Example Usage

### Basic event report

```hcl
resource "opentelekomcloud_ces_event_report_v1" "event_report_1" {
  event_name   = "Test_tf_event"
  event_source = "SYS.ECS"
  time         = 1257894000000
  detail {
    content     = "This is a test event"
    event_state = "normal"
    event_level = "Info"
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_name` - (Required, String, ForceNew) Specifies the event name. The value must be `1` to `64` characters long and contain letters, digits, and underscores (_). It must start with a letter.

* `event_source` - (Required, String, ForceNew) Specifies the event source. The format is `service.item`. Set this parameter based on the site requirements. `service` and `item` each must be a string that starts with a letter and contains 3 to 32 characters, including only letters, digits, and underscores (_).

* `time` - (Required, Int, ForceNew) Specifies when the event occurred, which is a UNIX timestamp (ms). 
  
  -> **NOTE:**
  Since there is a latency between the client and the server, the data timestamp to be inserted should be within the period that starts from one hour before the current time plus 20s to 10 minutes after the current time minus 20s. In this way, the timestamp will be inserted to the database without being affected by the latency. For example, if the current time is 2020.01.30 12:00:30, the timestamp inserted must be within the range [2020.01.30 11:00:50, 2020.01.30 12:10:10]. The corresponding UNIX timestamp is [1580353250, 1580357410].

* `detail` - (Required, List, ForceNew) Specifies the event details. The structure is described below.

The `detail` block supports:

* `content` - (Optional, String, ForceNew) Specifies the event content. Enter up to `4,096` characters.
* `group_id` - (Optional, String, ForceNew) Specifies the resource group the event belongs to. This ID must be an existing resource group ID. This can be attained from `Cloud Eye` in management console. Resource group ID is listed under `Name/ID` column in `Resource Groups`.
* `resource_id` - (Optional, String, ForceNew) Specifies the resource ID.
* `resource_name` - (Optional, String, ForceNew) Specifies the resource name.
* `event_state` - (Optional, String, ForceNew) Specifies the event status. The value can be `normal`, `warning`, or `incident`.
* `event_level` - (Optional, String, ForceNew) Specifies the event severity. The value can be `Critical`, `Major`, `Minor`, or `Info`.
* `event_user` - (Optional, String, ForceNew) Specifies the event user.
* `event_type` - (Optional, String, ForceNew) Specifies the event type. Its value can be `EVENT.SYS` or `EVENT.CUSTOM`. `EVENT.SYS` indicates system events that cannot be reported by users. **Only custom events can be reported.**

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `event_id` - Specifies the event ID.
