---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_configuration_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-configuration-v3"
description: |-
  Get available TaurusDB MySQL configuration from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_configuration_v3

Use this data source to get available OpenTelekomCloud TaurusDB MySQL configuration.

## Example Usage

```hcl
data "opentelekomcloud_taurusdb_mysql_configuration_v3" "this" {
  name = "Default-TaurusDB V2.0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specifies the name of the parameter template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the configuration.

* `region` - Indicates the region of the configuration.

* `description` - Indicates the description of the configuration.

* `datastore_name` - Indicates the datastore name of the configuration.

* `datastore_version` - Indicates the datastore version of the configuration.
