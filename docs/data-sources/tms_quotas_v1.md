---
subcategory: "Tag Management Service (TMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_tms_quotas_v1"
sidebar_current: "docs-opentelekomcloud-datasource-tms-quotas-v1"
description: |-
  List TMS Quotas within OpenTelekomCloud.
---

Up-to-date reference of API arguments for TMS quotas you can get at
[documentation portal](https://docs.otc.t-systems.com/tag-management-service/api-ref/api_description/quotas/index.html)

# opentelekomcloud_tms_quotas_v1

Use this data source to get the list of tag quotas within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_tms_quotas_v1" "test" {}
```

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the list of quotas. The structure is documented below.
    * `quota_key` - Indicates the quota key.
    * `quota_limit` - Indicates the quota value/limit.
    * `used` - Indicates the used quota.
    * `unit` - Indicates the unit.
