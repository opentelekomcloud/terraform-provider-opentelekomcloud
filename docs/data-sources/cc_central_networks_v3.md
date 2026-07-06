---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_networks_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cc-central-networks-v3"
description: |-
  Get the list of Cloud Connect central networks from OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect central networks you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_networks/index.html)

# opentelekomcloud_cc_central_networks_v3

Use this data source to query the Cloud Connect (CC) central networks within OpenTelekomCloud.

-> The central network APIs are account-level (domain-scoped). Make sure the credentials used by the provider
have permission to query Cloud Connect resources.

## Example Usage

```hcl
variable "central_network_name" {}

data "opentelekomcloud_cc_central_networks_v3" "test" {
  name = var.central_network_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source.
  If omitted, the provider-level region will be used.

* `central_network_id` - (Optional, String) The ID of the central network used to query.

* `name` - (Optional, String) The name of the central network used to query.

* `state` - (Optional, String) The status of the central network used to query.
  The value can be **AVAILABLE**, **CREATING**, **UPDATING**, **FAILED**, **DELETING**, **DELETED** or **RESTORING**.

* `enterprise_project_id` - (Optional, String) The enterprise project ID used to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `central_networks` - The list of central networks that match the filter parameters.
  The [central_networks](#cc_central_networks) structure is documented below.

<a name="cc_central_networks"></a>
The `central_networks` block supports:

* `id` - The central network ID.

* `name` - The central network name.

* `description` - The description of the central network.

* `state` - The status of the central network.

* `enterprise_project_id` - The ID of the enterprise project that the central network belongs to.

* `default_plane_id` - The ID of the default central network plane.

* `domain_id` - The ID of the account that the central network belongs to.

* `created_at` - The creation time of the central network.

* `updated_at` - The latest update time of the central network.
