---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_dead_letter_resend_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-dead-letter-resend-v2"
description: |-
  Resends DMS RocketMQ dead letter messages within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ dead letter message resending you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/managing_messages/resending_dead_letter_messages.html)

# opentelekomcloud_dms_rocketmq_dead_letter_resend_v2

Resends DMS RocketMQ dead letter messages within OpenTelekomCloud.

-> This is a one-time action resource. Deleting the resource only removes it from the state,
the dead letter messages remain resent.

-> The consumer group must have online consumers, otherwise the request fails.
Per-message results are reported in `resend_results`.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}
variable "msg_id" {}

resource "opentelekomcloud_dms_rocketmq_dead_letter_resend_v2" "test" {
  instance_id = var.instance_id
  topic       = "%DLQ%${var.group}"
  msg_id_list = [var.msg_id]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `topic` - (Optional, String, ForceNew) Specifies the dead letter topic name, in format `%DLQ%{group}`.
  Changing this creates a new resource.

* `msg_id_list` - (Optional, List, ForceNew) Specifies the IDs of the messages to be resent.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `resend_results` - Indicates the results of resending the dead letter messages.
  The [resend_results](#resend_results) structure is documented below.

<a name="resend_results"></a>
The `resend_results` block supports:

* `msg_id` - Indicates the message ID.

* `error_code` - Indicates the error code.

* `error_message` - Indicates the error message.
