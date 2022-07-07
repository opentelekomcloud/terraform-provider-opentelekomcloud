---
subcategory: "Cloud Backup and Recovery (CBR)"
---

# opentelekomcloud_cbr_backup_ids_v3

Use this data source to get details about backup resources from OpenTelekomCloud.

## Example Usage

```hcl
variable "checkpoint_id" {}

data "opentelekomcloud_cbr_backup_ids_v3" "cbr_backups" {
  checkpoint_id = var.checkpoint_id
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

The following attributes are exported:

* `ids` - A list of all the backup ids found. This data source will fail if none are found.
