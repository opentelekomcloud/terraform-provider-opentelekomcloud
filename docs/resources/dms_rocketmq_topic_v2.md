---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_topic_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-topic-v2"
description: |-
  Manages a DMS RocketMQ topic resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ topic you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/topic_management/index.html)

# opentelekomcloud_dms_rocketmq_topic_v2

Manages a DMS RocketMQ topic in the OpenTelekomCloud DMS Service.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_dms_rocketmq_topic_v2" "test" {
  instance_id  = var.instance_id
  name         = "topic_test"
  queue_num    = 4
  message_type = "NORMAL"
  description  = "topic description"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the topic. Enter 3 to 64 characters.
  Use only letters, digits, percent (%), vertical bars (|), hyphens (-), and underscores (_).
  Changing this creates a new resource.

* `queue_num` - (Optional, Int, ForceNew) Specifies the total number of queues. Value range: **1**–**50**.
  Changing this creates a new resource.

* `message_type` - (Optional, String, ForceNew) Specifies the message type.
  Value options: **NORMAL**, **FIFO** (ordered), **DELAY** (scheduled) and **TRANSACTION**.
  Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of the topic.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<name>`.
* `permission` - Indicates the permission of the topic. Values: **sub**, **pub**, **all**.
* `total_read_queue_num` - Indicates the total number of read queues.
* `total_write_queue_num` - Indicates the total number of write queues.

## Import

DMS RocketMQ topic can be imported using the instance ID and topic name separated by a slash, e.g.

```
 $ terraform import opentelekomcloud_dms_rocketmq_topic_v2.test 8d3c7938-dc47-4937-a30f-c80de381c5e3/topic_test
```

Note that `queue_num` and `description` are not returned by the API, so they are missing in the imported state.
You can ignore changes as below.

```hcl
resource "opentelekomcloud_dms_rocketmq_topic_v2" "test" {
  lifecycle {
    ignore_changes = [
      queue_num,
      description,
    ]
  }
}
```
