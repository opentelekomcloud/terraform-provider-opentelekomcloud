---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_sql_control_rule_v3"
sidebar_current: "docs-opentelekomcloud-resource-taurusdb-mysql-sql-control-rule-v3"
description: |-
  Manages a TaurusDB MySQL SQL concurrency control rule resource within OpenTelekomCloud.
---

# opentelekomcloud_taurusdb_mysql_sql_control_rule_v3

Manages a TaurusDB MySQL SQL concurrency control rule resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "opentelekomcloud_taurusdb_mysql_sql_control_rule_v3" "test" {
  instance_id     = var.instance_id
  node_id         = var.node_id
  sql_type        = "SELECT"
  pattern         = "select~from~t1"
  max_concurrency = 20
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the TaurusDB MySQL instance.
  Changing this parameter will create a new resource.

* `node_id` - (Required, String, ForceNew) Specifies the ID of the TaurusDB MySQL node.
  Changing this parameter will create a new resource.

* `sql_type` - (Required, String, ForceNew) Specifies the SQL statement type.
  Value options: **SELECT**, **UPDATE**, **DELETE**.
  Changing this parameter will create a new resource.

* `pattern` - (Required, String, ForceNew) Specifies the concurrency control rule of SQL statements. A rule can consist
  of up to 128 keywords. The keywords are separated by tildes (~), for example, select~from~t1. The rule cannot contain
  backslashes (\), commas (,), or double tildes (~~). It cannot end with tildes (~).
  Changing this parameter will create a new resource.

* `max_concurrency` - (Required, Int) Specifies the maximum number of concurrent SQL statements.
  Value: a non-negative integer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID which is formatted `<instance_id>/<node_id>/<sql_type>/<pattern>`.

* `region` - Indicates the region in which to create the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The TaurusDB MySQL SQL concurrency control rule can be imported using the `instance_id`, `node_id`, `sql_type` and
`pattern` separated by slashes, e.g.

```bash
$ terraform import opentelekomcloud_taurusdb_mysql_sql_control_rule_v3.test <instance_id>/<node_id>/<sql_type>/<pattern>
```
