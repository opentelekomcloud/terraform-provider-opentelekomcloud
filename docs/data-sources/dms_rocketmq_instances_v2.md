---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_instances_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-instances-v2"
description: |-
  Get the list of DMS RocketMQ instances within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ instances you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/lifecycle_management/querying_the_instance_list.html)

# opentelekomcloud_dms_rocketmq_instances_v2

Use this data source to get the list of DMS RocketMQ instances within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_name" {}

data "opentelekomcloud_dms_rocketmq_instances_v2" "test" {
  name             = var.instance_name
  exact_match_name = true
}
```

## Argument Reference

* `name` - (Optional, String) Specifies the instance name.

* `exact_match_name` - (Optional, Bool) Specifies whether to search for the instance that precisely matches
  the specified `name`. Defaults to **false**.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `status` - (Optional, String) Specifies the instance status, e.g. **RUNNING**.

* `include_failure` - (Optional, Bool) Specifies whether to return the instances that fail to be created.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the instances are queried.

* `instances` - The list of instances.
  The [instances](#dms_rocketmq_instances) structure is documented below.

<a name="dms_rocketmq_instances"></a>
The `instances` block supports:

* `id` - The instance ID.
* `name` - The instance name.
* `description` - The instance description.
* `engine` - The message engine.
* `engine_version` - The engine version.
* `status` - The instance status.
* `type` - The instance type. Values: **single.basic**, **cluster.basic**.
* `specification` - The instance specification.
* `flavor_id` - The flavor ID.
* `arch_type` - The CPU architecture type.
* `broker_num` - The number of brokers.
* `storage_space` - The total message storage space, in GB.
* `used_storage_space` - The used message storage space, in GB.
* `storage_spec_code` - The storage specification code.
* `disk_encrypted` - Whether disk encryption is enabled.
* `vpc_id` - The VPC ID.
* `vpc_name` - The VPC name.
* `subnet_id` - The subnet ID.
* `subnet_cidr` - The subnet CIDR.
* `security_group_id` - The security group ID.
* `security_group_name` - The security group name.
* `available_zones` - The list of availability zone IDs.
* `available_zone_names` - The list of availability zone names.
* `maintain_begin` - The time at which the maintenance window starts. The format is HH:mm:ss.
* `maintain_end` - The time at which the maintenance window ends. The format is HH:mm:ss.
* `ssl_enable` - Whether SSL is enabled.
* `tls_mode` - The security protocol used by the instance.
* `enable_acl` - Whether access control is enabled.
* `dns_enable` - Whether domain name access to the instance is enabled.
* `enable_publicip` - Whether public access is enabled.
* `publicip_id` - The IDs of the EIPs bound to the instance.
* `publicip_address` - The public IP address.
* `namesrv_address` - The metadata address.
* `namesrv_domain_name` - The metadata domain name.
* `broker_address` - The service data address.
* `public_namesrv_address` - The public network metadata address.
* `public_namesrv_domain_name` - The public metadata domain name.
* `public_broker_address` - The public network service data address.
* `grpc_address` - The gRPC connection address.
* `grpc_domain_name` - The gRPC connection domain name.
* `public_grpc_address` - The public gRPC connection address.
* `public_grpc_domain_name` - The public gRPC domain name.
* `cross_vpc_accesses` - The cross-VPC access information.
  The [cross_vpc_accesses](#dms_rocketmq_instances_cross_vpc_accesses) structure is documented below.
* `resource_spec_code` - The resource specification.
* `enterprise_project_id` - The enterprise project ID.
* `tags` - The key/value pairs associated with the instance.
* `created_at` - The creation time, as a timestamp in milliseconds.

<a name="dms_rocketmq_instances_cross_vpc_accesses"></a>
The `cross_vpc_accesses` block supports:

* `advertised_ip` - The advertised IP address or domain name.
* `listener_ip` - The listener IP address.
* `port` - The port number.
* `port_id` - The port ID associated with the address.
