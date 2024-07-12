---
subcategory: "GaussDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gaussdb_mysql_instance_v3"
sidebar_current: "docs-opentelekomcloud-resource-gaussdb-mysql-instance-v3"
description: |-
  Manages a GaussDB for MySql resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for GaussDB for MySql you can get at
[documentation portal](https://docs.otc.t-systems.com/gaussdb-mysql/api-ref/apis_recommended/managing_db_instances/index.html#gaussdb-04-0003).

# opentelekomcloud_gaussdb_mysql_instance_v3

GaussDB MySql instance management within OpenTelekomCloud.

## Example Usage

### Create a basic instance

```hcl
resource "opentelekomcloud_gaussdb_mysql_instance_v3" "instance" {
  name                     = "gaussdb_instance"
  vpc_id                   = var.vpc_id
  subnet_id                = var.subnet_id
  security_group_id        = var.secgroup_id
  flavor                   = "gaussdb.mysql.xlarge.x86.8"
  password                 = var.password
  availability_zone_mode   = "multi"
  master_availability_zone = "eu-de-01"
  read_replicas            = 1
}
```

### Create an instance with backup strategy

```hcl
resource "opentelekomcloud_gaussdb_mysql_instance_v3" "instance" {
  name              = "gaussdb_instance_1"
  password          = var.password
  flavor            = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 7
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the GaussDB mysql instance resource. If omitted,
  the provider-level region will be used.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name. The value
  must be 4 to 64 characters in length and start with a letter. It is case-sensitive and can contain only letters,
  digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String) Specifies the instance specifications.

* `password` - (Required, String) Specifies the database password. The value must be 8 to 32 characters in length,
  including uppercase and lowercase letters, digits, and special characters, such as ~!@#%^*-_=+?

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet.

* `security_group_id` - (Optional, String, ForceNew) Specifies the security group ID. Required if the selected subnet
  doesn't enable network ACL.

* `configuration_id` - (Optional, String, ForceNew) Specifies the configuration ID.

* `configuration_name` - (Optional, String, ForceNew) Specifies the configuration name.

* `dedicated_resource_id` - (Optional, String, ForceNew) Specifies the dedicated resource ID. Changing this parameter
  will create a new resource.

* `dedicated_resource_name` - (Optional, String, ForceNew) Specifies the dedicated resource name. Changing this parameter
  will create a new resource.

* `read_replicas` - (Optional, Int) Specifies the count of read replicas. Defaults to 1.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone. Defaults to "UTC+08:00". Changing this parameter
  will create a new resource.

* `availability_zone_mode` - (Optional, String, ForceNew) Specifies the availability zone mode: "single" or "multi".
  Defaults to "single". Changing this parameter will create a new resource.

* `master_availability_zone` - (Optional, String, ForceNew) Specifies the availability zone where the master node
  resides. The parameter is required in multi availability zone mode. Changing this parameter will create a new
  resource.

* `datastore` - (Optional, List, ForceNew) Specifies the database information. Structure is documented below.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below.

The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the database engine. Only "gaussdb-mysql" is supported now.

* `version` - (Required, String, ForceNew) Specifies the database version. Only "8.0" is supported now.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the "hh:mm-HH:MM" format. The current time is in the UTC format. The
  HH value must be 1 greater than the hh value. The values of mm and MM must be the same and must be set to 00. Example
  value: 08:00-09:00, 03:00-04:00.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files. The value ranges from
  0 to 35. If this parameter is set to 0, the automated backup policy is not set. If this parameter is not transferred,
  the automated backup policy is enabled by default. Backup files are stored for seven days by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.
* `status` - Indicates the DB instance status.
* `port` - Indicates the database port.
* `mode` - Indicates the instance mode.
* `db_user_name` - Indicates the default username.
* `private_write_ip` - Indicates the private IP address of the DB instance.
* `nodes` - Indicates the instance nodes information. Structure is documented below.
* `charging_mode` - Indicates the charging mode of the instance.
* `project_id` - Indicates the id of the project.
* `alias` - Indicates the alias of the instance.
* `public_ip` - Indicates the public IP address of the DB instance.
* `node_count` - Indicates the amount on nodes of the DB instance.
* `created` - Indicates the created time of the DB instance.
* `updated` - Indicates the updated time of the DB instance.

*
The `nodes` block contains:

* `id` - Indicates the node ID.
* `name` - Indicates the node name.
* `type` - Indicates the node type: master or slave.
* `status` - Indicates the node status.
* `port` - Indicates the database port.
* `private_read_ip` - Indicates the private IP address of a node.
* `az_code` - Indicates the availability zone where the node resides.
* `region_code` - Indicates the region where the node resides.
* `created` - Indicates the created time of the DB node.
* `updated` - Indicates the updated time of the DB node.
* `flavor_ref` - Indicates the specification code of DB node.
* `max_connections` - Indicates the maximum number of connections of DB node.
* `vcpus` - Indicates the vCPUs number of DB node.
* `ram` - Indicates the memory size in GB of the DB node.
* `need_restart` - Indicates whether the reboot of DB instance is needed for the parameter modifications to take effect.
* `priority` - Indicates the failover priority of the DB node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 30 minutes.

## Import

GaussDB instance can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_gaussdb_mysql_instance_v3.instance_1 1a801c1e01e6458d8eed810912e29d0cin07
```

Due to the security reasons, `password` can not be imported. It can be ignored as shown below.

```hcl
resource "opentelekomcloud_gaussdb_mysql_instance_v3" "instance_1" {
  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
