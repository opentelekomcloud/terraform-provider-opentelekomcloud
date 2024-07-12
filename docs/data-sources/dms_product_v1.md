---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_product_v1"
sidebar_current: "docs-opentelekomcloud-datasource-dms-product-v1"
description: |-
  Get available DMS product from OpenTelekomCloud
---

Up-to-date reference of API arguments for DMS product you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service/api-ref/apis_v2_recommended/other_apis/querying_product_specifications_list.html)

# opentelekomcloud_dms_product_v1

Use this data source to get the ID of an available DMS product within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_dms_product_v1" "product1" {
  engine            = "kafka"
  version           = "2.3.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}
```

## Argument Reference

* `engine` - (Required) Indicates the name of a message engine. Only "kafka" is supported now.

* `version` - (Optional) Indicates the version of a message engine. Only "2.3.0" is supported now.

* `instance_type` - (Required) Indicates an instance type. Only "cluster" is supported now.

* `vm_specification` - (Optional) Indicates VM specifications.

* `storage` - (Optional) Indicates the message storage space.

* `bandwidth` - (Optional) Indicates the baseline bandwidth of a Kafka instance.

* `partition_num` - (Optional) Indicates the maximum number of topics that can be created for a Kafka instance.

* `storage_spec_code` - (Optional) Indicates an I/O specification.

* `io_type` - (Optional) Indicates an I/O type.

* `node_num` - (Optional) Indicates the number of nodes in a cluster.

## Attributes Reference

`id` is set to the ID of the found product. In addition, the following attributes are exported:

* `engine` - See Argument Reference above.

* `version` - See Argument Reference above.

* `instance_type` - See Argument Reference above.

* `vm_specification` - See Argument Reference above.

* `bandwidth` - See Argument Reference above.

* `partition_num` - See Argument Reference above.

* `storage_spec_code` - See Argument Reference above.

* `io_type` - See Argument Reference above.

* `node_num` - See Argument Reference above.
