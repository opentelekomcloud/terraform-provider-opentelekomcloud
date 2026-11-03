---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_consumption_verify_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-consumption-verify-v2"
description: |-
  Verifies DMS RocketMQ message consumption within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ consumption verification you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/managing_messages/verifying_consumption.html)

# opentelekomcloud_dms_rocketmq_consumption_verify_v2

Verifies DMS RocketMQ message consumption within OpenTelekomCloud by making a consumer client consume the
specified messages again.

-> This is a one-time action resource. Deleting the resource only removes it from the state,
the messages remain re-delivered to the client.

-> The consumer group must have online consumers, otherwise the request fails.
Per-message results are reported in `resend_results`.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}
variable "topic" {}
variable "client_id" {}
variable "msg_id" {}

resource "opentelekomcloud_dms_rocketmq_consumption_verify_v2" "test" {
  instance_id = var.instance_id
  group       = var.group
  topic       = var.topic
  client_id   = var.client_id
  msg_id_list = [var.msg_id]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `group` - (Required, String, ForceNew) Specifies the consumer group name.
  Changing this creates a new resource.

* `client_id` - (Required, String, ForceNew) Specifies the ID of the consumer client.
  Changing this creates a new resource.

* `topic` - (Optional, String, ForceNew) Specifies the topic to which the messages belong.
  Changing this creates a new resource.

* `msg_id_list` - (Optional, List, ForceNew) Specifies the IDs of the messages to be consumed again.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `resend_results` - Indicates the consumption verification results.
  The [resend_results](#resend_results) structure is documented below.

<a name="resend_results"></a>
The `resend_results` block supports:

* `msg_id` - Indicates the message ID.

* `error_code` - Indicates the error code.

* `error_message` - Indicates the error message.
