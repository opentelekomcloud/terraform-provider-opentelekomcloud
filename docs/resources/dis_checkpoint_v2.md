---
subcategory: "Data Ingestion Service (DIS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dis_checkpoint_v2"
sidebar_current: "docs-opentelekomcloud-resource-dis-checkpoint-v2"
description: |-
Manages a DIS Checkpoint resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DIS checkpoint you can get at
[documentation portal](https://docs.otc.t-systems.com/data-ingestion-service/api-ref/api_description/checkpoint_management/index.html)

# opentelekomcloud_dis_checkpoint_v2

Manages a DIS Checkpoints in the OpenTelekomCloud DIS Service.

## Example Usage

```hcl
resource "opentelekomcloud_dis_stream_v2" "stream_1" {
  name                           = "my_stream"
  partition_count                = 3
  stream_type                    = "COMMON"
  retention_period               = 24
  auto_scale_min_partition_count = 1
  auto_scale_max_partition_count = 4
  compression_format             = "zip"

  data_type = "BLOB"

  tags = {
    foo = "bar"
  }
}

resource "opentelekomcloud_dis_app_v2" "app_1" {
  name = "my_app"
}

resource "opentelekomcloud_dis_checkpoint_v2" "checkpoint_1" {
  app_name        = opentelekomcloud_dis_app_v2.app_1.name
  stream_name     = opentelekomcloud_dis_stream_v2.stream_1.name
  partition_id    = "0"
  sequence_number = "0"
  metadata        = "my_first_checkpoint"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required) Name of the consumer application to be created
  The application name contains 1 to 200 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.

* `checkpoint_type` - (Required) Type of the checkpoint. `LAST_READ`: Only sequence numbers are recorded in databases.
  Default value: `LAST_READ`

* `stream_name` - (Required) Name of the stream. The stream name can contain 1 to 64 characters,
  including letters, digits, underscores (_), and hyphens (-).

* `partition_id` - (Required) Partition ID of the stream The value can be in either of the following formats:
  * `shardId-0000000000`
  * `0`

* `sequence_number` - (Required) Sequence number to be submitted, which is used to record the consumption
  checkpoint of the stream. Ensure that the sequence number is within the valid range.

* `metadata` - (Optional) Metadata information of the consumer application.
  Maximum length: 1000

## Attributes Reference

The following attributes are exported:

* `sequence_number` - Sequence number used to record the consumption checkpoint of the stream.

* `metadata` - Metadata information of the consumer application.
