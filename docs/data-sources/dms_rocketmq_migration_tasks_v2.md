---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_migration_tasks_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-migration-tasks-v2"
description: |-
  Get the list of DMS RocketMQ metadata migration tasks within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ migration tasks you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/metadata_migration/querying_all_migration_tasks_or_a_specified_migration_task_of_an_instance.html)

# opentelekomcloud_dms_rocketmq_migration_tasks_v2

Use this data source to get the list of metadata migration tasks of a DMS RocketMQ instance within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "opentelekomcloud_dms_rocketmq_migration_tasks_v2" "test" {
  instance_id = var.instance_id
  type        = "rocketmq"
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `task_id` - (Optional, String) Specifies the ID of the migration task.

* `name` - (Optional, String) Specifies the name of the migration task.

* `type` - (Optional, String) Specifies the migration type. Valid values are:
  + `rocketmq` - from RocketMQ to RocketMQ.
  + `rabbitToRocket` - from RabbitMQ to RocketMQ.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the migration tasks are queried.

* `tasks` - The list of migration tasks.
  The [tasks](#dms_rocketmq_migration_tasks) structure is documented below.

<a name="dms_rocketmq_migration_tasks"></a>
The `tasks` block supports:

* `id` - The ID of the migration task.

* `name` - The name of the migration task.

* `start_date` - The start time of the migration task.

* `status` - The status of the migration task. The value can be `creating`, `starting`, `failed` or `finished`.

* `type` - The migration type.
