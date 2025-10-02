---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_backup_v3"
sidebar_current: "docs-opentelekomcloud-resource-taurusdb-mysql-backup-v3"
description: |-
  Manages an TaurusDb mysql backup resource within OpenTelekomCloud.
---

# opentelekomcloud_taurusdb_mysql_backup_v3

Manages a TaurusDb MySQL backup resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_taurusdb_mysql_backup_v3" "test" {
  instance_id = var.instance_id
  name        = "test_backup_name"
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the TaurusDb MySQL instance.

* `name` - (Required, String, ForceNew) Specifies the name of the backup. It must start with a letter and consist of
  `4` to `64` characters. Only letters (case-sensitive), digits, hyphens (-), and underscores (_) are allowed.

* `description` - (Optional, String, ForceNew) Specifies the description of the backup.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `region` - The resource region

* `begin_time` - Indicates the backup start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the backup end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `take_up_time` - Indicates the backup duration in minutes.

* `size` - Indicates the backup size in MB.

* `datastore` - Indicates the database information.

  The [datastore](#backup_datastore_struct) structure is documented below.

<a name="backup_datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the database engine.

* `version` - Indicates the database version.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

## Import

The TaurusDb Mysql backup can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_taurusdb_mysql_backup_v3.test <id>
```
