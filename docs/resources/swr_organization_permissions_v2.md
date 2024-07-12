---
subcategory: "Software Repository for Container (SWR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_swr_organization_permissions_v2"
sidebar_current: "docs-opentelekomcloud-resource-swr-organization-permissions-v2"
description: |-
  Manages an SWR Organization Permissions resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SWR permission you can get at
[documentation portal](https://docs.otc.t-systems.com/software-repository-container/api-ref/api)

# opentelekomcloud_swr_organization_permissions_v2

Manages user permissions for the SWR organization resource within Open Telekom Cloud.

## Example Usage

```hcl
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "organization_1"
}

resource opentelekomcloud_swr_organization_permissions_v2 user_1 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name

  user_id  = var.user_id
  username = var.username
  auth     = 3
}
```

## Argument Reference

The following arguments are supported:

* `organization` - (Required) The name of the organization (namespace) to be accessed.

* `user_id` - (Required) The ID of the existing Open Telekom Cloud user.

* `username` - (Required) The username of the existing Open Telekom Cloud user.

* `auth` - (Required) User permission that is configured.
  The value can be `1`, `3`, or `7`. `7` ― manage, `3` ―  write, `1` ― read.

## Attributes Reference

The following attributes are exported:

* `organization` - See Argument Reference above.

* `user_id` - See Argument Reference above.

* `username` - See Argument Reference above.

* `auth` - See Argument Reference above.
