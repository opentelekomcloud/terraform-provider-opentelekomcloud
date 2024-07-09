---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_group_membership_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-group-membership-v3"
description: |-
Manages a IAM Group Membership resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM group membership you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/user_group_management)

# opentelekomcloud_identity_group_membership_v3

Manages a Group Membership resource within OpenTelekomCloud IAM service.

-> **Note:** You _must_ have admin privileges in your OpenTelekomCloud cloud to use this resource.

## Example Usage

```hcl
resource "opentelekomcloud_identity_group_v3" "group_1" {
  name        = "group1"
  description = "This is a test group"
}

resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "user1"
  enabled  = true
  password = "password12345!"
}

resource "opentelekomcloud_identity_user_v3" "user_2" {
  name     = "user2"
  enabled  = true
  password = "password12345!"
}

resource "opentelekomcloud_identity_group_membership_v3" "membership_1" {
  group = opentelekomcloud_identity_group_v3.group_1.id
  users = [opentelekomcloud_identity_user_v3.user_1.id,
  opentelekomcloud_identity_user_v3.user_2.id]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required) The group ID of this membership.

* `users` - (Required) A List of user IDs to associate to the group.

## Attributes Reference

The following attributes are exported:

* `group` - See Argument Reference above.

* `users` - See Argument Reference above.
