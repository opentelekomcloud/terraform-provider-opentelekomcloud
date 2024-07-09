---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_logtank_transfer_v2"
sidebar_current: "docs-opentelekomcloud-resource-logtank-transfer-v2"
description: |-
Manages a LTS Log Transfer resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log transfer you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/log_transfer/index.html)

# opentelekomcloud_logtank_transfer_v2

Manage a log transfer resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "test-bucket"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_logtank_group_v2" "testacc_group" {
  group_name  = "test_group"
  ttl_in_days = 7
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "test-topic-1"
}

resource "opentelekomcloud_logtank_topic_v2" "testacc_topic-2" {
  group_id   = opentelekomcloud_logtank_group_v2.testacc_group.id
  topic_name = "test-topic-2"
}

resource "opentelekomcloud_logtank_transfer_v2" "transfer" {
  log_group_id    = opentelekomcloud_logtank_group_v2.testacc_group.id
  log_stream_ids  = [opentelekomcloud_logtank_topic_v2.testacc_topic.id, opentelekomcloud_logtank_topic_v2.testacc_topic-2.id]
  obs_bucket_name = opentelekomcloud_obs_bucket.bucket.bucket
  storage_format  = "JSON"
  switch_on       = false
  period          = 30
  period_unit     = "min"
  prefix_name     = "prefix"
  dir_prefix_name = "dir"
}
```

## Argument Reference

The following arguments are supported:

* `log_group_id` - (Required) Specifies the ID of a log transfer.

* `log_stream_ids` - (Required) Specifies the log topics(streams) ids.

* `obs_bucket_name` - (Required) Specifies the name of an OBS bucket.

* `storage_format` - (Required) Indicates storage format for logs. Possible values are: `RAW`, `JSON`.

* `period` - (Required) Indicates the length of the log transfer interval.
  Possible values: `1`, `2`, `3`, `5`, `6`, `12`, and `30`.

* `period_unit` - (Required) Indicates the unit of the log transfer interval.
  Possible values: `min`, `hour`.

~> **Warning** The log transfer interval is specified by the combination of the values of `obs_period` and `obs_period_unit`,
and must be set to one of the following: `2 min`, `5 min`, `30 min`, `1 hour`, `3 hours`, `6 hours`, and `12 hours`.

* `switch_on` - (Optional) Indicates whether the log transfer is enabled. Default: `true`.

* `prefix_name` - (Optional) Indicates the file name prefix of the log files transferred to an OBS bucket.

* `dir_prefix_name` - (Optional) Indicates a custom path to store the log files.

## Attributes Reference

The following attributes are exported:

* `id` - The log transfer ID.

* `log_group_name` - The name of log group.

* `log_transfer_mode` - The log transfer mode. `cycle` indicates periodical transfer.

* `status` - The log transfer status.
  `ENABLE`/`DISABLE` indicates that log transfer is enabled/disabled.
  `EXCEPTION` indicates that log transfer is abnormal.

* `log_transfer_type` - The log transfer type.

* `create_time` - Specifies the time when a log transfer was created.

* `obs_encryption_id` - Specifies the KMS key ID for an OBS transfer task.

* `obs_encryption_enable` - Specifies whether OBS bucket encryption is enabled.
