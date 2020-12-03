---
subcategory: "Cloud Trace Service (CTS)"
---

# opentelekomcloud_cts_tracker_v1

Use this data source to get details about OpenTelekomCloud Cloud Trace Service.

## Example Usage
```hcl
variable "bucket_name" { }

data "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name = var.bucket_name
}

```

## Argument Reference

The following arguments are supported:

* `tracker_name` - (Optional) The tracker name.

* `bucket_name` - (Optional) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket.

* `status` - (Optional) Status of a tracker.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `is_support_smn` - Specifies SMN support.

* `topic_id` - The theme of the SMN service.

* `operations` - The trigger conditions for sending a notification

* `is_send_all_key_operation` - Specifies Typical or All operations for Trigger Condition.

* `need_notify_user_list` - The users using the login function.
