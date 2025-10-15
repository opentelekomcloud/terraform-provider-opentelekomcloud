---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_backups_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-backups-v3"
description: |-
  Get the list of TaurusDB MySQL backups from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_backups_v3

Use this data source to get the list of TaurusDB MySQL backups.

## Example Usage

```hcl
variable "instance_id" {}

data "opentelekomcloud_taurusdb_mysql_backups_v3" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, String) Specifies the ID of the TaurusDB MySQL instance.

* `backup_id` - (Optional, String) Specifies the ID of the backup.

* `backup_type` - (Optional, String) Specifies the backup type.
  Value options:
    + **auto**: automated full backup.
    + **manual**: manual full backup.

* `begin_time` - (Optional, String) Specifies the backup start time.
  The format is **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - (Optional, String) Specifies the backup end time.
  The format is **yyyy-mm-ddThh:mm:ssZ**.
  The end time must be later than the start time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - Indicates the list of backups.
  The [backups](#backups) structure is documented below.

<a name="backups"></a>
The `backups` block supports:

* `id` - Indicates the ID of the backup.

* `name` - Indicates the name of the backup.

* `instance_id` - Indicates the ID of the TaurusDB MySQL instance.

* `begin_time` - Indicates the backup start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the backup end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `take_up_time` - Indicates the backup duration in minutes.

* `type` - Indicates the backup type.

* `size` - Indicates the backup size in MB.

* `datastore` - Indicates the database information.
  The [datastore](#backups_datastore) structure is documented below.

* `status` - Indicates the backup status.

* `description` - Indicates the description of the backup.

<a name="backups_datastore"></a>
The `datastore` block supports:

* `type` - Indicates the database engine.

* `version` - Indicates the database version.
