---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_event_details_v1"
sidebar_current: "docs-opentelekomcloud-datasource-ces-event-details-v1"
description: |-
  Get details about a CES event within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES event details v1 you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/event_monitoring/index.html)

# opentelekomcloud_ces_event_details_v1

Get details about a CES event within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_ces_event_details_v1" "event_details_1" {
  event_name = "Test_Acc_tf_event"
  event_type = "EVENT.CUSTOM"
  from       = 1257893000000
  to         = 1257895000000
}
```

## Argument Reference

The following arguments are supported:

* `event_name` - (Required, String) Specifies the event name.

* `event_type` - (Required, String) Specifies the event type. The value can be `EVENT.SYS` (system event) or `EVENT.CUSTOM` (custom event).

* `event_source` - (Optional, String) Specifies the event name. The name can be a system event name or a custom event name.

* `event_level` - (Optional, String) Specifies the event severity. The value can be `Critical`, `Major`, `Minor`, or `Info`.

* `event_user` - (Optional, String) Specifies the name of the user who reports the event monitoring data. It can also be a project ID.

* `event_state` - (Optional, String) Specifies the event status. The value can be `normal`, `warning`, or `incident`.

* `from` - (Optional, Int) Specifies the start time of the query, which is a UNIX timestamp (ms). 

* `to` - (Optional, Int) Specifies the end time of the query, which is a UNIX timestamp (ms). 

* `limit` - (Optional, Int) Specifies the maximum number of records that can be queried at a time. Supported range: `1` to `100` (default)

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `event_users` - Specifies the name of the user who reports the event. It can also be a project ID.
* `event_sources` - Specifies the event source. For a system event, the source is the namespace of each service. 
* `event_info` - Specifies details of one or more events. The structure is described below. 
* `meta_data` - Specifies the number of metadata records in the query result. The structure is described below.
    * `total` - Specifies the total number of events.


The `event_info` block supports:

* `event_name` - Specifies the event name.
* `event_source` - Specifies the event source in the format of service.item.
* `time` - Specifies when the event occurred, which is a UNIX timestamp (ms).
* `detail` - Specifies the event details. The structure is described below.
  * `content` - Specifies the event content.
  * `group_id` - Specifies the resource group the event belongs to.
  * `resource_id` - Specifies the resource ID.
  * `resource_name` - Specifies the resource name.
  * `event_state` - Specifies the event status.
  * `event_level` - Specifies the event severity.
  * `event_user` - Specifies the event user.
  * `event_type` - Specifies the event type.
* `event_id` - Specifies the event ID.
