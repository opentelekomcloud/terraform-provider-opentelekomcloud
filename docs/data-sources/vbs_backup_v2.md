---
subcategory: "Volume Backup Service (VBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vbs_backup_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vbs-backup-v2"
description: |-
Get details about a specific VBS backup from OpenTelekomCloud
---

Up-to-date reference of API arguments for VBS backup you can get at
[documentation portal](https://docs.otc.t-systems.com/volume-backup-service/api-ref/api_description/vbs_backups/querying_details_about_vbs_backups_native_openstack_api.html#en-us-topic-0020237259)

# opentelekomcloud_vbs_backup_v2

Use this data source to get details about a specific VBS Backup.

## Example Usage

```hcl
variable "backup_id" {}

data "opentelekomcloud_vbs_backup_v2" "mybackup" {
  id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the vbs backup.

* `name` - (Optional) The name of the vbs backup.

* `volume_id` - (Optional) The source volume ID of the backup.

* `snapshot_id` - (Optional) ID of the snapshot associated with the backup.

* `status` - (Optional) The status of the VBS backup.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - The description of the vbs backup.

* `availability_zone` - The AZ where the backup resides.

* `size` - The size of the vbs backup.

* `container` - The container of the backup.

* `service_metadata` - The metadata of the vbs backup.

* `to_project_ids` - IDs of projects with which the backup is shared.

* `share_ids` - The backup share IDs.
