---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_proxy_flavors_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-proxy-flavors-v3"
description: |-
  Get available TaurusDB MySQL proxy flavors from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_proxy_flavors_v3

Use this data source to get the list of OpenTelekomCloud TaurusDB MySQL proxy flavors.

## Example Usage
```hcl
variable "instance_id" {}

data "opentelekomcloud_taurusdb_mysql_proxy_flavors_v3" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Specifies the ID of TaurusDB MySQL Instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which proxy flavors are obtained.

* `flavor_groups` - Indicates the list of flavor groups.
  The [flavor_groups](#flavor_groups) structure is documented below.

<a name="flavor_groups"></a>
The `flavor_groups` block supports:

* `type` - Indicates the group type. The value can be **arm** or **x86**.

* `flavors` - Indicates the list of flavors.
  The [flavors](#flavors) structure is documented below.

<a name="flavors"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the proxy flavor.

* `db_type` - Indicates the database type.

* `vcpus` - Indicates the number of vCPUs.

* `ram` - Indicates the memory size in GB.

* `spec_code` - Indicates the proxy specification code.

* `az_status` - Indicates the key/value pairs of the availability zone status.
  **key** indicates the AZ ID, and **value** indicates the specification status in the AZ.
