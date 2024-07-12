---
subcategory: "Data Ingestion Service (DIS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dis_dump_task_v2"
sidebar_current: "docs-opentelekomcloud-resource-dis-dump-task-v2"
description: |-
  Manages a DIS Dump Task resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DIS dump task you can get at
[documentation portal](https://docs.otc.t-systems.com/data-ingestion-service/api-ref/api_description/dump_task_management/index.html)

# opentelekomcloud_dis_dump_task_v2

Manages a DIS Dump Task in the OpenTelekomCloud DIS Service.

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

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "my-dis-bucket"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "opentelekomcloud_dis_dump_task_v2" "task_1" {
  stream_name = opentelekomcloud_dis_stream_v2.stream_1.name
  destination = "OBS"

  obs_destination_descriptor {
    task_name             = "my_task"
    agency_name           = "dis_admin_agency"
    deliver_time_interval = 30
    consumer_strategy     = "LATEST"
    file_prefix           = "_pf"
    partition_format      = "yyyy/MM/dd/HH/mm"
    obs_bucket_path       = opentelekomcloud_obs_bucket.bucket.bucket
    destination_file_type = "text"
    record_delimiter      = "|"
  }
}
```

## Argument Reference

The following arguments are supported:

* `stream_name` - (Required) Name of the stream.

* `destination` - (Required) Dump destination. Possible values:
  `OBS`: Data is dumped to OBS.

* `obs_destination_descriptor` - (Optional) Parameter list of OBS to which data in the DIS stream will be dumped.
  * `task_name` - (Required) Name of the dump task. The task name consists of letters, digits, hyphens (-), and underscores (_). It must be a string of 1 to 64 characters.
  * `agency_name` - (Required) Name of the agency created on IAM. DIS uses an agency to access your specified resources.
    The parameters for creating an agency are as follows:
    * Agency Type: Cloud service
    * Cloud Service: DIS
    * Validity Period: unlimited
    * Scope: Global service
    * Project: OBS.
    * Select the Tenant Administrator role for the global service project.
  * `deliver_time_interval` - (Required) User-defined interval at which data is imported from the current DIS stream into OBS.
    If no data is pushed to the DIS stream during the current interval, no dump file package will be generated. Value range: `30`-`900`.
  * `obs_bucket_path` - (Required) Name of the OBS bucket used to store data from the DIS stream.
  * `consumer_strategy` - (Optional) Offset.
    `LATEST`: Maximum offset, indicating that the latest data will be extracted.
    `TRIM_HORIZON`: Minimum offset, indicating that the earliest data will be extracted.
  * `file_prefix` - (Optional) Directory to store files that will be dumped to OBS.
    Different directory levels are separated by slashes (/) and cannot start with slashes.
  * `partition_format` - (Optional) Directory structure of the object file written into OBS.
    The directory structure is in the format of yyyy/MM/dd/HH/mm (time at which the dump task was created).
    Possible values:
    * `yyyy`
    * `yyyy/MM`
    * `yyyy/MM/dd`
    * `yyyy/MM/dd/HH`
    * `yyyy/MM/dd/HH/mm`
  * `destination_file_type` - (Optional) Dump file format. Possible values: `text`
  * `record_delimiter` - (Optional) Delimiter for the dump file, which is used to separate the user data that is written into the dump file.

* `obs_processing_schema` - (Optional) Dump time directory generated based on the timestamp
  of the source data and the configured partition_format.
  Directory structure of the object file written into OBS.
  The directory structure is in the format of yyyy/MM/dd/HH/mm.
  * `timestamp_name` - (Required) Attribute name of the source data timestamp.
  * `timestamp_type` - (Required) Type of the source data timestamp.
    Possible values:
    * `String`
    * `Timestamp`
  * `timestamp_format` - (Required) OBS directory generated based on the timestamp format.
    This parameter is mandatory when the timestamp type of the source data is String.
  * yyyy/MM/dd HH:mm:ss
    Possible values:
    * `MM/dd/yyyy HH:mm:ss`
    * `dd/MM/yyyy HH:mm:ss`
    * `yyyy-MM-dd HH:mm:ss`
    * `MM-dd-yyyy HH:mm:ss`
    * `dd-MM-yyyy HH:mm:ss`

* `action` - (Optional) Dump task operation. The value can only be `start` or `stop`.

## Attributes Reference

All above argument parameters can be exported as attribute parameters.

* `name` - Name of the dump task.

* `task_id` - ID of the dump task.

* `created_at` - Time when the dump task is created.

* `last_transfer_timestamp` - Latest dump time of the dump task.

* `status` - Current status of the stream, can be:
  * `ERROR`: creating
  * `STARTING`: running
  * `PAUSED`: deleting
  * `RUNNING`: deleted
  * `DELETE`: deleted
  * `ABNORMAL`: deleted

* `partitions` - List of partition dump details.
  * `id`: Unique identifier of the partition.
  * `status`: Current status of the partition.
  * `hash_range`: Possible value range of the hash key used by the partition.
  * `sequence_number_range`: Sequence number range of the partition.
  * `parent_partitions`: Parent partition.
