---
subcategory: "Distributed Message Service (DMS)"
---

# opentelekomcloud_dms_instance_v1

Manages a DMS instance in the OpenTelekomCloud DMS Service (Kafka Premium/Platinum).

## Example Usage

### Automatically detect the correct network

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}

data "opentelekomcloud_dms_az_v1" "az_1" {
  name = "eu-de-01"
}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine            = "kafka"
  version           = "2.3.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}

resource "opentelekomcloud_dms_instance_v1" "instance_1" {
  name              = "kafka-test"
  engine            = "kafka"
  product_id        = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version    = data.opentelekomcloud_dms_product_v1.product_1.version
  specification     = data.opentelekomcloud_dms_product_v1.product_1.bandwidth
  partition_num     = data.opentelekomcloud_dms_product_v1.product_1.partition_num
  storage_spec_code = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
  storage_space     = data.opentelekomcloud_dms_product_v1.product_1.storage
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  access_user       = var.access_user
  password          = var.password
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Indicates the name of an instance. An instance name starts with a letter,
  consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).

* `description` - (Optional) Indicates the description of an instance. It is a character
  string containing not more than 1024 characters.

* `engine` - (Required) Indicates a message engine. Only `kafka` is supported now.

* `engine_version` - (Required) Indicates the version of a message engine.
  Only `2.3.0` is supported now.

* `specification` - (Optional) This parameter is mandatory if the engine is `kafka`.
  Indicates the baseline bandwidth of a Kafka instance, that is, the maximum amount
  of data transferred per unit time. Unit: `byte/s`. Options: `100MB`, `300MB`,
  `600MB`, `1200MB`.

* `storage_space` - (Required) Indicates the message storage space. Value range:
  * Kafka instance with `specification` being `100MB`: `600`–`90000` GB
  * Kafka instance with `specification` being `300MB`: `1200`–`90000` GB
  * Kafka instance with `specification` being `600MB`: `2400`–`90000` GB
  * Kafka instance with `specification` being `1200MB`: `4800`–`90000` GB

* `partition_num` - (Optional) This parameter is mandatory when a `kafka` instance is created.
  Indicates the maximum number of topics in a Kafka instance.
  * When `specification` is `100MB`: `300`
  * When `specification` is `300MB`: `900`
  * When `specification` is `600MB`: `1800`
  * When `specification` is `1200MB`: `1800`

* `access_user` - (Optional) Indicates a username. A username consists of 4 to 64 characters
  and supports only letters, digits, and hyphens (-).

* `password` - (Optional) Indicates the password of an instance. An instance password
  must meet the following complexity requirements: Must be 8 to 32 characters long.
  Must contain at least 2 of the following character types: lowercase letters, uppercase
  letters, digits, and special characters (`~!@#$%^&*()-_=+\|[{}]:'",<.>/?`).

* `vpc_id` - (Required) Indicates the ID of a VPC (OpenStack router ID).

* `security_group_id` - (Required) Indicates the ID of a security group.

* `subnet_id` - (Required) Indicates the ID of the subnet (OpenStack network ID).

* `available_zones` - (Required) Indicates the ID of an AZ. The parameter value can not be
  left blank or an empty array. For details, see section
  [Querying AZ Information](https://docs.otc.t-systems.com/en-us/api/dms/dms-api-180514008.html).

* `product_id` - (Required) Indicates a product ID.

* `maintain_begin` - (Optional) Indicates the time at which a maintenance time window starts.
  Format: `HH:mm`.
  * The start time and end time of a maintenance time window must indicate the time segment of
  a supported maintenance time window.
  * The start time must be set to `22:00`, `02:00`, `06:00`, `10:00`, `14:00`, or `18:00`.
  * Parameters `maintain_begin` and `maintain_end` must be set in pairs. If parameter `maintain_begin`
  is left blank, parameter `maintain_end` is also blank. In this case, the system automatically
  allocates the default start time `02:00`.

* `maintain_end` - (Optional) Indicates the time at which a maintenance time window ends.
  Format: `HH:mm`.
  * The start time and end time of a maintenance time window must indicate the time segment of
  a supported maintenance time window.
  * The end time is four hours later than the start time. For example, if the start time is `22:00`,
  the end time is `02:00`.
  * Parameters `maintain_begin` and `maintain_end` must be set in pairs. If parameter `maintain_end` is left
  blank, parameter `maintain_begin` is also blank. In this case, the system automatically allocates
  the default end time `06:00`.

* `storage_spec_code` - (Required) Indicates the storage I/O specification. Options for a Kafka instance:
  * When `specification` is `100MB`: `dms.physical.storage.high` or `dms.physical.storage.ultra`
  * When `specification` is `300MB`: `dms.physical.storage.high` or `dms.physical.storage.ultra`
  * When `specification` is `600MB`: `dms.physical.storage.ultra`
  * When `specification` is `1200MB`: `dms.physical.storage.ultra`

* `retention_policy` - (Optional) Indicates the action to be taken when the memory usage reaches
  the disk capacity threshold. The possible values are:
  * `produce_reject`: New messages cannot be created
  * `time_base`: The earliest messages are deleted.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `engine` - See Argument Reference above.

* `engine_version` - See Argument Reference above.

* `specification` - See Argument Reference above.

* `storage_space` - Indicates the time when an instance is created.

* `partition_num` - See Argument Reference above.

* `access_user` - See Argument Reference above.

* `password` - See Argument Reference above.

* `vpc_id` - See Argument Reference above.

* `security_group_id` - See Argument Reference above.

* `security_group_name` - Indicates the name of a security group.

* `subnet_id` - See Argument Reference above.

* `subnet_name` - Indicates the name of a subnet.

* `subnet_cidr` - Indicates a subnet segment.

* `available_zones` - See Argument Reference above.

* `product_id` - See Argument Reference above.

* `maintain_begin` - See Argument Reference above.

* `maintain_end` - See Argument Reference above.

* `storage_spec_code` - See Argument Reference above.

* `used_storage_space` - Indicates the used message storage space. Unit: GB

* `connect_address` - Indicates the IP address of an instance.

* `port` - Indicates the port number of an instance.

* `status` - Indicates the status of an instance. For details, see section Instance Status.

* `instance_id` - Indicates the ID of an instance.

* `resource_spec_code` - Indicates a resource specifications identifier.

* `type` - Indicates an instance type. Options: `single` and `cluster`.

* `created_at` - Indicates the time when an instance is created. The time is in the format
  of timestamp, that is, the offset milliseconds from 1970-01-01 00:00:00 UTC to the specified time.

* `user_id` - Indicates a user ID.

* `user_name` -	Indicates a username.
