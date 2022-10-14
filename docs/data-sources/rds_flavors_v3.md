---
subcategory: "Relational Database Service (RDS)"
---

# opentelekomcloud_rds_flavors_v3

Use this data source to get available OpenTelekomCloud RDSv3 flavors.

## Example Usage

```hcl
data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "9.5"
  instance_mode = "ha"
}
```

## Argument Reference

* `db_type` - (Required) Specifies the DB engine. Possible values are: `MySQL`, `PostgreSQL`, `SQLServer`.

* `db_version` - (Required) Specifies the database version. `MySQL` databases support `5.6`,
  `5.7` and `8.0`. `PostgreSQL` databases support `9.5`, `9.6`, `10`, `11`, `12`  and `13`.
  `SQLServer` databases support `2014_SE`, `2016_SE`, `2016_EE`, `2017_SE` and `2017_EE`.

* `instance_mode` - (Required) The mode of instance. Possible values are:
  `ha` indicates primary/standby instance, `single` indicates single instance
  and `replica` indicates read-replica instance.

## Attributes Reference

In addition, the following attributes are exported:

* `flavors` - Indicates the `flavors` information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the rds flavor.

* `vcpus` - Indicates the CPU size.

* `memory` - Indicates the memory size in GB.

* `mode` - Indicates the DB instance type.

* `az_status` - Indicates the status of the AZ to which the DB instance specifications belong.
