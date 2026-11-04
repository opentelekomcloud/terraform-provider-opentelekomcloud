---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_flavors_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-flavors-v2"
description: |-
  Get the list of available DMS RocketMQ flavors within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ flavors you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/other_apis/querying_flavor_list.html)

# opentelekomcloud_dms_rocketmq_flavors_v2

Use this data source to get the list of available DMS RocketMQ flavors within OpenTelekomCloud.

## Example Usage

### Query the RocketMQ flavor details of the specified ID

```hcl
data "opentelekomcloud_dms_rocketmq_flavors_v2" "test" {
  flavor_id = "rocketmq.b1.large.1"
}
```

### Query the list of RocketMQ flavors available in the availability zones

```hcl
variable "az1" {}

data "opentelekomcloud_dms_rocketmq_flavors_v2" "test" {
  type               = "single.basic"
  storage_spec_code  = "dms.physical.storage.ultra.v2"
  availability_zones = [var.az1]
}
```

## Argument Reference

* `flavor_id` - (Optional, String) Specifies the flavor ID, e.g. **rocketmq.b1.large.1**.

* `type` - (Optional, String) Specifies the flavor type, e.g. **single.basic** or **cluster.basic**.

* `arch_type` - (Optional, String) Specifies the type of CPU architecture, e.g. **X86**.

* `storage_spec_code` - (Optional, String) Specifies the storage flavor code.
  + **dms.physical.storage.high.v2**: high I/O disk.
  + **dms.physical.storage.ultra.v2**: ultra-high I/O disk.
  + **dms.physical.storage.general**: general-purpose SSD.
  + **dms.physical.storage.extreme**: extreme SSD.

* `availability_zones` - (Optional, List) Specifies the list of availability zones with available resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the flavors are queried.

* `versions` - The supported engine versions.

* `flavors` - The list of flavor details.
  The [flavors](#dms_rocketmq_flavors) structure is documented below.

<a name="dms_rocketmq_flavors"></a>
The `flavors` block supports:

* `id` - The flavor ID.

* `type` - The flavor type.

* `vm_specification` - The underlying ECS flavor.

* `billing_code` - The CBC flavor code.

* `arch_types` - The list of supported CPU architectures.

* `charging_modes` - The list of supported billing modes.

* `qingtian_incompatible` - Whether the instance is a QingTian one.

* `ios` - The list of supported disk IO types.
  The [ios](#dms_rocketmq_flavor_ios) structure is documented below.

* `properties` - The properties of the flavor.
  The [properties](#dms_rocketmq_flavor_properties) structure is documented below.

<a name="dms_rocketmq_flavor_ios"></a>
The `ios` block supports:

* `storage_spec_code` - The storage flavor code.

* `type` - The service type to which storage belongs.

* `availability_zones` - The list of availability zones with available resources.

* `unavailability_zones` - The list of unavailable availability zones.

<a name="dms_rocketmq_flavor_properties"></a>
The `properties` block supports:

* `engine_versions` - The message engine version.

* `flavor_alias` - The flavor ID alias.

* `broker_num` - The broker quantity.

* `min_broker` - The minimum number of brokers.

* `max_broker` - The maximum number of brokers.

* `core` - The number of billing cores of an entire instance.

* `rcu` - The traffic unit, `rcu` x `max_tps_per_rcu` = maximum flavor TPS.

* `max_tps_per_rcu` - The maximum TPS of each RCU.

* `max_topic` - The maximum number of topics.

* `max_topic_per_broker` - The maximum number of topics that can be created on each broker.

* `max_consumer` - The maximum number of consumers in an instance.

* `max_consumer_per_broker` - The maximum number of consumers on each broker.

* `min_storage` - The minimum storage space, in GB.

* `max_storage` - The maximum storage space, in GB.

* `min_storage_per_node` - The minimum storage space of each node, in GB.

* `max_storage_per_node` - The maximum storage space of each node, in GB.
