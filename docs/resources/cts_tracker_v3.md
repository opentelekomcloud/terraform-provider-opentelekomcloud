---
subcategory: "Cloud Trace Service (CTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cts_tracker_v3"
sidebar_current: "docs-opentelekomcloud-resource-cts-tracker-v3"
description: |-
  Manages a CTS Tracker v3 resource within T Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for CTS tracker you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-trace-service/api-ref/v3_apis_recommended/tracker_management/index.html#cts-api-0320)

# opentelekomcloud_cts_tracker_v3

Allows you to collect, store, and query cloud resource operation records.

-> A single tracker can be created for current CTS version.

## Example Usage

```hcl
variable "bucket_name" {}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name      = var.bucket_name
  file_prefix_name = "prefix"
  is_lts_enabled   = true
  status           = "enabled"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, String) Specifies whether tracker is `enabled` or `disabled`.

* `is_lts_enabled` - (Optional, Boolean) Specifies whether to enable trace analysis.

* `is_support_validate` (Optional, Boolean) Specifies Whether trace file verification is enabled for trace transfer. When this function is enabled, integrity verification will be performed to check whether trace files in OBS buckets have been tampered with.

* `bucket_name` - (Optional, String) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional, String) The prefix of a log that needs to be stored in an OBS bucket.

* `is_obs_created` - (Optional, Boolean) Specifies whether the OBS bucket is automatically created by the tracker.

* `is_sort_by_service` - (Optional, Boolean) Specifies whether to sort the path by cloud service. If this option is enabled,
  the cloud service name is added to the transfer file path. Default: `true`.

* `compress_type` - (Optional, String) Specifies the compression type. Default value is `gzip`.
  The valid values are as follows:
    + **gzip**: compression.
    + **json**: no compression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tracker_name` - The tracker name. Currently, only tracker `system` is available.

* `tracker_type` - The tracker type. Currently, only tracker `system` is available.

* `id` - Specifies the tracker id.

* `domain_id` - Specifies domain id of the tracker.

* `project_id` - Specifies project id of the tracker.

* `log_group_name` - Specifies LTS log group name.

* `log_topic_name` - Specifies LTS log stream.

* `detail` - Specifies the cause of the abnormal status, and its value in case of errors.

## Import

CTS tracker can be imported using `tracker_name`, e.g.

```shell
$ terraform import opentelekomcloud_cts_tracker_v3.tracker system
```
