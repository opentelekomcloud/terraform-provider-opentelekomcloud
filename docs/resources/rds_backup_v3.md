---
subcategory: "Relational Database Service (RDS)"
---

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

* `backup_id` - The ID of the backup.
* `instance_id` - The ID of the RDS instance to which the backup belongs.
* `name` - The name of the backup.
* `type` - The type of the backup.
* `description` - The description of the backup.
* `created_at` - The creation time of the backup.
* `updated_at` - The update time of the backup.

## Import

Backups can be imported using the id, e.g.

```sh
terraform import opentelekomcloud_rds_backup_v3.my_backup 7117d38e-4c8f-1234-a707-bd96b97d024c
```
