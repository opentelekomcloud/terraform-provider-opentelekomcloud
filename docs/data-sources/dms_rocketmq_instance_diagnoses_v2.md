---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_instance_diagnoses_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-instance-diagnoses-v2"
description: |-
  Get the list of DMS RocketMQ instance diagnosis reports within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ instance diagnosis reports you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/instance_diagnosis/querying_an_instance_diagnosis_report_list.html)

# opentelekomcloud_dms_rocketmq_instance_diagnoses_v2

Use this data source to get the list of diagnosis reports of a DMS RocketMQ instance within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "opentelekomcloud_dms_rocketmq_instance_diagnoses_v2" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the diagnosis reports are queried.

* `reports` - The list of diagnosis reports.
  The [reports](#dms_rocketmq_instance_diagnoses) structure is documented below.

<a name="dms_rocketmq_instance_diagnoses"></a>
The `reports` block supports:

* `report_id` - The report ID.

* `group_name` - The consumer group name.

* `consumer_nums` - The number of consumers.

* `status` - The report status. The value can be `diagnosing`, `failed`, `deleted`, `finished`,
  `normal` or `abnormal`.

* `created_at` - The creation time.

* `abnormal_item_sum` - The number of abnormal items.

* `faulted_node_sum` - The number of abnormal nodes.
