---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_project_quotas_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-project-quotas-v3"
description: |-
  Get TaurusDB MySQL project quotas from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_project_quotas_v3

Use this data source to get the project quotas of a specified tenant.

## Example Usage
```hcl
data "opentelekomcloud_taurusdb_mysql_project_quotas_v3" "test" {}
```

## Argument Reference

The following arguments are supported:


* `type` - (Optional) Specifies the resource type used to filter quotas. Value options: **instance**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The data source region.

* `quotas` - Indicates the tenant instance quota information.
  The [quotas](#quotas) structure is documented below.

<a name="quotas"></a>
The `quotas` block supports:

* `resources` - Indicates the resource list objects.
  The [resources](#quotas_resources) structure is documented below.

<a name="quotas_resources"></a>
The `resources` block supports:

* `type` - Indicates the quota of the specified type.

* `used` - Indicates the number of created resources.

* `quota` - Indicates the maximum resource quota.
