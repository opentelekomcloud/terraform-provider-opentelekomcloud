---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_configurations_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-configurations-v3"
description: |-
  Get available TaurusDB MySQL configurations from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_configurations_v3

Use this data source to get the list of parameter templates.

## Example Usage

```hcl
data "opentelekomcloud_taurusdb_mysql_configurations_v3" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The data source region.

* `configurations` - Indicates the list of parameter templates.

  The [configurations](#configurations_struct) structure is documented below.

<a name="configurations_struct"></a>
The `configurations` block supports:

* `id` - Indicates the ID of the parameter template.

* `name` - Indicates the name of the parameter template.

* `user_defined` - Indicates whether the parameter template is a custom template.
    + **false**: default parameter template.
    + **true**: custom template.

* `description` - Indicates the description of parameter template.

* `datastore_name` - Indicates the engine name.

* `datastore_version` - Indicates the engine version.

* `created_at` - Indicates the creation time in the **yyyy-MM-ddTHH:mm:ssZ** format.

* `updated_at` - Indicates the update time in the **yyyy-MM-ddTHH:mm:ssZ** format.
