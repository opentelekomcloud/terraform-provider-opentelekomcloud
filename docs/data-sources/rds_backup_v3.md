---
subcategory: "Relational Database Service (RDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_backup_v3"
sidebar_current: "docs-opentelekomcloud-datasource-rds-backup-v3"
description: |-
  Get details about RDSv3 instance backup from OpenTelekomCloud
---

Up-to-date reference of API arguments for RDSv3 instance backup you can get at
[documentation portal](https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/backup_and_restoration/obtaining_details_about_backups.html#rds-09-0005)

# opentelekomcloud_rds_backup_v3

Use this data source to get information about RDSv3 instance backup.

## Example Usage

Finding the latest automatic backup:

```hcl
data "opentelekomcloud_rds_backup_v3" "backup" {
  instance_id = var.rds_instance_id
  type        = "auto"
}
```

## Argument Reference

* `instance_id` - (Required) Specifies the DB instance ID.

* `backup_id` - (Optional) Specifies the backup ID.

* `type` - (Optional) Specifies the backup type.

  Possible values:
    * `auto`: automated full backup.
    * `manual`: manual full backup.
    * `fragment`: differential full backup.
    * `incremental`: automated incremental backup.

## Attributes Reference

In addition, the following attributes are exported:

* `name` - Indicates the backup name.

* `status` - Indicates the status of the backup.

  Possible values:
    * `BUILDING`: Backup in progress
    * `COMPLETED`: Backup completed
    * `FAILED`: Backup failed
    * `DELETING`: Backup being deleted

* `type` - Indicates the backup type.

  Possible values:
    * `auto`: automated full backup.
    * `manual`: manual full backup.
    * `fragment`: differential full backup.
    * `incremental`: automated incremental backup.

* `size` - Indicates the backup size in kB.

* `begin_time` - Indicates the backup start time in the `yyyy-mm-ddThh:mm:ssZ` format.

* `end_time` - Indicates the backup end time in the `yyyy-mm-ddThh:mm:ssZ` format.

  -> In a full backup, `end_time` indicates the full backup end time. In a MySQL incremental backup, it indicates the
  time when the last transaction in the backup file is submitted.

* `databases` - Indicates a list of self-built Microsoft SQL Server databases that support partial backups.

* `db_type` - Indicates the DB engine.

* `db_version` - Indicates the database version.
