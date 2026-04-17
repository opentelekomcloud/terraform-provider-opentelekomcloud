---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_engine_versions_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-engine-versions_v3"
description: |-
  Get database specifications of a specified DB engine from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_engine_versions_v3

Use this data source to get the database specifications of a specified DB engine.

## Example Usage

```hcl
data "opentelekomcloud_taurusdb_mysql_engine_versions_v3" "test" {
  database_name = "gaussdb-mysql"
}
```

## Argument Reference

The following arguments are supported:

* `database_name` - (Required) Specifies the DB engine.
  Value options: **gaussdb-mysql**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastores` - Indicates the DB version list.
  The [datastores](#datastores) structure is documented below.

<a name="datastores"></a>
The `datastores` block supports:

* `id` - Indicates the DB version ID.

* `name` - Indicates the DB version number.
  Only the major version number with two digits is returned.
