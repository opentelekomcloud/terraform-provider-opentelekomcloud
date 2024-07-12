---
subcategory: "Document Database Service (DDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dds_instance_v3"
sidebar_current: "docs-opentelekomcloud-resource-dds-instance-v3"
description: |-
  Manages a DDS Instance resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DDS instance you can get at
[documentation portal](https://docs.otc.t-systems.com/document-database-service/api-ref/apis_v3.0_recommended/db_instance_management)

# opentelekomcloud_dds_instance_v3

Manages DDS instance resource within OpenTelekomCloud

## Example Usage: Creating a Replica Set
```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "opentelekomcloud_dds_instance_v3" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  availability_zone = var.availability_zone
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  password          = "5ecuredPa55w0rd@"
  mode              = "ReplicaSet"
  flavor {
    type      = "replica"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.s2.medium.4.repset"
  }
  tags = {
    foo      = "bar"
    new_test = "new_test2"
  }
}
```

## Example Usage: Creating a Cluster Community Edition
```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "opentelekomcloud_dds_instance_v3" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  availability_zone = var.availability_zone
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  password          = "5ecuredPa55w0rd2@"
  mode              = "Sharding"
  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s2.medium.4.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s2.medium.4.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s2.large.2.config"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }
}
```

## Example Usage: Creating a Single node instance
```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "opentelekomcloud_dds_instance_v3" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  availability_zone = var.availability_zone
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  password          = "5ecuredPa55w0rd@"
  mode              = "Single"
  flavor {
    type      = "single"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.s2.medium.4.single"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Specifies the region of the DDS instance. Changing this creates
  a new instance.

* `name` - (Required) Specifies the DB instance name. The DB instance name of the same
  type is unique in the same tenant.

* `datastore` - (Required) Specifies database information. The structure is described
  below. Changing this creates a new instance.

* `availability_zone` - (Required) Specifies the ID of the availability zone. Changing
  this creates a new instance.

* `vpc_id` - (Required) Specifies the VPC ID. Changing this creates a new instance.

* `subnet_id` - (Required) Specifies the subnet Network ID. Changing this creates a new instance.

* `security_group_id` - (Required) Specifies the security group ID of the DDS instance.

* `password` - (Required) Specifies the Administrator password of the database instance.

* `disk_encryption_id` - (Required) Specifies the disk encryption ID of the instance.
  Changing this creates a new instance.

* `mode` - (Required) Specifies the mode of the database instance. Changing this creates
  a new instance.

* `flavor` - (Required) Specifies the flavors information. The structure is described below.
  Changing this creates a new instance.

* `backup_strategy` - (Optional) Specifies the advanced backup policy. The structure is
  described below. Changing this creates a new instance.

* `ssl` - (Optional) Specifies whether to enable or disable SSL. Defaults to true.
-> The instance will be restarted in the background when switching SSL. Please operate with caution.

* `tags` - (Optional) Tags key/value pairs to associate with the volume.
  Changing this updates the existing volume tags.

The `datastore` block supports:

* `type` - (Required) Specifies the database type. DDS Community Edition is supported.
  The value is `DDS-Community`.

* `version` - (Required) Specifies the database version.
The values are `3.2`, `3.4`, `4.0`, `4.2`, `4.4`.

* `storage_engine` - (Optional) Specifies the storage engine. Currently, DDS supports the WiredTiger and RocksDB
   storage engine. The values are `wiredTiger`, `rocksDB`.
WiredTiger engine supports versions `3.2`, `3.4`, `4.0` while RocksDB supports versions `4.2`, `4.4`

The `flavor` block supports:

* `type` - (Required) Specifies the node type. Valid value:
  * For a cluster instance, the value can be `mongos`, `shard`, or `config`.
  * For a replica set instance, the value is `replica`.
  * For a single node instance, the value is `single`.

* `num` - (Required) Specifies the node quantity. Valid value:
  * `mongos`: The value ranges from `2` to `16`.
  * `shard`: The value ranges from `2` to `16`.
  * `config`: The value is `1`.
  * `replica`: The value is `1`.
  * `single`: The value is `1`.

* `storage` - (Optional) Specifies the disk type. Valid value: `ULTRAHIGH` which indicates the type SSD.

-> This parameter is optional for all nodes except `mongos`. This parameter is invalid for
  the `mongos` nodes.

* `size` - (Optional) Specifies the disk size. The value must be a multiple of `10`. The unit is GB.
  * For a `cluster` instance, the storage space of a shard node can be `10` to `1000` GB, and the config
  storage space is `20` GB. This parameter is invalid for `mongos` nodes. Therefore, you do not need
  to specify the storage space for `mongos` nodes.
  * For a `replica set` instance, the value ranges from `10` to `2000`.

-> This parameter is mandatory for all nodes except `mongos`. This parameter is invalid
  for the `mongos` nodes.

* `spec_code` - (Required) Specifies the resource specification code.

The `backup_strategy ` block supports:

* `start_time` - (Required) Specifies the backup time window. Automated backups will be triggered
	during the backup time window. The value cannot be empty. It must be a valid value in the
	`"hh:mm-HH:MM"` format. The current time is in the UTC format.
	* The `HH` value must be 1 greater than the `hh` value.
	* The values from `mm` and `MM` must be the same and must be set to any of the following `00`, `15`, `30`, or `45`.

* `keep_days` - (Optional) Specifies the number of days to retain the generated backup files. The
	value range is from `0` to `732`.
	* If this parameter is set to `0`, the automated backup policy is not set.
	* If this parameter is not transferred, the automated backup policy is enabled by default.
    Backup files are stored for seven days by default.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `datastore` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `security_group_id` - See Argument Reference above.
* `password` - See Argument Reference above.
* `disk_encryption_id` - See Argument Reference above.
* `ssl` - See Argument Reference above.
* `mode` - See Argument Reference above.
* `flavor` - See Argument Reference above.
* `backup_strategy` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `db_username` - Indicates the DB Administator name.
* `status` - Indicates the the DB instance status.
* `port` - Indicates the database port number. The port range is `2100` to `9500`.
* `nodes` - Indicates the instance nodes information. Structure is documented below.
* `pay_mode` - Indicates the billing mode. `0`: indicates the pay-per-use billing mode.

The `nodes` block contains:

  - `id` - Indicates the node ID.
  - `name` - Indicates the node name.
  - `role` - Indicates the node role.
  - `type` - Indicates the node type.
  - `private_ip` - Indicates the private IP address of a node. This parameter is valid only for
     mongos nodes, replica set instances.
  - `public_ip` - Indicates the EIP that has been bound on a node. This parameter is valid only for
     mongos nodes of cluster instances, primary nodes and secondary nodes of replica set instances.
  - `status` - Indicates the node status.

## Timeouts
This resource provides the following timeouts configuration options:
  - `create` - Default is 30 minute.
  - `delete` - Default is 30 minute.

## Import

DDSv3 Instance can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_dds_instance_v3.instance_1 c1851195-cdcb-4d23-96cb-032e6a3ee667
```
