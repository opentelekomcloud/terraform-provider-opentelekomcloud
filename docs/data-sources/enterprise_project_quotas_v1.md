---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_enterprise_project_quotas_v1"
sidebar_current: "docs-opentelekomcloud-datasource-enterprise-project-quotas-v1"
description: "Use this data source to get the list of quotas of EPS resources within T-Cloud Public (former OpenTelekomCloud)"
---

# opentelekomcloud_enterprise_project_quotas_v1

Use this data source to get the list of EPS resource quotas within T-Cloud Public (former OpenTelekomCloud).

## Example Usage

```hcl
data "opentelekomcloud_enterprise_project_quotas_v1" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The region used for the query.

* `region` - The region of the resource quotas.

* `quotas` - The list of the resource quotas.
  The [quotas](#eps_quotas) structure is documented below.

<a name="eps_quotas"></a>
The `quotas` block supports:

* `quota` - The total number of the resource quota.

* `type` - The resource type corresponding to quota.

* `used` - The used quota number.
