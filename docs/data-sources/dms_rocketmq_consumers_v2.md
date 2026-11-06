---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_consumers_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-consumers-v2"
description: |-
  Get the list of consumers of a DMS RocketMQ consumer group within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ consumers you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/consumer_group_management/querying_the_consumer_list.html)

# opentelekomcloud_dms_rocketmq_consumers_v2

Use this data source to get the list of consumers of a DMS RocketMQ consumer group within OpenTelekomCloud.

-> The request fails with the `consumer group ... not online` error if the consumer group has no online consumers.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}

data "opentelekomcloud_dms_rocketmq_consumers_v2" "test" {
  instance_id = var.instance_id
  group       = var.group_name
  is_detail   = true
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `group` - (Required, String) Specifies the consumer group name.

* `is_detail` - (Optional, Bool) Specifies whether to query the consumer details.
  If it is not `true`, a simple consumer list is queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the consumers are queried.

* `online` - Whether the consumer group is online.

* `subscription_consistency` - Whether subscriptions are consistent.

* `clients` - The list of consumers.
  The [clients](#dms_rocketmq_consumers_clients) structure is documented below.

<a name="dms_rocketmq_consumers_clients"></a>
The `clients` block supports:

* `language` - The client language.

* `version` - The client version.

* `client_id` - The client ID.

* `client_addr` - The client address.

* `subscriptions` - The subscription list.
  The [subscriptions](#dms_rocketmq_consumers_subscriptions) structure is documented below.

<a name="dms_rocketmq_consumers_subscriptions"></a>
The `subscriptions` block supports:

* `topic` - The name of the subscribed topic.

* `type` - The subscription type. The value can be `TAG` or `SQL92`.

* `expression` - The subscription tag.
