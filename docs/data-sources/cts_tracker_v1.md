---
subcategory: "Cloud Trace Service (CTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cts_tracker_v1"
sidebar_current: "docs-opentelekomcloud-datasource-cts-tracker-v1"
description: |-
  Get details about CTS tracker from OpenTelekomCloud
---

Up-to-date reference of API arguments for CTS tracker you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-trace-service/api-ref/out-of-date_apis/tracker_management/querying_a_tracker.html#cts-api-0016)

# opentelekomcloud_cts_tracker_v1

Use this data source to get details about OpenTelekomCloud Cloud Trace Service.

## Example Usage
```hcl
data "opentelekomcloud_cts_tracker_v1" "tracker_v1" {}

```

## Attributes Reference

The following arguments are supported:

* `bucket_name` - The OBS bucket name for a tracker to store trace info.

* `file_prefix_name` - The prefix of a log that needs to be stored in an OBS bucket.

* `is_lts_enabled` - Specifies whether to enable trace analysis.

* `tracker_name` - The tracker name. Currently, only tracker `system` is available.

* `region` - Specifies the tracker region.

* `status` - Specifies current status of the tracker.

* `log_group_name` - Specifies LTS log group name.

* `log_topic_name` - Specifies LTS log stream.
