---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_flavors_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-flavors-v3"
description: |-
  Get available TaurusDB MySQL flavors from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_flavors_v3

Use this data source to get available OpenTelekomCloud TaurusDB MySQL flavors.

## Example Usage

```hcl
data "opentelekomcloud_taurusdb_mysql_flavors_v3" "flavors" {
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Optional) Specifies the database engine. Only **gaussdb-mysql** is supported now.

* `version` - (Optional) Specifies the database version. Only **8.0** is supported now.

* `availability_zone_mode` - (Optional) Specifies the availability zone mode. Currently supports **single** and **multi**. Defaults to **single**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which flavors are obtained.

* `flavors` - Indicates the flavors information.
  The [flavors](#flavors) structure is documented below.

<a name="flavors"></a>
The `flavors` block supports:

* `name` - Indicates the name of the TaurusDB MySQL flavor.

* `vcpus` - Indicates the CPU size.

* `memory` - Indicates the memory size in GB.

* `type` - Indicates the arch type of the flavor.

* `mode` - Indicates the database mode.

* `version` - Indicates the database version.

* `az_status` - Indicates the flavor status in each availability zone.
