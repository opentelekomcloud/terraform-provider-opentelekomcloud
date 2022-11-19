---
subcategory: "Cloud Trace Service (CTS)"
---

# opentelekomcloud_cts_tracker_v1

Allows you to collect, store, and query cloud resource operation records.

-> A single tracker can be created for current CTS version.

## Example Usage

```hcl
variable "bucket_name" {}

resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = var.bucket_name
  file_prefix_name = "yO8Q"
  is_lts_enabled   = true
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket.

* `is_lts_enabled` - (Optional) Specifies whether to enable trace analysis.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tracker_name` - The tracker name. Currently, only tracker `system` is available.

* `region` - Specifies the tracker region.

* `status` - Specifies current status of the tracker.

* `log_group_name` - Specifies LTS log group name.

* `log_topic_name` - Specifies LTS log stream.

## Import

CTS tracker can be imported using  `tracker_name`, e.g.

```shell
$ terraform import opentelekomcloud_cts_tracker_v1.tracker system
```
