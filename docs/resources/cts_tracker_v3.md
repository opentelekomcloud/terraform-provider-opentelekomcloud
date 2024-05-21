---
subcategory: "Cloud Trace Service (CTS)"
---

Up-to-date reference of API arguments for CTS tracker you can get at
`https://docs.sc.otc.t-systems.com/api/cts/cts_api_0201.html`.

# opentelekomcloud_cts_tracker_v3

Allows you to collect, store, and query cloud resource operation records.

~> **Warning** `opentelekomcloud_cts_tracker_v3` is not available for `eu-nl` region.

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

* `status` - (Required) Specifies whether tracker is `enabled` or `disabled`.

* `is_lts_enabled` - (Optional) Specifies whether to enable trace analysis.

* `bucket_name` - (Optional) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket.

* `is_obs_created` - (Optional) Specifies whether the OBS bucket is automatically created by the tracker.

* `is_sort_by_service` - (Optional) Specifies whether to sort the path by cloud service. If this option is enabled,
  the cloud service name is added to the transfer file path. Default: `true`.

* `compress_type` - (Optional) Specifies the compression type. Default value is `gzip`.
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
