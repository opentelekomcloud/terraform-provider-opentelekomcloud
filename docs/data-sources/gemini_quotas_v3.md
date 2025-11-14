---
subcategory: "GeminiDB"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gemini_quotas_v3"
sidebar_current: "docs-opentelekomcloud-datasource-gemini-quotas-v3"
description: |-
  Get GeminiDB resource quotas from OpenTelekomCloud
---

# opentelekomcloud_gemini_quotas_v3

Use this data source to get the list of GeminiDB resource quotas for your tenant.

## Example Usage
```hcl
data "opentelekomcloud_gemini_quotas_v3" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The data source region.

* `quotas` - Indicates the list of resource quotas.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `type` - Indicates the quota resource type.

* `quota` - Indicates the current quota. If set to 0, no quantity limit is set for resources.

* `used` - Indicates the number of used resources.
