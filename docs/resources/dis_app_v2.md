---
subcategory: "Data Ingestion Service (DIS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dis_app_v2"
sidebar_current: "docs-opentelekomcloud-resource-dis-app-v2"
description: |-
  Manages a DIS App resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DIS app you can get at
[documentation portal](https://docs.otc.t-systems.com/data-ingestion-service/api-ref/api_description/app_management/index.html)

# opentelekomcloud_dis_app_v2

Manages a DIS Apps in the OpenTelekomCloud DIS Service.

## Example Usage

```hcl
resource "opentelekomcloud_dis_app_v2" "app_1" {
  name = "app_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the consumer application to be created
  The application name contains 1 to 200 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.

## Attributes Reference

Above argument parameter can be exported as attribute parameters.

* `created` - Time when the app is created. The value is a timestamp.

* `id` - Unique identifier of the app.

* `commit_checkpoint_stream_names` - List of associated streams.

* `partition_consuming_states` - Associated partitions details.
  * `id`: Partition Id.
  * `status`: Partition Status, can be:
    * `CREATING`
    * `ACTIVE`
    * `DELETED`
    * `EXPIRED`
  * `checkpoint_type`: Type of the checkpoint.
  * `sequence_number`: Partition Sequence Number
  * `latest_offset`: Partition data latest offset
  * `earliest_offset`: Partition data earliest offset


## Import

App can be imported using the app name, e.g.

```shell
terraform import opentelekomcloud_dis_app_v2.app_1 app_name
```
