---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_consumer_groups_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-consumer-groups-v2"
description: |-
  Get the list of DMS RocketMQ consumer groups within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ consumer groups you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/consumer_group_management/querying_the_consumer_group_list.html)

# opentelekomcloud_dms_rocketmq_consumer_groups_v2

Use this data source to get the list of consumer groups of a DMS RocketMQ instance within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}

data "opentelekomcloud_dms_rocketmq_consumer_groups_v2" "test" {
  instance_id = var.instance_id
  name        = var.group_name
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `name` - (Optional, String) Specifies the consumer group name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the consumer groups are queried.

* `groups` - The list of consumer groups.
  The [groups](#dms_rocketmq_consumer_groups) structure is documented below.

<a name="dms_rocketmq_consumer_groups"></a>
The `groups` block supports:

* `name` - The consumer group name.

* `description` - The consumer group description.

* `enabled` - Whether consumption is allowed.

* `broadcast` - Whether broadcast is enabled.

* `brokers` - The list of associated brokers.

* `retry_max_time` - The maximum number of retries.

* `permissions` - The permission set.

* `consume_orderly` - Whether orderly consumption is enabled.

* `group_online` - Whether the consumer group is online.

* `created_at` - The creation time, as a timestamp in milliseconds.
