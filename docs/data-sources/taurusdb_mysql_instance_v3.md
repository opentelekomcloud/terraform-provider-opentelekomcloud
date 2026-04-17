---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_instance_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-instance-v3"
description: |-
  Get available TaurusDB MySQL instance from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_instance_v3

Use this data source to get available OpenTelekomCloud TaurusDB MySQL instance.

## Example Usage

```hcl
data "opentelekomcloud_taurusdb_mysql_instance_v3" "this" {
  name = "taurusdb-instance"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specifies the name of the instance.

* `vpc_id` - (Optional) Specifies the VPC ID.

* `subnet_id` - (Optional) Specifies the network ID of a subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the instance.

* `region` - Indicates the instance region.

* `flavor` - Indicates the instance specifications.

* `security_group_id` - Indicates the security group ID.

* `configuration_id` - Indicates the configuration ID.

* `read_replicas` - Indicates the count of read replicas.

* `time_zone` - Indicates the time zone.

* `availability_zone_mode` - Indicates the availability zone mode: **single** or **multi**.

* `master_availability_zone` - Indicates the availability zone where the master node resides.

* `datastore` - Indicates the database information.
  The [datastore](#datastore) structure is documented below.

* `backup_strategy` - Indicates the advanced backup policy.
  The [backup_strategy](#backup_strategy) structure is documented below.

* `status` - Indicates the DB instance status.

* `port` - Indicates the database port.

* `mode` - Indicates the instance mode.

* `db_user_name` - Indicates the default username.

* `private_write_ip` - Indicates the private IP address of the DB instance.

* `nodes` - Indicates the instance nodes information.
  The [nodes](#nodes) structure is documented below.

<a name="datastore"></a>
The `datastore` block supports:

* `engine` - Indicates the database engine.

* `version` - Indicates the database version.

<a name="backup_strategy"></a>
The `backup_strategy` block supports:

* `start_time` - Indicates the backup time window.

* `keep_days` - Indicates the number of days to retain the generated backup.

<a name="nodes"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `type` - Indicates the node type: master or slave.

* `status` - Indicates the node status.

* `private_read_ip` - Indicates the private IP address of a node.

* `availability_zone` - Indicates the availability zone where the node resides.
