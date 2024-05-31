---
subcategory: "Distributed Cache Service (DCS)"
---

Up-to-date reference of API arguments for DCS V2 instance you can get at
`https://docs.otc.t-systems.com/distributed-cache-service/api-ref/apis_v2_recommended/index.html`.

# opentelekomcloud_dcs_instance_v1

Manages a DCSv2 instance in the OpenTelekomCloud DCS Service.

## Example Usage

### Engine version 3.0 (`security_group_id` is required):

```hcl
variable "network_id" {}
variable "vpc_id" {}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}

data "opentelekomcloud_dcs_az_v1" "az_1" {
  name = "eu-de-01"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name               = "test_dcs_instance"
  engine_version     = "3.0"
  password           = "0TCTestP@ssw0rd"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = var.vpc_id
  security_group_id  = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
  subnet_id          = var.network_id
  availability_zones = ["eu-de-01"]
  flavor             = "dcs.master_standby"

  backup_policy {
    save_days   = 1
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1, 2, 4, 6]
  }
  tags = {
    environment = "basic"
    managed_by  = "terraform"
  }
}
```

### Engine version 5.0 (please pay attention of proper selection of flavor):

```hcl
data "opentelekomcloud_compute_availability_zones_v2" "zones" {}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name               = "test_dcs_instance_5.0"
  engine_version     = "5.0"
  password           = "0TCTestP@ssw0rd"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  availability_zones = [data.opentelekomcloud_compute_availability_zones_v2.zones.names[0]]
  flavor             = "redis.single.xu1.tiny.128"

  enable_whitelist = true
  whitelist {
    group_name = "test-group-name"
    ip_list    = ["10.10.10.1", "10.10.10.2"]
  }
  whitelist {
    group_name = "test-group-name-2"
    ip_list    = ["10.10.10.11", "10.10.10.3", "10.10.10.4"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of an instance.
  The name must be 4 to 64 characters and start with a letter.
  Only chinese, letters (case-insensitive), digits, underscores (_) ,and hyphens (-) are allowed.

* `engine` - (Required, String, ForceNew) Specifies a cache engine. Options: *Redis* and *Memcached*.
  Changing this creates a new instance.

* `engine_version` - (Optional, String, ForceNew) Specifies the version of a cache engine.
  It is mandatory when the engine is *Redis*, the value can be 3.0, 4.0, 5.0 or 6.0.
  Changing this creates a new instance.

* `capacity` - (Required, Float) Specifies the cache capacity. Unit: GB.

* `flavor` - (Required, String) The flavor of the cache instance, which including the total memory, available memory,
  maximum number of connections allowed, maximum/assured bandwidth and reference performance.
  It also includes the modes of Redis instances. You can query the *flavor* as follows:
    + Query flavors
      in [DCS Instance Specifications](https://docs.otc.t-systems.com/distributed-cache-service/umn/service_overview/dcs_instance_specifications/index.html)
    + Log in to the DCS console, click *Create DCS Instance*, and find the corresponding instance specification.

* `availability_zones` - (Optional, List, ForceNew) The code of the AZ where the cache node resides.
  Master/Standby, Proxy Cluster, and Redis Cluster DCS instances support cross-AZ deployment.
  You can specify an AZ for the standby node. When specifying AZs for nodes, use commas (,) to separate AZs.
  Changing this creates a new instance.

* `vpc_id` - (Required, String, ForceNew) The ID of VPC which the instance belongs to.
  Changing this creates a new instance resource.

* `subnet_id` - (Required, String, ForceNew) The ID of subnet which the instance belongs to.
  Changing this creates a new instance resource.

* `security_group_id` - (Optional, String) The ID of the security group which the instance belongs to.
  This parameter is mandatory for Memcached and Redis 3.0 version.

* `ssl_enable` - (Optional, Bool) Specifies whether to enable the SSL. Value options: **true**, **false**.

* `private_ip` - (Optional, String, ForceNew) The IP address of the DCS instance,
  which can only be the currently available IP address the selected subnet.
  You can specify an available IP for the Redis instance (except for the Redis Cluster type).
  If omitted, the system will automatically allocate an available IP address to the Redis instance.
  Changing this creates a new instance resource.

* `template_id` - (Optional, String, ForceNew) The Parameter Template ID.
  Changing this creates a new instance resource.

* `port` - (Optional, Int) Port customization, which is supported only by Redis 4.0 and Redis 5.0 instances.
  Redis instance defaults to 6379. Memcached instance does not use this argument.

* `password` - (Optional, String) Specifies the password of a DCS instance.
  The password of a DCS instance must meet the following complexity requirements:
    + Must be a string of 8 to 32 bits in length.
    + Must contain three combinations of the following four characters: Lower case letters, uppercase letter, digital,
      Special characters include (`~!@#$^&*()-_=+\\|{}:,<.>/?).
    + The new password cannot be the same as the old password.

* `whitelist` - (Optional, List) Specifies the IP addresses which can access the instance.
  This parameter is valid for Redis 4.0 and 5.0 versions. The structure is described below.

* `enable_whitelist` - (Optional, Bool) Enable or disable the IP address whitelists. Defaults to true.
  If the whitelist is disabled, all IP addresses connected to the VPC can access the instance.

* `maintain_begin` - (Optional, String) Time at which the maintenance time window starts. Defaults to **02:00:00**.
    + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
      time window.
    + The start time must be on the hour, such as **18:00:00**.
    + If parameter `maintain_begin` is left blank, parameter `maintain_end` is also blank.
      In this case, the system automatically allocates the default start time **02:00:00**.

* `maintain_end` - (Optional, String) Time at which the maintenance time window ends. Defaults to **06:00:00**.
    + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
      time window.
    + The end time is one hour later than the start time. For example, if the start time is **18:00:00**, the end time is
      **19:00:00**.
    + If parameter `maintain_end` is left blank, parameter `maintain_begin` is also blank.
      In this case, the system automatically allocates the default end time **06:00:00**.

-> **NOTE:** Parameters `maintain_begin` and `maintain_end` must be set in pairs.

* `backup_policy` - (Optional, List) Specifies the backup configuration to be used with the instance.
  The structure is described below.

  -> **NOTE:** This parameter is not supported when the instance type is single.

* `parameters` - (Optional, List) Specify an array of one or more parameters to be set to the DCS instance after
  launched. You can check on console to see which parameters supported.
  The [parameters](#DcsInstance_Parameters) structure is documented below.

* `rename_commands` - (Optional, Map) Critical command renaming, which is supported only by Redis 4.0 and
  Redis 5.0 instances but not by Redis 3.0 instance.
  The valid commands that can be renamed are: **command**, **keys**, **flushdb**, **flushall** and **hgetall**.

* `tags` - (Optional, Map) The key/value pairs to associate with the dcs instance.

* `access_user` - (Optional, String, ForceNew) Specifies the username used for accessing a DCS instance.
  The username starts with a letter, consists of 1 to 64 characters, and supports only letters, digits, and
  hyphens (-). Changing this creates a new instance.

* `description` - (Optional, String) Specifies the description of an instance.
  It is a string that contains a maximum of 1024 characters.

* `deleted_nodes` - (Optional, List) Specifies the ID of the replica to delete. This parameter is mandatory when
  you delete replicas of a master/standby DCS Redis 4.0 or 5.0 instance. Currently, only one replica can be deleted
  at a time.

* `reserved_ips` - (Optional, List) Specifies IP addresses to retain. Mandatory during cluster scale-in. If this
  parameter is not set, the system randomly deletes unnecessary shards.

The `whitelist` block supports:

* `group_name` - (Required, String) Specifies the name of IP address group.

* `ip_list` - (Required, List) Specifies the list of IP address or CIDR which can be whitelisted for an instance.
  The maximum is 20.

The `backup_policy` block supports:

* `backup_type` - (Optional, String) Backup type. Default value is `auto`. The valid values are as follows:
    + `auto`: automatic backup.
    + `manual`: manual backup.

* `save_days` - (Optional, Int) Retention time. Unit: day, the value ranges from 1 to 7.
  This parameter is required if the backup_type is **auto**.

* `period_type` - (Optional, String) Interval at which backup is performed. Default value is `weekly`.
  Currently, only weekly backup is supported.

* `backup_at` - (Required, List) Day in a week on which backup starts, the value ranges from 1 to 7.
  Where: 1 indicates Monday; 7 indicates Sunday.

* `begin_at` - (Required, String) Time at which backup starts.
  Format: `hh24:00-hh24:00`, "00:00-01:00" indicates that backup starts at 00:00:00.

<a name="DcsInstance_Parameters"></a>
The `parameters` block supports:

* `id` - (Required, String) Specifies the ID of the configuration item.

* `name` - (Required, String) Specifies the name of the configuration item.

* `value` - (Required, String) Specifies the value of the configuration item.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `status` - Cache instance status. The valid values are as follows:
    + `RUNNING`: The instance is running properly.
      Only instances in the Running state can provide in-memory cache service.
    + `ERROR`: The instance is not running properly.
    + `RESTARTING`: The instance is being restarted.
    + `FROZEN`: The instance has been frozen due to low balance.
      You can unfreeze the instance by recharging your account in My Order.
    + `EXTENDING`: The instance is being scaled up.
    + `RESTORING`: The instance data is being restored.
    + `FLUSHING`: The DCS instance is being cleared.

* `domain_name` - Domain name of the instance. Usually, we use domain name and port to connect to the DCS instances.

* `max_memory` - Total memory size. Unit: MB.

* `used_memory` - Size of the used memory. Unit: MB.

* `vpc_name` - The name of VPC which the instance belongs to.

* `subnet_name` - The name of subnet which the instance belongs to.

* `subnet_cidr` - Indicates the subnet segment.

* `security_group_name` - The name of security group which the instance belongs to.

* `created_at` - Indicates the time when the instance is created, in RFC3339 format.

* `launched_at` - Indicates the time when the instance is started, in RFC3339 format.

* `bandwidth_info` - Indicates the bandwidth information of the instance.
  The [bandwidth_info](#dcs_bandwidth_info) structure is documented below.

* `cache_mode` - Indicates the instance type. The value can be **single**, **ha**, **cluster** or **proxy**.

* `cpu_type` - Indicates the CPU type of the instance. The value can be **x86_64** or **aarch64**.

* `readonly_domain_name` - Indicates the read-only domain name of the instance. This parameter is available
  only for master/standby instances.

* `replica_count` - Indicates the number of replicas in the instance.

* `transparent_client_ip_enable` - Indicates whether client IP pass-through is enabled.

* `product_type` - Indicates the product type of the instance. The value can be: **generic** or **enterprise**.

* `sharding_count` - Indicates the number of shards in a cluster instance.

* `region` - Indicates the region in which DCS instance resource is created.

<a name="dcs_bandwidth_info"></a>
The `bandwidth_info` block supports:

* `bandwidth` - Indicates the bandwidth size, the unit is **GB**.

* `begin_time` - Indicates the begin time of temporary increase.

* `current_time` - Indicates the current time.

* `end_time` - Indicates the end time of temporary increase.

* `expand_count` - Indicates the number of increases.

* `expand_effect_time` - Indicates the interval between temporary increases, the unit is **ms**.

* `expand_interval_time` - Indicates the time interval to the next increase, the unit is **ms**.

* `max_expand_count` - Indicates the maximum number of increases.

* `next_expand_time` - Indicates the next increase time.

* `task_running` - Indicates whether the increase task is running.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 120 minutes.
* `update` - Default is 120 minutes.
* `delete` - Default is 15 minutes.

## Import

DCS instance can be imported using the `id`, e.g.

```bash
terraform import opentelekomcloud_dcs_instance_v2.instance_1 80e373f9-872e-4046-aae9-ccd9ddc55511
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `backup_policy`, `parameters`, `password`,
`bandwidth_info.0.current_time`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also, you can ignore changes as below.

```
resource "opentelekomcloud_dcs_instance_v2" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      password, rename_commands, backup_policy, parameters,
      bandwidth_info.0.current_time
    ]
  }
}
```
