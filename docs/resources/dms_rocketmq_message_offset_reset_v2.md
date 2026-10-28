---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_message_offset_reset_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-message-offset-reset-v2"
description: |-
  Resets the consumer offset of a DMS RocketMQ consumer group within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ consumer offset reset you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/consumer_group_management/resetting_the_consumer_offset.html)

# opentelekomcloud_dms_rocketmq_message_offset_reset_v2

Resets the consumer offset of a DMS RocketMQ consumer group for a topic within OpenTelekomCloud.

-> This is a one-time action resource. Deleting the resource only removes it from the state,
the consumer offset remains reset.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}
variable "topic" {}

resource "opentelekomcloud_dms_rocketmq_message_offset_reset_v2" "test" {
  instance_id = var.instance_id
  group       = var.group
  topic       = var.topic
  timestamp   = "1662652800000"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `group` - (Required, String, ForceNew) Specifies the consumer group name.
  Changing this creates a new resource.

* `topic` - (Required, String, ForceNew) Specifies the topic to be reset.
  Changing this creates a new resource.

* `timestamp` - (Required, String, ForceNew) Specifies the reset time, Unix timestamp in milliseconds.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<group>/<topic>`.

* `queues` - Indicates the queues to be reset.
  The [queues](#queues) structure is documented below.

<a name="queues"></a>
The `queues` block supports:

* `broker_name` - Indicates the broker where the queue is located.

* `queue_id` - Indicates the queue ID.

* `timestamp_offset` - Indicates the reset consumer offset.
