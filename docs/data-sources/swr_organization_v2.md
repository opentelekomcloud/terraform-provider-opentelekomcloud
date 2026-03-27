---
subcategory: "Software Repository for Container (SWR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_swr_organization_v2"
sidebar_current: "docs-opentelekomcloud-data-source-swr-organization-v2"
description: |-
  Get details of SWR Organizations within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SWR organization you can get at
[documentation portal](https://docs.otc.t-systems.com/software-repository-container/api-ref/api)

# opentelekomcloud_swr_organization_v2

Get details of SWR organizations within Open Telekom Cloud.

## Example Usage

```hcl
data opentelekomcloud_swr_organization_v2 org_1 {}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the SWR organization. Use this to filter organizations list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `organizations` - List of organizations. The structure is documented below:
    * `name` - Organization name.
    * `organization_id` - Numeric ID of the organization.
    * `creator_name` - Username of the organization creator.
    * `auth` - User permission. The value can be `1`, `3`, or `7`. `7`: manage `3`: write `1`: read
