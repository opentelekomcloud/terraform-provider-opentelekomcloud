---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_quotas_v1"
sidebar_current: "docs-opentelekomcloud-datasource-ces-quotas-v1"
description: |-
  Get details about CES quotas within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES quotas v1 you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/quotas/index.html)

# opentelekomcloud_ces_quotas_v1

Get details about CES quotas within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_ces_quotas_v1" "quotas_1" {}
```

## Attributes Reference

The following attributes are exported:

* `quotas` - Specifies the quota list. The structure is described below.
    * `resources` - Specifies the resource quota list. The structure is described below.

The `resources` block supports:

* `type` - Specifies the quota type. `alarm` indicates the alarm rule.
* `used` - Specifies the used amount of the quota.
* `unit` - Specifies the quota unit.
* `quota` - Specifies the total amount of the quota.
