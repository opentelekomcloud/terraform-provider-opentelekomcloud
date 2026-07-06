---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_network_quotas_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cc-central-network-quotas-v3"
description: |-
  Get the list of Cloud Connect central network quotas from OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect central network quotas you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_network_quotas/index.html)

# opentelekomcloud_cc_central_network_quotas_v3

Use this data source to query the Cloud Connect (CC) central network quotas within OpenTelekomCloud.

-> The central network APIs are account-level (domain-scoped). Make sure the credentials used by the provider
have permission to query Cloud Connect resources.

## Example Usage

```hcl
data "opentelekomcloud_cc_central_network_quotas_v3" "test" {
  quota_type = ["central_network_count"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source.
  If omitted, the provider-level region will be used.

* `quota_type` - (Optional, List) The quota types used to query. Multiple quota types can be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of quotas that match the filter parameters.
  The [quotas](#cc_quotas) structure is documented below.

<a name="cc_quotas"></a>
The `quotas` block supports:

* `quota_key` - The quota identifier.

* `quota_limit` - The quota limit.

* `used` - The number of used quotas.

* `unit` - The unit of the quota value.
