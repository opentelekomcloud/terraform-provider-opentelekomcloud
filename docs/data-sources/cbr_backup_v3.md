---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cbr_backup_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cbr-backup-v3"
description: |-
  Get details about CBR backup resources from OpenTelekomCloud
---

Up-to-date reference of API arguments for CBR backups you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-backup-recovery/api-ref/cbr_apis/backups/querying_all_backups.html#listbackups)


# opentelekomcloud_cbr_backup_v3

Use this data source to get details about backup resources from OpenTelekomCloud.

## Example Usage

```hcl
variable "backup_id" {}

data "opentelekomcloud_cbr_backup_v3" "cbr_backup" {
  id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Specifies the ID of backup.

* `chackpoint_id` - (Optional) Specifies the restore point ID.

* `status` - (Optional) Specifies the backup status.

* `resource_name` - (Optional) Specifies the backup resource name.

* `image_type` - (Optional) Specifies the backup type.

* `resource_type` - (Optional) Specifies the type of backup objects.

* `resource_id` - (Optional) Specifies the backup object ID.

* `name` - (Optional) Specifies the backup name

* `parent_id` - (Optional) Specifies the ID of parent backup.

* `resource_az` - (Optional) Specifies the AZ of backup.

* `vault_id` - (Optional) Specifies the ID of backup vault.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - The time the backup was created.

* `description` - Backup description.

* `expired_at` - The time the backup will be expired.

* `project_id` - The project ID of backup.

* `resource_size` - Backup size in GB.

* `updated_at` - Indicates the update time.

* `provider_id` - Backup provider ID which is used to distinguish backup objects.

* `auto_trigger` - Specifies whether the backup is automatically generated.

* `bootable` - Specifies whether the backup is a system disk backup.

* `incremental` - Specifies whether the backup is an incremental backup.

* `snapshot_id` - The snapshot ID of the disk backup.

* `support_lld` - Specifies whether to allow lazyloading for fast restoration.

* `supported_restore_mode` - Restoration mode of the backup.

* `encrypted` - Specifies whether the backup is encrypted.

* `system_disk` - Specifies whether a disk is a system disk.
