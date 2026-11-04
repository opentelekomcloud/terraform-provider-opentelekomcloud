---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_topics_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-topics-v2"
description: |-
  Get the list of DMS RocketMQ topics within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ topics you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/topic_management/querying_the_topic_list.html)

# opentelekomcloud_dms_rocketmq_topics_v2

Use this data source to get the list of topics of a DMS RocketMQ instance within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "opentelekomcloud_dms_rocketmq_topics_v2" "test" {
  instance_id  = var.instance_id
  message_type = "FIFO"
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `name` - (Optional, String) Specifies the topic name.

* `permission` - (Optional, String) Specifies the topic permission.
  Valid values are **sub**, **pub** and **all**.

* `message_type` - (Optional, String) Specifies the message type.
  Valid values are **NORMAL**, **FIFO**, **DELAY** and **TRANSACTION**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the topics are queried.

* `topics` - The list of topics.
  The [topics](#dms_rocketmq_topics) structure is documented below.

<a name="dms_rocketmq_topics"></a>
The `topics` block supports:

* `name` - The topic name.

* `total_read_queue_num` - The total number of read queues.

* `total_write_queue_num` - The total number of write queues.

* `permission` - The topic permission.

* `message_type` - The message type.

* `brokers` - The list of associated brokers.
  The [brokers](#dms_rocketmq_topics_brokers) structure is documented below.

<a name="dms_rocketmq_topics_brokers"></a>
The `brokers` block supports:

* `broker_name` - The broker name.

* `read_queue_num` - The number of read queues.

* `write_queue_num` - The number of write queues.
