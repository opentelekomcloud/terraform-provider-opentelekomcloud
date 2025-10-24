---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_slow_logs_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-slow-logs-v3"
description: |-
  Get TaurusDB MySQL slow logs from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_slow_logs_v3

Use this data source to get the list of TaurusDB MySQL slow logs.

## Example Usage
```hcl
variable "instance_id" {}
variable "node_id" {}
variable "start_date" {}
variable "end_date" {}

data "opentelekomcloud_taurusdb_mysql_slow_logs_v3" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_date  = var.start_date
  end_date    = var.end_date
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Specifies the ID of the TaurusDB MySQL instance.

* `node_id` - (Required) Specifies the ID of the TaurusDB MySQL instance node.

* `start_date` - (Required) Specifies the start date in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_date` - (Required) Specifies the end date in the **yyyy-mm-ddThh:mm:ssZ** format.

* `type` - (Optional) Specifies the SQL statement type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The data source region.

* `slow_log_list` - Indicates the list of the slow logs.
  The [slow_log_list](#slow_log_list) structure is documented below.

<a name="slow_log_list"></a>
The `slow_log_list` block supports:

* `node_id` - Indicates the ID of the TaurusDB MySQL instance node.

* `count` - Indicates the number of executions.

* `time` - Indicates the execution time.

* `lock_time` - Indicates the lock wait time.

* `rows_sent` - Indicates the number of sent rows.

* `rows_examined` - Indicates the number of scanned rows.

* `database` - Indicates the database that slow query logs belong to.

* `users` - Indicates the name of the account.

* `query_sample` - Indicates the execution syntax.

* `type` - Indicates the statement type.

* `start_time` - Indicates the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `client_ip` - Indicates the IP address of the client.
