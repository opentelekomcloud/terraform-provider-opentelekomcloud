---
subcategory: "Relational Database Service (RDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_instance_v3"
sidebar_current: "docs-opentelekomcloud-datasource-rds-instance-v3"
description: |-
Get available RDSv3 instance from OpenTelekomCloud
---

Up-to-date reference of API arguments for RDSv3 instance you can get at
[documentation portal](https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/db_instance_management/querying_details_about_db_instances.html#rds-01-0004)

# opentelekomcloud_rds_instance_v3

Use the `opentelekomcloud_rds_instance_v3` datasource to query DB instances according to search criteria.

## Example Usage

```hcl
data "opentelekomcloud_rds_instance_v3" "instance" {
  name          = "rds_instance_1"
  id            = "rds_instance_1_id"
  type          = "single"
  database_type = "PostgreSQL"
  vpc_id        = "vpc-id"
  subnet_id     = "subnet-id"
}
```

## Argument Reference

* `name` - (Optional) Specifies the DB instance ID.
* `id` - (Optional) ID of the RDS instance.
* `type` - (Optional) Specifies the instance type based query.
           The value is Single, Ha, or Replica, which correspond to single instance,
           primary/standby instances, and read replica, respectively.
* `database_type` - (Optional) Specifies the database type.
                    Its value can be any of the following and is case-sensitive:
                    `MySQL`, ` PostgreSQL`, `SQLServer`
* `vpc_id` - (Optional) Specifies the VPC ID.
* `subnet_id` - (Optional) Specifies the network ID of the subnet.


## Attributes Reference

The following attributes are exported:

* `id` - Indicates the DB instance ID.

* `name` - Indicates created the DB instance name.

* `status` - Indicates the DB instance status.

* `private_ips` - Indicates the private IP address. It is a blank string until an ECS is created.

* `public_ips` - Indicates the public IP address.

* `port` - Indicates the database port number.

* `type` - The value is Single, Ha, or Replica, which correspond to single instance,
           primary/standby instances, and read replica, respectively.

* `ha/replication_mode` - Indicates the replication mode for the standby DB instance.
                          The value cannot be empty.

* `region` - Indicates the region where the DB instance is deployed.

* `datastore_version` - Indicates the database version.

* `database_type` - Indicates the database type.

* `vpc_id` - Indicates the VPC ID.

* `subnet_id` - Indicates the network ID of the subnet.

* `security_group_id` - Indicates the security group ID.

* `security_group_name` - Indicates the security group name.

* `volume_size` - Indicates the volume size.

* `volume_type` - Indicates the volume type.

* `backup_strategy/start_time` - Indicates the backup time.

* `backup_strategy/keep_days` - Indicates the backup retention period.

* `availability_zone` - Indicates the availability zone.

* `created` - Indicates the creation time.

* `updated` - Indicates the update time.

* `flavor` - Indicates the flavor ID.

* `tags` - Indicates the tags.

* `timezone` - Indicates the time zone.

* `db_username` - Indicates the database username.

* `disk_encryption_id` - Indicates the disk encryption ID.

* `nodes` - Indicates the node information.

* `nodes/availability_zone` - Indicates the availability zone.

* `nodes/id` - Indicates the node ID.

* `nodes/name` - Indicates the node name.

* `nodes/role` - Indicates the node role.
                 The value can be master, slave, or readreplica, indicating the primary node,
                 standby node, and read replica node, respectively.

* `nodes/status` - Indicates the node status.
