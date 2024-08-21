---
subcategory: "Data Replication Service (DRS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_drs_task_v3"
sidebar_current: "docs-opentelekomcloud-resource-drs-task-v3"
description: |-
  Manages a DRS Task resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DRS task you can get at
[documentation portal](https://docs.otc.t-systems.com/data-replication-service/api-ref/api/public_api_management/index.html#drs-03-0101)

# opentelekomcloud_drs_task_v3

Manages DRS task resource within OpenTelekomCloud.

## Example Usage

### Create a DRS task to migrate data to the OpenTelekomCloud RDS database

```hcl
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}
resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {}

resource "opentelekomcloud_rds_instance_v3" "mysql_1" {}

resource "opentelekomcloud_rds_instance_v3" "mysql_2" {}

resource "opentelekomcloud_drs_task_v3" "test" {
  name           = "test-drs-task"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "down"
  net_type       = "eip"
  migration_type = "FULL_TRANS"
  description    = "TEST"
  force_destroy  = "true"

  source_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_networking_floatingip_v2.fip_1.address
    port        = "3306"
    user        = "root"
    password    = "MySql_120521"
    instance_id = opentelekomcloud_rds_instance_v3.mysql_1.id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  }

  destination_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_networking_floatingip_v2.fip_2.address
    port        = 3306
    user        = "root"
    password    = "MySql_120521"
    instance_id = opentelekomcloud_rds_instance_v3.mysql_2.id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the job name. The name consists of 4 to 50 characters, starting with
  a letter. Only letters, digits, underscores (\_) and hyphens (-) are allowed.

* `type` - (Required, String, ForceNew) Specifies the job type. Changing this parameter will create a new
  resource. The options are as follows:
    + **migration**: Online Migration.
    + **sync**: Data Synchronization.
    + **cloudDataGuard**: Disaster Recovery.

* `engine_type` - (Required, String, ForceNew) Specifies the migration engine type.
  Changing this parameter will create a new resource. The options are as follows:
    + **mysql**:  MySQL migration, MySQL synchronization use.
    + **mongodb**: Mongodb migration use.
    + **cloudDataGuard-mysql**: Disaster recovery use.
    + **gaussdbv5**: GaussDB (for openGauss) synchronization use.
    + **mysql-to-kafka**: Synchronization from MySQL to Kafka use.
    + **taurus-to-kafka**: Synchronization from GaussDB(for MySQL) to Kafka use.
    + **gaussdbv5ha-to-kafka**: Synchronization from GaussDB primary/standby to Kafka use.
    + **postgresql**: Synchronization from PostgreSQL to PostgreSQL use.

* `direction` - (Required, String, ForceNew) Specifies the direction of data flow.
  Changing this parameter will create a new resource. The options are as follows:
    + **up**: To the cloud. The destination database must be a database in the current cloud.
    + **down**: Out of the cloud. The source database must be a database in the current cloud.
    + **non-dbs**: self-built database.

* `source_db` - (Required, List, ForceNew) Specifies the source database configuration.
  The `db_info` object structure of the `source_db` is documented below.
  Changing this parameter will create a new resource.

* `destination_db` - (Required, List, ForceNew) Specifies the destination database configuration.
  The `db_info` object structure of the `destination_db` is documented below.
  Changing this parameter will create a new resource.

* `net_type` - (Optional, String, ForceNew) Specifies the network type.
  Changing this parameter will create a new resource. The options are as follows:
    + **eip**: suitable for migration from an on-premises or other cloud database to a destination cloud database.
      An EIP will be automatically bound to the replication instance and released after the replication task is complete.
    + **vpc**: suitable for migration from one cloud database to another.
    + **vpn**: suitable for migration from an on-premises self-built database to a destination cloud database,
      or from one cloud database to another in a different region.

The default value is `eip`.

* `migration_type` - (Optional, String, ForceNew) Specifies migration type.
  Changing this parameter will create a new resource. The options are as follows:
    + **FULL_TRANS**: Full migration. Suitable for scenarios where services can be interrupted. It migrates all database
      objects and data, in a non-system database, to a destination database at a time.
    + **INCR_TRANS**: Incremental migration. Suitable for migration from an on-premises self-built database to a
      destination cloud database, or from one cloud database to another in a different region.
    + **FULL_INCR_TRANS**:  Full+Incremental migration. This allows to migrate data with minimal downtime. After a full
      migration initializes the destination database, an incremental migration parses logs to ensure data consistency
      between the source and destination databases.

The default value is `FULL_INCR_TRANS`.

* `migrate_definer` - (Optional, Bool, ForceNew) Specifies whether to migrate the definers of all source database
  objects to the `user` of `destination_db`. The default value is `true`.
  Changing this parameter will create a new resource.

* `limit_speed` - (Optional, List, ForceNew) Specifies the migration speed by setting a time period.
  The default is no speed limit. The maximum length is 3. Structure is documented below.
  Changing this parameter will create a new resource.

* `multi_write` - (Optional, Bool, ForceNew) Specifies whether to enable multi write. It is mandatory when `type`
  is `cloudDataGuard`. When the disaster recovery type is dual-active disaster recovery, set `multi_write` to `true`,
  otherwise to `false`. The default value is `false`. Changing this parameter will create a new resource.

* `expired_days` - (Optional, Int, ForceNew) Specifies how many days after the task is abnormal, it will automatically
  end. The value ranges from 14 to 100. the default value is `14`. Changing this parameter will create a new resource.

* `start_time` - (Optional, String, ForceNew) Specifies the time to start the job. The time format
  is `yyyy-MM-dd HH:mm:ss`. Start immediately by default. Changing this parameter will create a new resource.

* `destination_db_readonly` - (Optional, Bool, ForceNew) Specifies the destination DB instance as read-only helps
  ensure the migration is successful. Once the migration is complete, the DB instance automatically changes to
  Read/Write. The default value is `true`. Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the job, which contain a
  maximum of 256 characters, and certain special characters (including !<>&'"\\) are not allowed.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the DRS job.
  Changing this parameter will create a new resource.

* `force_destroy` - (Optional, Bool) Specifies whether to forcibly destroy the job even if it is running.
  The default value is `false`.

The `db_info` block supports:

* `engine_type` - (Required, String, ForceNew) Specifies the engine type of database. Changing this parameter will
  create a new resource. The options are as follows: `mysql`, `mongodb`, `gaussdbv5`, `postgresql`.

* `ip` - (Required, String, ForceNew) Specifies the IP of database. Changing this parameter will create a new resource.

* `port` - (Required, Int, ForceNew) Specifies the port of database. Changing this parameter will create a new resource.

* `user` - (Required, String, ForceNew) Specifies the user name of database.
  Changing this parameter will create a new resource.

* `password` - (Required, String, ForceNew) Specifies the password of database.
  Changing this parameter will create a new resource.

* `instance_id` - (Optional, String, ForceNew) Specifies the instance id of database when it is a RDS database.
  Changing this parameter will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) Specifies subnet ID of database when it is a RDS database.
  It is mandatory when `direction` is `down`. Changing this parameter will create a new resource.

* `region` - (Optional, String, ForceNew) Specifies the region which the database belongs when it is a RDS database.
  Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the name of database.
  Changing this parameter will create a new resource.

* `ssl_enabled` - (Optional, Bool, ForceNew) Specifies whether to enable SSL connection.
  Changing this parameter will create a new resource.

* `ssl_cert_key` - (Optional, String, ForceNew) Specifies the SSL certificate content, encrypted with base64.
  It is mandatory when `ssl_enabled` is `true`. Changing this parameter will create a new resource.

* `ssl_cert_name` - (Optional, String, ForceNew) Specifies SSL certificate name.
  It is mandatory when `ssl_enabled` is `true`. Changing this parameter will create a new resource.

* `ssl_cert_check_sum` - (Optional, String, ForceNew) Specifies the checksum of SSL certificate content.
  It is mandatory when `ssl_enabled` is `true`. Changing this parameter will create a new resource.

* `ssl_cert_password` - (Optional, String, ForceNew) Specifies SSL certificate password. It is mandatory when
  `ssl_enabled` is `true` and the certificate file suffix is `.p12`. Changing this parameter will create a new resource.

The `limit_speed` block supports:

* `speed` - (Required, String, ForceNew) Specifies the transmission speed, the value range is 1 to 9999, unit: `MB/s`.
  Changing this parameter will create a new resource.

* `start_time` - (Required, String, ForceNew) Specifies the time to start speed limit, this time is UTC time. The start
  time is the whole hour, if there is a minute, it will be ignored, the format is `hh:mm`, and the hour number
  is two digits, for example: 01:00. Changing this parameter will create a new resource.

* `end_time` - (Required, String, ForceNew) Specifies the time to end speed limit, this time is UTC time. The input must
  end at 59 minutes, the format is `hh:mm`, for example: 15:59. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The resource ID in UUID format.

* `created_at` - Create time. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ

* `status` - Status.

* `public_ip` - Public IP.

* `private_ip` - Private IP.

* `region` - The region in which to create the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.

* `delete` - Default is 10 minute.

## Import

The DRS job can be imported by `id`. For example,

```
terraform import opentelekomcloud_drs_task_v3.test b11b407c-e604-4e8d-8bc4-92398320b847
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `tags`, `force_destroy`,
`source_db.0.password` and `destination_db.0.password`.It is generally recommended running
`terraform plan` after importing a job. You can then decide if changes should be applied to the job, or the resource
definition should be updated to align with the job. Also you can ignore changes as below.

```
resource "opentelekomcloud_drs_job" "test" {
    ...

  lifecycle {
    ignore_changes = [
      source_db.0.password,destination_db.0.password
    ]
  }
}
```
