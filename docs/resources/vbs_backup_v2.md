---
subcategory: "Volume Backup Service (VBS)"
---

# opentelekomcloud_vbs_backup_v2

Provides an VBS Backup resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "backup_name" {}
variable "volume_id" {}

resource "opentelekomcloud_vbs_backup_v2" "mybackup" {
  volume_id = var.volume_id
  name      = var.backup_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the vbs backup. Changing the parameter will create new resource.

* `volume_id` - (Required) The id of the disk to be backed up. Changing the parameter will create new resource.

* `snapshot_id` - (Optional) The snapshot id of the disk to be backed up. Changing the parameter will create new resource.

* `description` - (Optional) The description of the vbs backup. Changing the parameter will create new resource.

* `tags` - (Optional) List of tags to be configured for the backup resources. Changing the parameter will create new resource.

  * `key` - (Required) Specifies the tag key.

  * `value` - (Required) Specifies the tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the vbs backup.

* `container` - The container of the backup.

* `status` - The status of the VBS backup.

* `availability_zone` - The AZ where the backup resides.

* `size` - The size of the vbs backup.

* `service_metadata` - The metadata of the vbs backup.

## Timeouts

This resource provides the following timeouts configuration options:

- `create` - Default is 10 minutes.

- `delete` - Default is 3 minutes.

## Import

VBS Backup can be imported using the `backup id`, e.g.

```sh
terraform import opentelekomcloud_vbs_backup_v2.mybackup 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
