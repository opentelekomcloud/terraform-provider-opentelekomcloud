---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_instance_diagnosis_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-instance-diagnosis-v2"
description: |-
  Manages a DMS RocketMQ instance diagnosis resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ instance diagnosis you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/instance_diagnosis/index.html)

# opentelekomcloud_dms_rocketmq_instance_diagnosis_v2

Manages a DMS RocketMQ instance diagnosis in the OpenTelekomCloud DMS Service.
The resource creates a diagnosis task for a consumer group and waits until the diagnosis report is generated.

-> Deleting the resource deletes the diagnosis report.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}

resource "opentelekomcloud_dms_rocketmq_instance_diagnosis_v2" "test" {
  instance_id = var.instance_id
  group_name  = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `group_name` - (Required, String, ForceNew) Specifies the name of the consumer group to diagnose.
  Changing this creates a new resource.

* `node_id_list` - (Optional, List, ForceNew) Specifies the list of node IDs to diagnose.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<report_id>`.

* `report_id` - Indicates the ID of the diagnosis report.

* `consumer_nums` - Indicates the number of consumers.

* `status` - Indicates the status of the diagnosis report.
  Values: **diagnosing**, **failed**, **deleted**, **finished**, **normal** (no problems found),
  **abnormal** (problems found).

* `created_at` - Indicates the time when the diagnosis report was generated, in RFC3339 format.

* `abnormal_item_sum` - Indicates the number of abnormal items.

* `faulted_node_sum` - Indicates the number of abnormal nodes.

* `online` - Indicates whether the consumer group is online.

* `message_accumulation` - Indicates the number of stacked messages.

* `subscription_consistency` - Indicates whether the subscription relationship is consistent.

* `duplicate_client_id` - Indicates whether there are duplicate client IDs.

* `different_consumer_type` - Indicates whether there are inconsistent consumption types.

* `subscriptions` - Indicates the subscriber list.
  The [subscriptions](#subscriptions) structure is documented below.

* `diagnosis_node_report_list` - Indicates the diagnosis node report list.
  The [diagnosis_node_report_list](#diagnosis_node_report_list) structure is documented below.

<a name="subscriptions"></a>
The `subscriptions` block supports:

* `topic_name` - Indicates the topic name.

* `consumers_in_tags` - Indicates the consumer tag list.
  The [consumers_in_tags](#consumers_in_tags) structure is documented below.

<a name="consumers_in_tags"></a>
The `consumers_in_tags` block supports:

* `consumers` - Indicates the consumer list.

* `tag_name` - Indicates the tag name.

<a name="diagnosis_node_report_list"></a>
The `diagnosis_node_report_list` block supports:

* `node_id` - Indicates the node ID.

* `is_faulted` - Indicates whether the node is faulty.

* `abnormal_item_sum` - Indicates the total number of exceptions.

* `message_accumulation` - Indicates the number of stacked messages.

* `dead_lock` - Indicates whether a deadlock occurs.

* `deadlock_thread` - Indicates the deadlock thread.

* `stack_id` - Indicates the thread ID.

* `is_pop` - Indicates whether pop consumption is used.

* `consume_type` - Indicates the consumption type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

DMS RocketMQ instance diagnosis can be imported using the instance ID and report ID separated by a slash, e.g.

```
 $ terraform import opentelekomcloud_dms_rocketmq_instance_diagnosis_v2.test 8d3c7938-dc47-4937-a30f-c80de381c5e3/b00a070e-75e0-4463-b53e-fbc29dcff951
```

Note that `node_id_list` is not returned by the API, so it is missing in the imported state.
You can ignore changes as below.

```hcl
resource "opentelekomcloud_dms_rocketmq_instance_diagnosis_v2" "test" {
  lifecycle {
    ignore_changes = [
      node_id_list,
    ]
  }
}
```
