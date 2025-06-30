---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_events_v1"
sidebar_current: "docs-opentelekomcloud-datasource-ces-events-v1"
description: |-
  Get details about CES events within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES events v1 you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/event_monitoring/index.html)

# opentelekomcloud_ces_events_v1

Get details about CES events within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_ces_events_v1" "events_1" {
  event_type = "EVENT.CUSTOM"
  from       = 1257893000000
  to         = 1257895000000
}
```

## Argument Reference

The following arguments are supported:

* `event_name` - (Optional, String) Specifies the event name.

* `event_type` - (Optional, String) Specifies the event type. The value can be `EVENT.SYS` (system event) or `EVENT.CUSTOM` (custom event).

* `from` - (Optional, Int) Specifies the start time of the query, which is a UNIX timestamp (ms). 

* `to` - (Optional, Int) Specifies the end time of the query, which is a UNIX timestamp (ms). 

* `limit` - (Optional, Int) Specifies the maximum number of records that can be queried at a time. Supported range: `1` to `100` (default)

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `events` - Specifies one or more pieces of event data. The structure is described below. 
* `meta_data` - Specifies the number of metadata records in the query result. The structure is described below.
    * `total` - Specifies the total number of events.


The `events` block supports:

* `event_name` - Specifies the event name.
* `event_type` - Specifies the event type.
* `event_count` - Specifies the number of occurrences of this event within the specified query time range.
* `latest_occur_time` - Specifies when the event occurred, which is a UNIX timestamp (ms).
* `latest_event_source` - Specifies the event source. If the event is a system event, the source is the namespace of each service. If the event is a custom event, the event source is defined by the user.
