---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_instance_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-instance-v2"
description: |-
  Manages a DMS RocketMQ instance resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ instance you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/lifecycle_management/index.html)

# opentelekomcloud_dms_rocketmq_instance_v2

Manages a DMS RocketMQ instance in the OpenTelekomCloud DMS Service.

## Example Usage

### Create a single-node RocketMQ instance

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

data "opentelekomcloud_dms_az_v1" "az_1" {}

resource "opentelekomcloud_dms_rocketmq_instance_v2" "test" {
  name        = "rocketmq_test"
  description = "rocketmq test"

  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id

  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  engine_version    = "5.x"
  flavor_id         = "rocketmq.b1.large.1"
  storage_spec_code = "dms.physical.storage.ultra.v2"
  storage_space     = 300
  broker_num        = 1

  configs {
    name  = "fileReservedTime"
    value = "72"
  }

  tags = {
    foo = "bar"
  }
}
```

### Create a RocketMQ instance with public access

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

data "opentelekomcloud_dms_az_v1" "az_1" {}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_dms_rocketmq_instance_v2" "test" {
  name = "rocketmq_public"

  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id

  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  engine_version    = "5.x"
  flavor_id         = "rocketmq.b1.large.1"
  storage_spec_code = "dms.physical.storage.ultra.v2"
  storage_space     = 300
  broker_num        = 1

  enable_publicip = true
  publicip_id     = [opentelekomcloud_networking_floatingip_v2.fip_1.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the DMS RocketMQ instance. An instance name starts with a letter,
  consists of 4 to 64 characters, and can contain only letters, digits, underscores (_) and hyphens (-).

* `description` - (Optional, String) Specifies the description of the DMS RocketMQ instance.
  It is a character string containing not more than 1,024 characters.

* `engine_version` - (Required, String, ForceNew) Specifies the version of the RocketMQ engine. Value: **5.x**.
  Changing this creates a new instance resource.

* `flavor_id` - (Required, String) Specifies the RocketMQ [flavor ID](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/other_apis/querying_flavor_list.html),
  e.g. **rocketmq.b1.large.1** (single-node) or **rocketmq.b2.large.4** (cluster).
  The flavor of a cluster instance can be changed, the flavor of a single-node instance can't.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification.
  Value options: **dms.physical.storage.high.v2**, **dms.physical.storage.ultra.v2**,
  **dms.physical.storage.general** and **dms.physical.storage.extreme**.
  Changing this creates a new instance resource.

* `storage_space` - (Required, Int) Specifies the message storage capacity, the unit is GB.
  The supported range depends on the flavor, see `min_storage_per_node` and `max_storage_per_node` of the flavor.
  The storage space can only be expanded, the storage space of each broker must be expanded by at least 100 GB.

* `broker_num` - (Required, Int) Specifies the number of brokers.
  The broker number can't be changed separately, the instance can only be scaled via `flavor_id`.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC. Changing this creates a new instance resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet.
  Changing this creates a new instance resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `available_zones` - (Required, Set, ForceNew) Specifies the IDs of the AZs where the brokers reside.
  A RocketMQ instance can be deployed in 1 AZ or at least 3 AZs.
  Changing this creates a new instance resource.

* `ssl_enable` - (Optional, Bool, ForceNew) Specifies whether to enable SSL-encrypted access. Defaults to **false**.
  Changing this creates a new instance resource.

* `tls_mode` - (Optional, String, ForceNew) Specifies the security protocol used by the instance.
  Changing this creates a new instance resource.

* `arch_type` - (Optional, String, ForceNew) Specifies the CPU architecture. Value options: **X86**, **ARM**.
  Changing this creates a new instance resource.

* `enable_acl` - (Optional, Bool) Specifies whether access control is enabled.

* `enable_publicip` - (Optional, Bool) Specifies whether to enable public access. By default, public access is disabled.

* `publicip_id` - (Optional, Set) Specifies the IDs of the EIPs bound to the instance.
  This parameter is mandatory if public access is enabled.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the instance.

* `configs` - (Optional, Set) Specifies the RocketMQ configurations of the instance.
  The [configs](#dms_rocketmq_configs) structure is documented below.

* `tags` - (Optional, Map) The key/value pairs to associate with the instance. A maximum of 20 tags can be added.

<a name="dms_rocketmq_configs"></a>
The `configs` block supports:

* `name` - (Required, String) Specifies the configuration name, e.g. **fileReservedTime**.

* `value` - (Required, String) Specifies the configuration value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `region` - The region in which the DMS RocketMQ instance is created.
* `engine` - Indicates the message engine.
* `status` - Indicates the status of the DMS RocketMQ instance.
* `type` - Indicates the instance type. Values: **single.basic**, **cluster.basic**.
* `specification` - Indicates the instance specification.
* `maintain_begin` - Indicates the time at which the maintenance window starts. The format is HH:mm:ss.
* `maintain_end` - Indicates the time at which the maintenance window ends. The format is HH:mm:ss.
* `used_storage_space` - Indicates the used message storage space. Unit: GB.
* `publicip_address` - Indicates the public IP address.
* `ipv6_enable` - Indicates whether IPv6 is enabled.
* `namesrv_address` - Indicates the metadata address.
* `broker_address` - Indicates the service data address.
* `public_namesrv_address` - Indicates the public network metadata address.
* `public_broker_address` - Indicates the public network service data address.
* `grpc_address` - Indicates the gRPC connection address.
* `public_grpc_address` - Indicates the public gRPC connection address.
* `resource_spec_code` - Indicates the resource specification.
* `created_at` - Indicates the creation time, as a timestamp in milliseconds.
* `cross_vpc_accesses` - Indicates the cross-VPC access information.
  The [cross_vpc_accesses](#dms_rocketmq_cross_vpc_accesses) structure is documented below.

<a name="dms_rocketmq_cross_vpc_accesses"></a>
The `cross_vpc_accesses` block supports:

* `advertised_ip` - The advertised IP address or domain name.
* `listener_ip` - The listener IP address.
* `port` - The port number.
* `port_id` - The port ID associated with the address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minutes.
* `update` - Default is 50 minutes.
* `delete` - Default is 15 minutes.

## Import

DMS RocketMQ instance can be imported using the instance id, e.g.

```
 $ terraform import opentelekomcloud_dms_rocketmq_instance_v2.test 8d3c7938-dc47-4937-a30f-c80de381c5e3
```

Note that the imported state may not be identical to your resource definition, as only the `configs`
listed in the resource definition are tracked. It is generally recommended running `terraform plan` after importing
a DMS RocketMQ instance. You can then decide if changes should be applied to the instance, or the resource definition
should be updated to align with the instance. Also, you can ignore changes as below.

```hcl
resource "opentelekomcloud_dms_rocketmq_instance_v2" "test" {
  lifecycle {
    ignore_changes = [
      configs,
    ]
  }
}
```
