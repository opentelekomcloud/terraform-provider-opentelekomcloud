---
subcategory: "Cloud Server Backup Service (CSBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_csbs_backup_v1"
sidebar_current: "docs-opentelekomcloud-resource-csbs-backup-v1"
description: |-
  Manages a CSBS Backup resource within OpenTelekomCloud.
---


Up-to-date reference of API arguments for CSBS backup you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-server-backup-service/api-ref/api_description/backup_management)

# opentelekomcloud_csbs_backup_v1

Provides an OpenTelekomCloud Backup of Resources.

## Example Usage

```hcl
variable "backup_name" {}
variable "resource_id" {}

resource "opentelekomcloud_csbs_backup_v1" "backup_v1" {
  backup_name   = var.backup_name
  resource_id   = var.resource_id
  resource_type = "OS::Nova::Server"
}
```

## Argument Reference

The following arguments are supported:

* `backup_name` - (Optional) Name for the backup. The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-). Changing backup_name creates a new backup.

* `description` - (Optional) Backup description. The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<). Changing description creates a new backup.

* `resource_id` - (Required) ID of the target to which the backup is restored. Changing this creates a new backup.

* `resource_type` - (Optional) Type of the target to which the backup is restored. The default value is `OS::Nova::Server` for an ECS. Changing this creates a new backup.

* `tags` - (Optional) block supports the following arguments:

* `key` - (Required) Tag key. It cannot be an empty string.Changing key creates a new backup.

* `value` - (Required) Tag value. It can be an empty string.Changing value creates a new backup.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - It specifies the status of backup.

* `backup_record_id` - Specifies backup record ID.

The `volume_backups` block supports the following arguments:

* `status` -  Status of backup Volume.

* `space_saving_ratio` -  Specifies space saving rate.

* `name` -  It gives EVS disk backup name.

* `bootable` -  Specifies whether the disk is bootable.

* `average_speed` -  Specifies the average speed.

* `source_volume_size` -  Shows source volume size in GB.

* `source_volume_id` -  It specifies source volume ID.

* `incremental` -  Shows whether incremental backup is used.

* `snapshot_id` -  ID of snapshot.

* `source_volume_name` -  Specifies source volume name.

* `image_type` -  It specifies backup. The default value is backup.

* `id` -  Specifies Cinder backup ID.

* `size` -  Specifies accumulated size (MB) of backups.

The `vm_metadata` block supports the following arguments:

* `name` - Name of backup data.

* `eip` - Specifies elastic IP address of the ECS.

* `cloud_service_type` - Specifies ECS type.

* `ram` - Specifies memory size of the ECS, in MB.

* `vcpus` - Specifies CPU cores corresponding to the ECS.

* `private_ip` - It specifies internal IP address of the ECS.

* `disk` - Shows system disk size corresponding to the ECS specifications.

* `image_type` - Specifies image type.

## Import

Backup can be imported using  `backup_record_id`, e.g.

```sh
terraform import opentelekomcloud_csbs_backup_v1.backup_v1.backup_v1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```
