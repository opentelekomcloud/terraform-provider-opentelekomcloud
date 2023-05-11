---
subcategory: "Relational Database Service (RDS)"
---

Up-to-date reference of API arguments for RDS backup rule you can get at
`https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/backup_and_restoration`.

# opentelekomcloud_rds_backup_v3

Manages a manual RDS backup.

## Example Usage

```hcl
resource "opentelekomcloud_rds_instance_v3" "instance" {
  name       = "test-instance"
  engine     = "mysql"
  datastore  = "percona"
  flavor_ref = "rds.mysql.s1.large"
}

resource "opentelekomcloud_rds_backup_v3" "test" {
  instance_id = opentelekomcloud_rds_instance_v3.instance.id
  name        = "rds-backup-test-01"
  description = "manual"
}
```

## Argument Reference
The following arguments are supported:

* `instance_id` - (Required) The ID of the RDS instance to which the backup belongs.
* `name` - (Required) The name of the backup.
* `description` - (Optional) Specifies the backup description.
                  It contains a maximum of 256 characters and cannot contain the following special characters: >!<"&'=
* `databases` - (Optional) Specifies a list of self-built Microsoft SQL Server databases that are partially backed up.
                (Only Microsoft SQL Server support partial backups.)

## Attributes Reference
The following attributes are exported:

* `id` - The ID of the backup.
* `instance_id` - The ID of the RDS instance to which the backup belongs.
* `name` - The name of the backup.
* `description` - The description of the backup.
* `databases` - The list of self-built Microsoft SQL Server databases that are partially backed up.
                (Only Microsoft SQL Server support partial backups.)
* `begin_time` - Indicates the backup start time in the "yyyy-mm-ddThh:mm:ssZ" format,
                 where "T" indicates the start time of the time field, and "Z" indicates the time zone offset.
* `status` - Indicates the backup status. Value:
             - BUILDING: Backup in progress
             - COMPLETED: Backup completed
             - FAILED: Backup failed
             - DELETING: Backup being deleted
* `type` - Indicates the backup type. Value:
           - auto: automated full backup
           - manual: manual full backup
           - fragment: differential full backup
           - incremental: automated incremental backup
