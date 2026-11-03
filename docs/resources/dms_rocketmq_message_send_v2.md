---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_message_send_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-message-send-v2"
description: |-
  Sends a message to a DMS RocketMQ topic within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ message sending you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/managing_messages/sending_messages.html)

# opentelekomcloud_dms_rocketmq_message_send_v2

Sends a message to a DMS RocketMQ topic within OpenTelekomCloud.

-> This is a one-time action resource. Deleting the resource only removes it from the state,
the message remains in the topic.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic" {}

resource "opentelekomcloud_dms_rocketmq_message_send_v2" "test" {
  instance_id = var.instance_id
  topic       = var.topic
  body        = "message content"

  property_list {
    name  = "KEYS"
    value = "owner"
  }
  property_list {
    name  = "TAGS"
    value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `topic` - (Required, String, ForceNew) Specifies the topic name.
  Changing this creates a new resource.

* `body` - (Required, String, ForceNew) Specifies the message content.
  Changing this creates a new resource.

* `property_list` - (Optional, List, ForceNew) Specifies the feature list of the message.
  The [property_list](#property_list) structure is documented below.
  Changing this creates a new resource.

<a name="property_list"></a>
The `property_list` block supports:

* `name` - (Optional, String, ForceNew) Specifies the feature name.
  Changing this creates a new resource.

* `value` - (Optional, String, ForceNew) Specifies the feature value.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `msg_id`.

* `msg_id` - Indicates the message ID.

* `queue_id` - Indicates the queue ID.

* `queue_offset` - Indicates the queue offset.

* `broker_name` - Indicates the broker name.
