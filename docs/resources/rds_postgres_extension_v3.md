---
subcategory: "Relational Database Service (RDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_postgres_extension_v3"
sidebar_current: "docs-opentelekomcloud-resource-rds-postgres-extension-v3"
description: |-
  Manages an RDS Postgres extension v3 resource within T CLoud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for RDS Postgres extension you can get at
[documentation portal](https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/extension_management_rds_for_postgresql/index.html)

# opentelekomcloud_rds_postgres_extension_v3

Manages a PostgreSQL extension of an RDS instance.

## Example Usage

```hcl
variable "instance_id" {}
variable "database_name" {}

resource "opentelekomcloud_rds_postgres_extension_v3" "extension" {
  instance_id    = var.instance_id
  database_name  = var.database_name
  extension_name = "hstore"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RDS PostgreSQL instance.

* `database_name` - (Required, String, ForceNew) Specifies the name of the database in which the extension is created.

* `extension_name` - (Required, String, ForceNew) Specifies the name of the extension.

## Attribute Reference

The following attributes are exported:

* `id` - Specifies the resource ID, which is `{instance_id}/{database_name}/{extension_name}`.

* `version` - Specifies the current version of the extension.

* `version_update` - Specifies the new version that the extension can be upgraded to. If the value of this parameter is different from that of `version`, the extension can be upgraded.

* `description` - Specifies the description of the extension.

* `shared_preload_libraries` - Specifies the dependent preloaded library..

* `enable_install` - Specifies whether the extension can be installed..

## Import

RDS Postgres extensions can be imported using `instance_id/database_name/extension_name`, e.g.

```
$ terraform import opentelekomcloud_rds_postgres_extension_v3.extension 58760797-f992-44e9-a128-0ff3129989b5/my_db/hstore
```
