---
subcategory: "Relational Database Service (RDS)"
---

# opentelekomcloud_rds_versions_v3

Use this data source to get available OpenTelekomCloud rds versions.

## Example Usage

```hcl
data "opentelekomcloud_rds_versions_v3" "versions" {
  database_name = "mysql"
}
```

## Argument Reference

* `database_name` - (Required) Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer. Case-insensitive.

## Attributes Reference

In addition, the following attributes are exported:

* `versions` - List of version names, sorted by a version (higher to lower). Example: `["11", "10", "9.6", "9.5"]`.
