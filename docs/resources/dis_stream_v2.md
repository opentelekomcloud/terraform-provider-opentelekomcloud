---
subcategory: "Data Ingestion Service (DIS)"
---

Up-to-date reference of API arguments for DIS stream you can get at
`https://docs.otc.t-systems.com/data-ingestion-service/api-ref/api_description/stream_management/index.html`.

# opentelekomcloud_dis_stream_v2

Manages a DIS Stream in the OpenTelekomCloud DIS Service.

## Example Usage

```hcl
resource "opentelekomcloud_dis_stream_v2" "stream_1" {
  stream_name                    = "MyStream"
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
```

## Argument Reference

The following arguments are supported:

* `stream_name` - (Required) Name of the stream. The stream name can contain 1 to 64 characters,
  including letters, digits, underscores (_), and hyphens (-).

* `partition_count` - (Required) Number of partitions. Partitions are the base throughput unit of a DIS stream.

* `retention_period` - (Optional) Period of time for which data is retained in the stream.
  Value range: 24-72 Unit: hour. If this parameter is left blank, the default value is used.
  `Maximum`: 72
  `Default`: 24

* `stream_type` - (Optional) Stream type.
  * `COMMON`: a common stream with a bandwidth of 1 MB/s
  * `ADVANCED`: an advanced stream with a bandwidth of 5 MB/s

* `data_type` - (Optional) Source data type.
  `BLOB`: a collection of binary data stored as a single entity in a database management system.
  Default value: `BLOB`.

* `compression_format` - (Optional) Data compression type. The following types are available:
  `snappy`, `gzip`, `zip`. Data is not compressed by default.

* `auto_scale_min_partition_count` - (Optional) Minimum number of partitions for automatic scale-down
  when auto-scaling is enabled. Minimum: 1.

* `auto_scale_max_partition_count` - (Optional) Maximum number of partitions for automatic scale-up when auto-scaling is enabled.

* `tags` - (Optional) Tags key/value pairs to associate with the instance.

## Attributes Reference

All above argument parameters can be exported as attribute parameters.

* `created` - Time when the stream is created. The value is a 13-bit timestamp..

* `readable_partition_count` - Total number of readable partitions (including partitions in ACTIVE and DELETED state).

* `writable_partition_count` - Total number of writable partitions (including partitions in ACTIVE state only).

* `status` - Current status of the stream
  `CREATING`: creating
  `RUNNING`: running
  `TERMINATING`: deleting
  `TERMINATED`: deleted

* `stream_id` - Unique identifier of the stream.

## Import

Stream can be imported using the stream id, e.g.

```shell
terraform import opentelekomcloud_dis_stream_v2.stream_1 Mh1NYdBHWdV0JPoNFTy
```
