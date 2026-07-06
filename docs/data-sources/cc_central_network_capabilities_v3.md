---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_network_capabilities_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cc-central-network-capabilities-v3"
description: |-
  Get the list of Cloud Connect central network capabilities from OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect central network capabilities you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_network_capabilities/index.html)

# opentelekomcloud_cc_central_network_capabilities_v3

Use this data source to query the capabilities supported by the Cloud Connect (CC) central network within
OpenTelekomCloud. Capabilities describe, for example, the supported regions, sites, billing options and
bandwidth ranges of the account.

-> The central network APIs are account-level (domain-scoped). Make sure the credentials used by the provider
have permission to query Cloud Connect resources.

## Example Usage

```hcl
data "opentelekomcloud_cc_central_network_capabilities_v3" "test" {
  capability = "central-network.connect-er-instances"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source.
  If omitted, the provider-level region will be used.

* `capability` - (Optional, String) The capability used to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `capabilities` - The list of capabilities that match the filter parameters.
  The [capabilities](#cc_capabilities) structure is documented below.

<a name="cc_capabilities"></a>
The `capabilities` block supports:

* `id` - The instance ID.

* `domain_id` - The ID of the account that the capability belongs to.

* `capability` - The capability name.

* `specifications` - The specifications of the capability.
  The [specifications](#cc_capabilities_specifications) structure is documented below.

<a name="cc_capabilities_specifications"></a>
The `specifications` block supports:

* `is_support` - Whether the capability is supported.

* `charge_mode` - The supported billing options.

* `support_regions` - The supported regions.

* `support_ipv6_regions` - The regions that support IPv6.

* `support_dscp_regions` - The regions that support DSCP.

* `support_sts5_regions` - The regions that support STS5.

* `support_freeze_regions` - The regions that support freezing.

* `size_range` - The supported bandwidth size range.
  The [size_range](#cc_capabilities_size_range) structure is documented below.

* `free_lines` - The free lines.
  The [free_lines](#cc_capabilities_free_lines) structure is documented below.

* `support_sites` - The supported sites.
  The [support_sites](#cc_capabilities_support_sites) structure is documented below.

<a name="cc_capabilities_size_range"></a>
The `size_range` block supports:

* `min` - The minimum bandwidth in Mbit/s.

* `max` - The maximum bandwidth in Mbit/s.

<a name="cc_capabilities_free_lines"></a>
The `free_lines` block supports:

* `local_site_code` - The local site code.

* `remote_site_code` - The remote site code.

<a name="cc_capabilities_support_sites"></a>
The `support_sites` block supports:

* `region_id` - The region ID.

* `site_code` - The site code.
