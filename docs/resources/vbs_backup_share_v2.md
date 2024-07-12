---
subcategory: "Volume Backup Service (VBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vbs_backup_share_v2"
sidebar_current: "docs-opentelekomcloud-resource-vbs-backup-share-v2"
description: |-
  Manages an VBS Backup Share resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VBS backup share you can get at
[documentation portal](https://docs.otc.t-systems.com/volume-backup-service/api-ref/api_description/vbs_backups)

# opentelekomcloud_vbs_backup_share_v2

Provides an VBS Backup Share resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "backup_id" {}

variable "to_project_ids" {}

resource "opentelekomcloud_vbs_backup_share_v2" "backupshare" {
  backup_id      = var.backup_id
  to_project_ids = var.to_project_ids
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Required) The ID of the backup to be shared. Changing the parameter will create new resource.

* `to_project_ids` - (Required) The IDs of projects with which the backup is shared. Changing the parameter will create new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `container` - The container of the backup.

* `backup_status` - The status of the VBS backup.

* `description` - The status of the VBS backup.

* `availability_zone` - The AZ where the backup resides.

* `size` - The size of the vbs backup.

* `backup_name` - The backup name.

* `snapshot_id` - The ID of the snapshot associated with the backup.

* `volume_id` - The ID of the tenant to which the backup belongs.

* `share_ids` - The backup share IDs.

* `service_metadata` - The metadata of the vbs backup.

## Import

VBS Backup Share can be imported using the `backup id`, e.g.

```sh
terraform import opentelekomcloud_vbs_backup_share_v2.backupshare 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
