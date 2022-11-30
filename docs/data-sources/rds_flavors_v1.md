---
subcategory: "Relational Database Service (RDS)"
---

**DEPRECATED**
# opentelekomcloud_rds_flavors_v1

Use this data source to get the ID of an available OpenTelekomCloud RDS flavor.

## Example Usage

```hcl
data "opentelekomcloud_rds_flavors_v1" "flavor" {
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
  speccode          = "rds.pg.s1.medium"
}
```

## Argument Reference

* `datastore_name` - (Required) The datastore name of the rds.

* `datastore_version` - (Required) The datastore version of the rds.

* `speccode` - (Optional) The spec code of a rds flavor.

## Available value for attributes

datastore_name | datastore_version | speccode
---- | ---- | ----
PostgreSQL | 9.5.5 <br> 9.6.3 <br> 9.6.5| rds.pg.s1.xlarge rds.pg.m1.2xlarge rds.pg.c2.xlarge rds.pg.s1.medium rds.pg.c2.medium rds.pg.s1.large rds.pg.c2.large rds.pg.m1.large rds.pg.s1.2xlarge rds.pg.m1.xlarge
MySQL| 5.6.33 <br>5.6.30  <br>5.6.34 <br>5.6.35 <br>5.6.36 <br>5.7.17 <br>5.7.20| rds.mysql.s1.medium rds.mysql.s1.large rds.mysql.s1.xlarge rds.mysql.s1.2xlarge rds.mysql.m1.2xlarge rds.mysql.c2.medium rds.mysql.c2.large rds.mysql.c2.xlarge rds.mysql.m1.large rds.mysql.m1.xlarge
SQLServer| 2014 SP2 SE | rds.mssql.s1.xlarge rds.mssql.m1.2xlarge rds.mssql.c2.xlarge rds.mssql.s1.2xlarge rds.mssql.m1.xlarge

## Attributes Reference

`id` is set to the ID of the found rds flavor. In addition, the following attributes are exported:

* `datastore_name` - See Argument Reference above.

* `datastore_version` - See Argument Reference above.

* `speccode` - See Argument Reference above.

* `name` - The name of the rds flavor.

* `ram` - The name of the rds flavor.
