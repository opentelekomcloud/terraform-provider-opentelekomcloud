---
subcategory: "Distributed Cache Service (DCS)"
---

# opentelekomcloud_dcs_instance_v1

Manages a DCSv1 instance in the OpenTelekomCloud DCS Service.

## Example Usage

```hcl
variable "network_id" {}
variable "vpc_id" {}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}

data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "dcs.master_standby"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name           = "test_dcs_instance"
  engine_version = "3.0.7"
  password       = "0TCTestP@ssw0rd"
  engine         = "Redis"
  capacity       = 2
  vpc_id         = var.vpc_id

  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
  subnet_id         = var.network_id
  available_zones   = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dcs_product_v1.product_1.id
  backup_policy {
    save_days   = 1
    backup_type = "manual"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1, 2, 4, 6]
  }
}
```

### Engine version 5.0:

```hcl
data "opentelekomcloud_dcs_az_v1" "az_1" {
  port = "8002"
  code = "eu-de"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.single.xu1.tiny.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name              = "test_dcs_instance_5.0"
  engine_version    = "5.0"
  password          = "0TCTestP@ssw0rd"
  engine            = "Redis"
  capacity          = 0.125
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones   = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dcs_product_v1.product_1.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Indicates the name of an instance. An instance name starts with a letter,
  consists of `4` to `64` characters, and supports only letters, digits, and hyphens (-).

* `description` - (Optional) Indicates the description of an instance. It is a character
  string containing not more than `1024` characters.

* `engine` - (Required) Indicates a cache engine. Only `Redis` is supported. Changing this
  creates a new instance.

* `engine_version` - (Required) Indicates the version of a cache engine, which is `3.0.7`.
  Changing this creates a new instance.

* `capacity` - (Required) Indicates the Cache capacity. Unit: GB.
  For a DCS Redis or Memcached instance in single-node or master/standby mode, the cache
  capacity can be `2`, `4`, `8`, `16`, `32`, or `64` GB.
  For a DCS Redis instance in cluster mode, the cache capacity can be `64`, `128`, `256`, `512` GB.
  Changing this creates a new instance.

* `access_user` - (Optional) Username used for accessing a DCS instance after password
  authentication. A username starts with a letter, consists of `1` to `64` characters,
  and supports only letters, digits, and hyphens (-). Changing this creates a new instance.

* `password` - (Required) Indicates the password of an instance. An instance password
  must meet the following complexity requirements: Must be 8 to 32 characters long.
  Must contain at least 3 of the following character types: lowercase letters, uppercase
  letters, digits, and special characters (`~!@#$%^&*()-_=+\|[{}]:'",<.>/?).
  Changing this creates a new instance.

* `vpc_id` - (Required) Specifies the VPC ID. Changing this creates a new instance.

* `security_group_id` - (Required) Security group ID.

* `subnet_id` - (Required) Specifies the subnet Network ID. Changing this creates a new instance.

* `available_zones` - (Required) IDs of the AZs where cache nodes reside. For details
  on how to query AZs, see [Querying AZ Information](https://docs.otc.t-systems.com/en-us/api/dcs/dcs-api-0312039.html).
  Changing this creates a new instance.

* `product_id` - (Required) Product ID used to differentiate DCS instance types.
  Changing this creates a new instance.

* `maintain_begin` - (Optional) Indicates the time at which a maintenance time window starts.
  Format: `HH:mm:ss`. The start time and end time of a maintenance time window must indicate the time segment of
  a supported maintenance time window. For details, see section
  [Querying Maintenance Time Windows](https://docs.otc.t-systems.com/api/dcs/dcs-api-0312041.html).
  The start time must be set to `22:00`, `02:00`, `06:00`, `10:00`, `14:00`, or `18:00`.

* `maintain_end` - (Optional) Indicates the time at which a maintenance time window ends.
  Format: `HH:mm:ss`. The start time and end time of a maintenance time window must indicate the time segment of
  a supported maintenance time window. For details, see section
  [Querying Maintenance Time Windows](https://docs.otc.t-systems.com/api/dcs/dcs-api-0312041.html).
  The end time is four hours later than the start time. For example, if the start time is `22:00`,
  the end time is `02:00`.

-> Parameters `maintain_begin` and `maintain_end` must be set in pairs. If parameter `maintain_end` is left
blank, parameter `maintain_begin` is also blank. In this case, the system automatically allocates
the default start time `02:00` and the default end time `06:00`.

* `backup_policy` - (Optional) Describes the backup configuration to be used with the instance.
  * `save_days` - (Optional) Retention time. Unit: day. Range: `1`–`7`.
  * `backup_type` - (Optional) Backup type. Valid values are: `auto` automatic backup,
    `manual` manual backup (default).
  * `begin_at` - (Required) Time at which backup starts. `00:00-01:00` indicates that backup
    starts at `00:00:00`.
  * `period_type` - (Required) Interval at which backup is performed.
    Currently, only weekly backup is supported.
  * `backup_at` - (Required) Day in a week on which backup starts. Range: `1`–`7`. Where: `1`
    indicates Monday; `7` indicates Sunday.

* `configuration` - (Optional) Describes the array of configuration items of the DCS instance.
  Configured values can be found [here](https://docs.otc.t-systems.com/en-us/api/dcs/dcs-api-0312015.html#dcs-api-0312015__table1439111281351).
  * `parameter_id` - (Required) Configuration item ID.
  * `parameter_name` - (Required) Configuration item name.
  * `parameter_value` - (Required) Value of the configuration item.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `engine` - See Argument Reference above.

* `engine_version` - See Argument Reference above.

* `capacity` - See Argument Reference above.

* `access_user` - See Argument Reference above.

* `password` - See Argument Reference above.

* `vpc_id` - See Argument Reference above.

* `vpc_name` - Indicates the name of a vpc.

* `security_group_id` - See Argument Reference above.

* `security_group_name` - Indicates the name of a security group.

* `subnet_id` - See Argument Reference above.

* `subnet_name` - Indicates the name of a subnet.

* `available_zones` - See Argument Reference above.

* `product_id` - See Argument Reference above.

* `maintain_begin` - See Argument Reference above.

* `maintain_end` - See Argument Reference above.

* `save_days` - See Argument Reference above.

* `backup_type` - See Argument Reference above.

* `begin_at` - See Argument Reference above.

* `period_type` - See Argument Reference above.

* `backup_at` - See Argument Reference above.

* `backup_policy` - See Argument Reference above.

* `order_id` - An order ID is generated only in the monthly or yearly billing mode.
  In other billing modes, no value is returned for this parameter.

* `port` - Port of the cache node.

* `resource_spec_code` - Resource specifications.
  * `dcs.single_node`: indicates a DCS instance in single-node mode.
  * `dcs.master_standby`: indicates a DCS instance in master/standby mode.
  * `dcs.cluster`: indicates a DCS instance in cluster mode.

* `used_memory` - Size of the used memory. Unit: MB.

* `internal_version` - Internal DCS version.

* `max_memory` - Overall memory size. Unit: MB.

* `user_id` - Indicates a user ID.

* `user_name` - Username.

* `ip` - Cache node's IP address in the tenant's VPC.

* `status` - Cache instance status. One of `CREATING`, `CREATEFAILED`, `RUNNING`, `ERROR`,
  `RESTARTING`, `EXTENDING`, `RESTORING`

* `created_at` - Time at which the DCS instance is created. For example, `2017-03-31T12:24:46.297Z`.
