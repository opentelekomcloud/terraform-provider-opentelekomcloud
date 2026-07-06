---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_central_network_policies_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cc-central-network-policies-v3"
description: |-
  Get the list of Cloud Connect central network policies from OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect central network policies you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/central_network_policies/index.html)

# opentelekomcloud_cc_central_network_policies_v3

Use this data source to query the policies of a Cloud Connect (CC) central network within OpenTelekomCloud.

-> The central network APIs are account-level (domain-scoped). Make sure the credentials used by the provider
have permission to query Cloud Connect resources.

## Example Usage

```hcl
variable "central_network_id" {}

data "opentelekomcloud_cc_central_network_policies_v3" "test" {
  central_network_id = var.central_network_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source.
  If omitted, the provider-level region will be used.

* `central_network_id` - (Required, String) The ID of the central network to which the policies belong.

* `policy_id` - (Optional, String) The ID of the policy used to query.

* `state` - (Optional, String) The status of the policy used to query.
  The value can be **AVAILABLE**, **CANCELING**, **APPLYING**, **FAILED** or **DELETED**.

* `is_applied` - (Optional, String) Whether the policy is applied. The value can be **true** or **false**.

* `version` - (Optional, List) The policy versions used to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - The list of policies that match the filter parameters.
  The [policies](#cc_policies) structure is documented below.

<a name="cc_policies"></a>
The `policies` block supports:

* `id` - The policy ID.

* `central_network_id` - The ID of the central network that the policy belongs to.

* `domain_id` - The ID of the account that the policy belongs to.

* `state` - The status of the policy.

* `document_template_version` - The policy document template version.

* `is_applied` - Whether the policy is applied.

* `version` - The policy version.

* `created_at` - The creation time of the policy.

* `document` - The policy document.
  The [document](#cc_policies_document) structure is documented below.

<a name="cc_policies_document"></a>
The `document` block supports:

* `default_plane` - The name of the default central network plane.

* `er_instances` - The list of enterprise routers on the central network.
  The [er_instances](#cc_policies_er_instances) structure is documented below.

* `planes` - The list of central network planes.
  The [planes](#cc_policies_planes) structure is documented below.

<a name="cc_policies_er_instances"></a>
The `er_instances` block supports:

* `project_id` - The project ID.

* `region_id` - The region ID.

* `enterprise_router_id` - The enterprise router ID.

<a name="cc_policies_planes"></a>
The `planes` block supports:

* `name` - The name of the central network plane.

* `associate_er_tables` - The enterprise router route tables associated with the central network plane.
  The [associate_er_tables](#cc_policies_associate_er_tables) structure is documented below.

* `exclude_er_connections` - The connections between enterprise routers excluded from the central network plane.
  The [exclude_er_connections](#cc_policies_exclude_er_connections) structure is documented below.

<a name="cc_policies_associate_er_tables"></a>
The `associate_er_tables` block supports:

* `project_id` - The project ID.

* `region_id` - The region ID.

* `enterprise_router_id` - The enterprise router ID.

* `enterprise_router_table_id` - The enterprise router route table ID.

<a name="cc_policies_exclude_er_connections"></a>
The `exclude_er_connections` block supports:

* `exclude_er_instances` - The excluded enterprise router instances.
  The [exclude_er_instances](#cc_policies_er_instances) structure is the same as the `er_instances` block above.
