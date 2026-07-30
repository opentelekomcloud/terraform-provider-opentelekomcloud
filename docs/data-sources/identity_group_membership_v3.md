---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_group_membership_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-group-membership-v3"
description: |-
  Get list of members in an IAM group within T Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for IAM group membership you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/user_group_management)

# opentelekomcloud_identity_group_membership_v3

Get list of members in an IAM group within T Cloud Public IAM service.

-> **Note:** You _must_ have `Security Administrator` privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).

## Example Usage

```hcl
variable group_id {}

data "opentelekomcloud_identity_group_membership_v3" "membership_1" {
  group = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String) The group ID.

## Attributes Reference

The following attributes are exported:

* `users` -  List of user IDs of group members.
