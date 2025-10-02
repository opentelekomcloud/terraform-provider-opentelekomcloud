---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_instance_v3"
sidebar_current: "docs-opentelekomcloud-resource-taurusdb-mysql-instance-v3"
description: |-
  Manages an TaurusDb mysql instance resource within OpenTelekomCloud.
---

# opentelekomcloud_taurusdb_mysql_instance_v3

TaurusDB mysql instance management within OpenTelekomCloud.

## Example Usage

### create a basic instance

```hcl
resource "opentelekomcloud_taurusdb_mysql_instance_v3" "instance_1" {
  name              = "taurusdb_instance_1"
  password          = var.password
  flavor            = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
}
```

### create a taurusdb mysql instance with backup strategy

```hcl
resource "opentelekomcloud_taurusdb_mysql_instance_v3" "instance_1" {
  name              = "taurusdb_instance_1"
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

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name.
  The value must be `4` to `64` characters in length and start with a letter.
  It is case-sensitive and can contain only letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String) Specifies the instance specifications.

* `password` - (Required, String) Specifies the database password. The value must be `8` to `32` characters in length,
  including uppercase and lowercase letters, digits, and special characters, such as ~!@#%^*-_=+?

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet. Changing this parameter will create a
  new resource.

* `dedicated_resource_id` - (Optional, String, ForceNew) Specifies the dedicated resource ID. Changing this parameter
  will create a new resource.

* `security_group_id` - (Optional, String) Specifies the security group ID. Required if the selected subnet doesn't
  enable network ACL.

* `port` - (Optional, Int) Specifies the database port.

* `configuration_id` - (Optional, String) Specifies the configuration ID.

* `seconds_level_monitoring_enabled` - (Optional, Bool) Specifies whether to enable seconds level monitoring.

* `seconds_level_monitoring_period` - (Optional, Int) Specifies the seconds level collection period.
    + This parameter is valid only when `seconds_level_monitoring_enabled` is set to **true**.
    + This parameter can not be specified when `seconds_level_monitoring_enabled` is set to **false**.
    + Value options:
        - **1**: The collection period is 1s.
        - **5** (default value): The collection period is 5s.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id. Required if EPS enabled.

* `table_name_case_sensitivity` - (Optional, Bool) Whether the kernel table name is case sensitive. The value can
  be `true` (case sensitive) and `false` (case insensitive). Defaults to `false`. This parameter only works during
  creation.

* `read_replicas` - (Optional, Int) Specifies the count of read replicas. Defaults to `1`.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone. Defaults to "UTC+08:00". Changing this parameter
  will create a new resource.

* `availability_zone_mode` - (Optional, String, ForceNew) Specifies the availability zone mode: "single" or "multi".
  Defaults to "single". Changing this parameter will create a new resource.

* `master_availability_zone` - (Optional, String, ForceNew) Specifies the availability zone where the master node
  resides. The parameter is required in multi availability zone mode. Changing this parameter will create a new
  resource.

* `datastore` - (Optional, List, ForceNew) Specifies the database information. Structure is documented below. Changing
  this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below.

* `volume_size` - (Optional, Int) Specifies the volume size of the instance. The new storage space must be greater than
  the current storage and must be a multiple of `10` GB. Only valid when in prePaid mode.

The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the database engine. Only "gaussdb-mysql" is supported now.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the database version. Only "8.0" is supported now.
  Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the "hh:mm-HH:MM" format. The current time is in the UTC format. The
  HH value must be 1 greater than the hh value. The values of mm and MM must be the same and must be set to 00. Example
  value: **08:00-09:00**, **03:00-04:00**.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files.
  The value ranges from `0` to `35`. If this parameter is set to `0`, the automated backup policy is not set.
  If this parameter is not transferred, the automated backup policy is enabled by default.
  Backup files are stored for seven days by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.

* `region` - Indicates the DB instance region.

* `status` - Indicates the DB instance status.

* `mode` - Indicates the instance mode.

* `db_user_name` - Indicates the default username.

* `private_dns_name` - Indicates the private domain name.

* `private_write_ip` - Indicates the private IP address of the DB instance.

* `created_at` - Indicates the creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `updated_at` - Indicates the Update time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `nodes` - Indicates the instance nodes information.
  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block contains:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `type` - Indicates the node type: master or slave.

* `status` - Indicates the node status.

* `private_read_ip` - Indicates the private IP address of a node.

* `availability_zone` - Indicates the availability zone where the node resides.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 30 minutes.

## Import

TaurusDB instance can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_taurusdb_mysql_instance_v3.test <id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attribute is: `table_name_case_sensitivity`, `enterprise_project_id`, `password`.
It is generally recommended running `terraform plan` after importing
a TaurusDB MySQL instance. You can then decide if changes should be applied to the TaurusDB MySQL instance, or the resource
definition should be updated to align with the TaurusDB MySQL instance. Also you can ignore changes as below.

```hcl
resource "opentelekomcloud_taurusdb_mysql_instance_v3" "test" {
  lifecycle {
    ignore_changes = [
      table_name_case_sensitivity, enterprise_project_id, password
    ]
  }
}
```
