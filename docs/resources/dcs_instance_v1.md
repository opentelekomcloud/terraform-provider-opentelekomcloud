---
subcategory: "Distributed Cache Service (DCS)"
---

Up-to-date reference of API arguments for DCS instance you can get at
`https://docs.otc.t-systems.com/distributed-cache-service/api-ref/lifecycle_management_apis`.

# opentelekomcloud_dcs_instance_v1

Manages a DCSv1 instance in the OpenTelekomCloud DCS Service.

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

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "dcs.master_standby"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name           = "test_dcs_instance"
  engine_version = "3.0"
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

### Engine version 5.0 (please pay attention of proper selection of the spec_code):

```hcl
data "opentelekomcloud_dcs_az_v1" "az_1" {
  name = "eu-de-01"
}

data "opentelekomcloud_dcs_product_v1" "product_1" {
  spec_code = "redis.single.xu1.tiny.128"
}

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
  name            = "test_dcs_instance_5.0"
  engine_version  = "5.0"
  password        = "0TCTestP@ssw0rd"
  engine          = "Redis"
  capacity        = 0.125
  vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones = [data.opentelekomcloud_dcs_az_v1.az_1.id]
  product_id      = data.opentelekomcloud_dcs_product_v1.product_1.id

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

* `name` - (Required, String) Indicates the name of an instance. An instance name starts with a letter,
  consists of `4` to `64` characters, and supports only letters, digits, and hyphens (-).

* `description` - (Optional, String) Indicates the description of an instance. It is a character
  string containing not more than `1024` characters.

* `engine` - (Required, ForceNew, String) Indicates a cache engine. Only `Redis` is supported. Changing this
  creates a new instance.

* `engine_version` - (Required, ForceNew, String) Indicates the version of a cache engine, which can be `3.0`/`4.0`/`5.0`/`6.0`.
  Changing this creates a new instance.

* `capacity` - (Required, ForceNew, Float) Indicates the Cache capacity. Unit: GB.
  + **Redis4.0, Redis5.0 and Redis6.0**: Stand-alone and active/standby type instance values: `0.125`, `0.25`,
    `0.5`, `1`, `2`, `4`, `8`, `16`, `32` and `64`.
    Cluster instance specifications support `4`,`8`,`16`, `24`, `32`, `48`, `64`, `96`, `128`, `192`, `256`,
    `384`, `512`, `768` and `1024`.
  + **Redis3.0**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.
    Proxy cluster instance specifications support `64`, `128`, `256`, `512`, and `1024`.
  + **Memcached**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.

* `password` - (Optional, ForceNew, String) Indicates the password of an instance. An instance password
  must meet the following complexity requirements: Must be 8 to 32 characters long.
  Must contain at least 3 of the following character types: lowercase letters, uppercase
  letters, digits, and special characters: `~!@#$^&*()-_=+|{}:,<>./?
  Changing this creates a new instance.

* `vpc_id` - (Required, ForceNew, String) Specifies the VPC ID. Changing this creates a new instance.

* `security_group_id` - (Optional, String) Security group ID. This parameter is mandatory when `engine_version` is `3.0`.

* `subnet_id` - (Required, ForceNew, String) Specifies the subnet Network ID. Changing this creates a new instance.

* `available_zones` - (Required, ForceNew, List) IDs of the AZs where cache nodes reside. For details
  on how to query AZs, see [Querying AZ Information](https://docs.otc.t-systems.com/en-us/api/dcs/dcs-api-0312039.html)
  or use [opentelekomcloud_dcs_az_v1 data source](https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs/data-sources/dcs_az_v1):
  ```hcl
  data "opentelekomcloud_dcs_az_v1" "az1" {
    name = "eu-de-01"
  }
  ```
  Changing this creates a new instance.

* `product_id` - (Required, ForceNew, String) Product ID used to differentiate DCS instance types.
  Changing this creates a new instance.

* `maintain_begin` - (Optional, String) Indicates the time at which a maintenance time window starts.
  Format: `HH:mm:ss`. The start time and end time of a maintenance time window must indicate the time segment of
  a supported maintenance time window. For details, see section
  [Querying Maintenance Time Windows](https://docs.otc.t-systems.com/api/dcs/dcs-api-0312041.html).
  The start time must be set to `22:00`, `02:00`, `06:00`, `10:00`, `14:00`, or `18:00`.

* `maintain_end` - (Optional, String) Indicates the time at which a maintenance time window ends.
  Format: `HH:mm:ss`. The start time and end time of a maintenance time window must indicate the time segment of
  a supported maintenance time window. For details, see section
  [Querying Maintenance Time Windows](https://docs.otc.t-systems.com/api/dcs/dcs-api-0312041.html).
  The end time is four hours later than the start time. For example, if the start time is `22:00`,
  the end time is `02:00`.

-> Parameters `maintain_begin` and `maintain_end` must be set in pairs. If parameter `maintain_end` is left
blank, parameter `maintain_begin` is also blank. In this case, the system automatically allocates
the default start time `02:00` and the default end time `06:00`.

* `backup_policy` - (Optional, List) Describes the backup configuration to be used with the instance.
  * `save_days` - (Optional, Int) Retention time. Unit: day. Range: `1`–`7`.
  * `backup_type` - (Optional, String) Backup type. Valid values are: `auto` automatic backup,
    `manual` manual backup (default).
  * `begin_at` - (Required, String) Time at which backup starts. `00:00-01:00` indicates that backup
    starts at `00:00:00`.
  * `period_type` - (Required, String) Interval at which backup is performed.
    Currently, only weekly backup is supported.
  * `backup_at` - (Required, List) Day in a week on which backup starts. Range: `1`–`7`. Where: `1`
    indicates Monday; `7` indicates Sunday.

* `configuration` - (Optional, List) Describes the array of configuration items of the DCS instance.
  Configured values can be found [here](https://docs.otc.t-systems.com/en-us/api/dcs/dcs-api-0312015.html#dcs-api-0312015__table1439111281351).
  * `parameter_id` - (Required, String) Configuration item ID.
  * `parameter_name` - (Required, String) Configuration item name.
  * `parameter_value` - (Required, String) Value of the configuration item.

* `enable_whitelist` - (Optional, Bool) Specifies whether to enable or disable `whitelist`. Only available when
  `engine_version` is set to `4.0`/`5.0`. Parameter have to be used together with `whitelist`.

* `whitelist` - (Optional, List) Describes the `whitelist` groups to be used with the instance. Only available when
  `engine_version` is set to `4.0`/`5.0`. Parameter have to be used together with `enable_whitelist`.
  Resource fields:
  * `group_name` - (Required, String) Whitelist group name. A maximum of four groups can be created for each instance.
  * `ip_list` - (Required, List) List of IP addresses in the whitelist group. A maximum of 20 IP addresses or IP address
  ranges can be added to an instance. Separate multiple IP addresses or IP address ranges with commas (,).
  IP address 0.0.0.0 and IP address range 0.0.0/0 are not supported.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `engine` - See Argument Reference above.

* `engine_version` - See Argument Reference above.

* `capacity` - See Argument Reference above.

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
  * `dcs.cluster`: indicates a DCS instance in cluster mode. Not available with version `6.0`.

* `used_memory` - Size of the used memory. Unit: MB.

* `internal_version` - Internal DCS version.

* `max_memory` - Overall memory size. Unit: MB.

* `user_id` - Indicates a user ID.

* `user_name` - Username.

* `ip` - Cache node's IP address in the tenant's VPC.

* `status` - Cache instance status. One of `CREATING`, `CREATEFAILED`, `RUNNING`, `ERROR`,
  `RESTARTING`, `EXTENDING`, `RESTORING`

* `created_at` - Time at which the DCS instance is created. For example, `2017-03-31T12:24:46.297Z`.

* `no_password_access` - An indicator of whether a DCS instance can be accessed in password-free mode.
  `true` when password not set.

## Import

DCS instance can be imported using  `instance_name`, e.g.
```shell
$ terraform import opentelekomcloud_dcs_instance_v1.instance instance_name
```
