---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_network_connections_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cc-central-network-connections-v3"
description: |-
  Get the list of Cloud Connect central network connections from OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect central network connections you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_network_connections/index.html)

# opentelekomcloud_cc_central_network_connections_v3

Use this data source to query the connections on a Cloud Connect (CC) central network within OpenTelekomCloud.

-> The central network APIs are account-level (domain-scoped). Make sure the credentials used by the provider
have permission to query Cloud Connect resources.

## Example Usage

```hcl
variable "central_network_id" {}

data "opentelekomcloud_cc_central_network_connections_v3" "test" {
  central_network_id = var.central_network_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source.
  If omitted, the provider-level region will be used.

* `central_network_id` - (Required, String) The ID of the central network to which the connections belong.

* `connection_id` - (Optional, String) The ID of the connection used to query.

* `name` - (Optional, String) The name of the connection used to query.

* `state` - (Optional, String) The status of the connection used to query.

* `bandwidth_type` - (Optional, String) The bandwidth type used to query.
  The value can be **BandwidthPackage** or **TestBandwidth**.

* `connection_type` - (Optional, String) The connection type used to query.
  The value can be **ER-ER**, **ER-GDGW** or **ER-ER_ROUTE_TABLE**.

* `global_connection_bandwidth_id` - (Optional, String) The global connection bandwidth ID used to query.

* `is_cross_region` - (Optional, String) Whether to query only cross-region connections.
  The value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of connections that match the filter parameters.
  The [connections](#cc_connections) structure is documented below.

<a name="cc_connections"></a>
The `connections` block supports:

* `id` - The connection ID.

* `name` - The connection name.

* `description` - The description of the connection.

* `domain_id` - The ID of the account that the connection belongs to.

* `enterprise_project_id` - The enterprise project ID of the connection.

* `central_network_id` - The ID of the central network that the connection belongs to.

* `central_network_plane_id` - The ID of the central network plane that the connection belongs to.

* `global_connection_bandwidth_id` - The ID of the global connection bandwidth bound to the connection.

* `bandwidth_type` - The bandwidth type.

* `bandwidth_size` - The bandwidth size in Mbit/s.

* `state` - The status of the connection.

* `is_frozen` - Whether the connection is frozen.

* `connection_type` - The connection type.

* `created_at` - The creation time of the connection.

* `updated_at` - The latest update time of the connection.

* `connection_point_pair` - The two connection endpoints.
  The [connection_point_pair](#cc_connections_point_pair) structure is documented below.

<a name="cc_connections_point_pair"></a>
The `connection_point_pair` block supports:

* `id` - The instance ID.

* `project_id` - The project ID.

* `region_id` - The region ID.

* `site_code` - The geographic site code.

* `instance_id` - The connection endpoint ID.

* `parent_instance_id` - The parent resource ID.

* `type` - The resource type. The value can be **ER**, **GDGW** or **ER_ROUTE_TABLE**.
