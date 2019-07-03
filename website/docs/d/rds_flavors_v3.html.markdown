---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_flavors_v3"
sidebar_current: "docs-opentelekomcloud-datasource-rds-flavors-v3"
description: |-
  Get the flavor information on an OpenTelekomCloud rds service.
---

# opentelekomcloud\_rds\_flavors\_v3

Use this data source to get available OpenTelekomCloud rds flavors.

## Example Usage

```hcl
data "opentelekomcloud_rds_flavors_v3" "flavor" {
    db_type = "PostgreSQL"
    db_version = "9.5"
    instance_mode = "ha"
}
```

## Argument Reference

* `db_type` - (Required) Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer.

* `db_version` -
  (Required)
  Specifies the database version. MySQL databases support MySQL 5.6
  and 5.7. PostgreSQL databases support
  PostgreSQL 9.5 and 9.6. Microsoft SQL Server
  databases support 2014_SE, 2016_SE, and 2016_EE.

* `instance_mode` - (Required) The mode of instance. Value: ha(indicates primary/standby instance), single(indicates single instance)

## Attributes Reference

In addition, the following attributes are exported:

* `flavors` -
  Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the rds flavor.
* `vcpus` - Indicates the CPU size.
* `memory` - Indicates the memory size in GB.
* `mode` - See 'instance_mode' above.
