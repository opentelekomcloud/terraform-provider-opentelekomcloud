---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_error_logs_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-error-logs-v3"
description: |-
  Get TaurusDB MySQL error logs from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_error_logs_v3

Use this data source to get the list of TaurusDB MySQL error logs.

## Example Usage
```hcl
variable "instance_id" {}
variable "node_id" {}
variable "start_time" {}
variable "end_time" {}

data "opentelekomcloud_taurusdb_mysql_error_logs_v3" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Specifies the ID of the TaurusDB MySQL instance.

* `node_id` - (Required) Specifies the ID of the TaurusDB MySQL instance node.

* `start_time` - (Required) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `level` - (Optional) Specifies the log level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The data source region.

* `error_log_list` - Indicates the list of the error logs.
  The [error_log_list](#error_log_list) structure is documented below.

<a name="error_log_list"></a>
The `error_log_list` block supports:

* `node_id` - Indicates the ID of the TaurusDB MySQL instance node.

* `time` - Indicates the execution time.

* `level` - Indicates the error log level.

* `content` - Indicates the error log content.
