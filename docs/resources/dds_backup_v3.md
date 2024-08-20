---
subcategory: "Document Database Service (DDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dds_backup_v3"
sidebar_current: "docs-opentelekomcloud-resource-dds-backup-v3"
description: |-
  Manages a DDS backup resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DDS instance you can get at
[documentation portal](https://docs.otc.t-systems.com/document-database-service/api-ref/apis_v3.0_recommended/backup_and_restoration/index.html)

# opentelekomcloud_dds_backup_v3

Manages a DDS backup resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}

resource "opentelekomcloud_dds_backup_v3" "backup" {
  instance_id = var.instance_id
  name        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDS instance.

* `name` - (Required, String, ForceNew) Specifies the manual backup name.
  The value must be `4` to `64` characters in length and start with a letter (from A to Z or from a to z).
  It is case-sensitive and can contain only letters, digits (from 0 to 9), hyphens (-), and underscores (_).

* `description` - (Optional, String, ForceNew) Specifies the manual backup description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `instance_name` - Indicates the name of a DDS instance.

* `datastore` - Indicates the database version.
  The [datastore](#datastore_struct) structure is documented below.

* `region` - Indicates the region in which resource was created.

* `type` - Indicates the backup type. Valid value:
  * `Manual`: indicates manual full backup.

* `begin_time` - Indicates the start time of the backup. The format is yyyy-mm-dd hh:mm:ss. The value is in UTC format.

* `end_time` - Indicates the end time of the backup. The format is yyyy-mm-dd hh:mm:ss. The value is in UTC format.

* `status` - Indicates the backup status. Valid value:
  + `BUILDING`: Backup in progress
  + `COMPLETED`: Backup completed
  + `FAILED`: Backup failed
  + `DISABLED`: Backup being deleted

* `size` - Indicates the backup size in KB.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the DB engine.

* `version` - Indicates the database version.

* `storage_engine` - Indicates the database storage engine.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The DDS backup can be imported using the `instance_id` and the `id` separated by a slash, e.g.:

```bash
$ terraform import opentelekomcloud_dds_backup_v3.backup <instance_id>/<id>
```
