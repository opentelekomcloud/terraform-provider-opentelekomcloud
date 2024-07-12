---
subcategory: "Distributed Cache Service (DCS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dcs_maintainwindow_v1"
sidebar_current: "docs-opentelekomcloud-datasource-dcs-maintainwindow-v1"
description: |-
  Get the ID of an available DCS maintain window from OpenTelekomCloud
---

Up-to-date reference of API arguments for DCS certificate you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-cache-service/api-ref/apis_v2_recommended/other_apis/listing_maintenance_time_windows.html)

# opentelekomcloud_dcs_maintainwindow_v1

Use this data source to get the ID of an available OpenTelekomCloud DCS maintain window.

## Example Usage

```hcl
data "opentelekomcloud_dcs_maintainwindow_v1" "maintainwindow1" {
  seq = 1
}
```

## Argument Reference

* `seq` - (Required) Indicates the sequential number of a maintenance time window.

* `begin` - (Optional) Indicates the time at which a maintenance time window starts.

* `end` - (Required) Indicates the time at which a maintenance time window ends.

* `default` - (Required) Indicates whether a maintenance time window is set to the default time segment.

## Attributes Reference

`id` is set to the ID of the found maintainwindow. In addition, the following attributes are exported:

* `begin` - See Argument Reference above.

* `end` - See Argument Reference above.

* `default` - See Argument Reference above.
