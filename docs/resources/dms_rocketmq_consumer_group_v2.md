---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_consumer_group_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-consumer-group-v2"
description: |-
  Manages a DMS RocketMQ consumer group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ consumer group you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/consumer_group_management/index.html)

# opentelekomcloud_dms_rocketmq_consumer_group_v2

Manages a DMS RocketMQ consumer group in the OpenTelekomCloud DMS Service.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_dms_rocketmq_consumer_group_v2" "test" {
  instance_id    = var.instance_id
  name           = "group_test"
  retry_max_time = 16
  description    = "group description"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the consumer group. Enter 3 to 64 characters.
  Use only letters, digits, percent (%), vertical bars (|), hyphens (-), and underscores (_).
  Changing this creates a new resource.

* `enabled` - (Optional, Bool) Specifies whether consumption is allowed. Defaults to **true**.

* `broadcast` - (Optional, Bool) Specifies whether to broadcast. Defaults to **false**.

* `retry_max_time` - (Optional, Int) Specifies the maximum number of retries. Value range: **1**–**16**.

* `consume_orderly` - (Optional, Bool) Specifies whether orderly consumption is enabled. Defaults to **false**.

* `description` - (Optional, String) Specifies the description of the consumer group.
  It contains 0 to 200 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 5 minutes.

## Import

DMS RocketMQ consumer group can be imported using the instance ID and group name separated by a slash, e.g.

```
 $ terraform import opentelekomcloud_dms_rocketmq_consumer_group_v2.test 8d3c7938-dc47-4937-a30f-c80de381c5e3/group_test
```
