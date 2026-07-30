---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_user_group_membership_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-user-group-membership-v3"
description: |-
 Get detail of group memberships of an IAM user within T Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for IAM user group membership you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/user_group_management)

# opentelekomcloud_identity_user_group_membership_v3

Get detail of group memberships of an IAM user within T Cloud Public.

-> **Note:** You _must_ have `Security Administrator` privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).

## Example Usage

```hcl
variable user_id {}

resource "opentelekomcloud_identity_user_group_membership_v3" "membership_1" {
  user = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) ID of the user.

## Attributes Reference

The following attributes are exported:

* `groups` - IDs of the groups of which the user is a member.
