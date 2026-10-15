---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_migration_task_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-migration-task-v2"
description: |-
  Manages a DMS RocketMQ metadata migration task resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ metadata migration task you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/metadata_migration/index.html)

# opentelekomcloud_dms_rocketmq_migration_task_v2

Manages a DMS RocketMQ metadata migration task in the OpenTelekomCloud DMS Service.
The task migrates metadata of a self-built RocketMQ or RabbitMQ instance to a DMS RocketMQ instance.
The resource waits until the task is finished.

-> Deleting the resource deletes only the task, the migrated topics and consumer groups stay in the instance.

## Example Usage

### Migrating metadata from RocketMQ

```hcl
variable "instance_id" {}

resource "opentelekomcloud_dms_rocketmq_migration_task_v2" "test" {
  instance_id = var.instance_id
  name        = "task_test"
  type        = "rocketmq"
  overwrite   = true

  topic_configs {
    topic_name        = "topic_test"
    perm              = 6
    read_queue_nums   = 3
    write_queue_nums  = 3
    topic_filter_type = "SINGLE_TAG"
  }

  subscription_groups {
    group_name      = "group_test"
    retry_max_times = 2
  }
}
```

### Migrating metadata from RabbitMQ

```hcl
variable "instance_id" {}

resource "opentelekomcloud_dms_rocketmq_migration_task_v2" "test" {
  instance_id = var.instance_id
  name        = "task_test"
  type        = "rabbitToRocket"

  vhosts {
    name = "vhost_test"
  }

  queues {
    vhost   = "vhost_test"
    name    = "queue_test"
    durable = true
  }

  exchanges {
    vhost   = "vhost_test"
    name    = "exchange_test"
    type    = "direct"
    durable = true
  }

  bindings {
    vhost            = "vhost_test"
    source           = "exchange_test"
    destination      = "queue_test"
    destination_type = "queue"
    routing_key      = "routing_key_test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the migration task. Starts with a letter, consists of
  4 to 64 characters, and contains only letters, digits, underscores (_), and hyphens (-).
  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the migration task type.
  Value options: **rocketmq** (from RocketMQ to RocketMQ), **rabbitToRocket** (from RabbitMQ to RocketMQ).
  Changing this creates a new resource.

* `overwrite` - (Optional, Bool, ForceNew) Specifies whether to overwrite configurations of the existing topics and
  consumer groups with the same name. If **false**, an error is reported when a topic or group already exists.
  Defaults to **false**. Changing this creates a new resource.

* `topic_configs` - (Optional, List, ForceNew) Specifies the topic metadata. Mandatory for **rocketmq** tasks.
  The [topic_configs](#topic_configs) structure is documented below.
  Changing this creates a new resource.

* `subscription_groups` - (Optional, List, ForceNew) Specifies the consumer group metadata.
  Mandatory for **rocketmq** tasks.
  The [subscription_groups](#subscription_groups) structure is documented below.
  Changing this creates a new resource.

* `vhosts` - (Optional, List, ForceNew) Specifies the virtual host metadata. Mandatory for **rabbitToRocket** tasks.
  The [vhosts](#vhosts) structure is documented below.
  Changing this creates a new resource.

* `queues` - (Optional, List, ForceNew) Specifies the queue metadata. Mandatory for **rabbitToRocket** tasks.
  The [queues](#queues) structure is documented below.
  Changing this creates a new resource.

* `exchanges` - (Optional, List, ForceNew) Specifies the exchange metadata. Mandatory for **rabbitToRocket** tasks.
  The [exchanges](#exchanges) structure is documented below.
  Changing this creates a new resource.

* `bindings` - (Optional, List, ForceNew) Specifies the binding metadata. Mandatory for **rabbitToRocket** tasks.
  The [bindings](#bindings) structure is documented below.
  Changing this creates a new resource.

<a name="topic_configs"></a>
The `topic_configs` block supports:

* `topic_name` - (Required, String, ForceNew) Specifies the topic name.

* `order` - (Optional, Bool, ForceNew) Specifies whether messages are ordered. Defaults to **false**.

* `perm` - (Optional, Int, ForceNew) Specifies the topic permissions. Defaults to **6**.

* `read_queue_nums` - (Optional, Int, ForceNew) Specifies the number of read queues. Defaults to **16**.

* `write_queue_nums` - (Optional, Int, ForceNew) Specifies the number of write queues. Defaults to **16**.

* `topic_filter_type` - (Optional, String, ForceNew) Specifies the topic filtering type.
  Value options: **SINGLE_TAG**, **MULTI_TAG**.

* `topic_sys_flag` - (Optional, Int, ForceNew) Specifies the topic system flag. Defaults to **0**.

<a name="subscription_groups"></a>
The `subscription_groups` block supports:

* `group_name` - (Required, String, ForceNew) Specifies the consumer group name.

* `consume_broadcast_enable` - (Optional, Bool, ForceNew) Specifies whether to allow broadcast.
  Defaults to **false**.

* `consume_enable` - (Optional, Bool, ForceNew) Specifies whether to enable consumption. Defaults to **true**.

* `consume_from_min_enable` - (Optional, Bool, ForceNew) Specifies whether to enable consumption from the earliest
  offset. Defaults to **true**.

* `notify_consumer_ids_changed_enable` - (Optional, Bool, ForceNew) Specifies whether to notify consumer ID changes.
  Defaults to **true**.

* `retry_max_times` - (Optional, Int, ForceNew) Specifies the maximum number of consumption retries.
  Defaults to **16**.

* `retry_queue_nums` - (Optional, Int, ForceNew) Specifies the number of retry queues. Defaults to **1**.

* `which_broker_when_consume_slow` - (Optional, Int, ForceNew) Specifies the ID of the broker selected for slow
  consumption. Defaults to **1**.

<a name="vhosts"></a>
The `vhosts` block supports:

* `name` - (Optional, String, ForceNew) Specifies the virtual host name.

<a name="queues"></a>
The `queues` block supports:

* `vhost` - (Optional, String, ForceNew) Specifies the virtual host name.

* `name` - (Optional, String, ForceNew) Specifies the queue name.

* `durable` - (Optional, Bool, ForceNew) Specifies whether to enable data persistence. Defaults to **false**.

<a name="exchanges"></a>
The `exchanges` block supports:

* `vhost` - (Optional, String, ForceNew) Specifies the virtual host name.

* `name` - (Optional, String, ForceNew) Specifies the exchange name.

* `type` - (Optional, String, ForceNew) Specifies the exchange type.
  Value options: **topic**, **direct**, **fanout**, **headers**.

* `durable` - (Optional, Bool, ForceNew) Specifies whether to enable data persistence. Defaults to **false**.

<a name="bindings"></a>
The `bindings` block supports:

* `vhost` - (Optional, String, ForceNew) Specifies the virtual host name.

* `source` - (Optional, String, ForceNew) Specifies the message source.

* `destination` - (Optional, String, ForceNew) Specifies the message target.

* `destination_type` - (Optional, String, ForceNew) Specifies the target type.
  Value options: **exchange**, **queue**.

* `routing_key` - (Optional, String, ForceNew) Specifies the routing key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<task_id>`.
* `start_date` - Indicates the start time of the migration task.
* `status` - Indicates the status of the migration task. Values: **creating**, **starting**, **failed**, **finished**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

DMS RocketMQ migration task can be imported using the instance ID and task ID separated by a slash, e.g.

```
 $ terraform import opentelekomcloud_dms_rocketmq_migration_task_v2.test 8d3c7938-dc47-4937-a30f-c80de381c5e3/6cf4dcd3-8471-4139-8b5b-8a3a71f704c7
```

Note that `overwrite` and the metadata blocks are not returned by the API, so they are missing in the imported state.
You can ignore changes as below.

```hcl
resource "opentelekomcloud_dms_rocketmq_migration_task_v2" "test" {
  lifecycle {
    ignore_changes = [
      overwrite,
      topic_configs,
      subscription_groups,
      vhosts,
      queues,
      exchanges,
      bindings,
    ]
  }
}
```
