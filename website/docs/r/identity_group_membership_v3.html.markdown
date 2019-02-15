---
layout: "opentelekomcloud"
page_title: "OpentelekomCloud: opentelekomcloud_identity_group_membership_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-group-membership-v3"
description: |-
  Manages the membership combine User Group resource and User resource  within
  OpentelekomCloud IAM service.
---

# opentelekomcloud\_identity\_group_membership_v3

Manages a User Group Membership resource within OpentelekomCloud IAM service.

Note: You _must_ have admin privileges in your OpentelekomCloud cloud to use
this resource.

## Example Usage

```hcl
resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "group1"
  description = "This is a test group"
}

resource "opentelekomcloud_identity_user_v3" "user_1" {
      name = "user1"
      enabled = true
      password = "password12345!"
}

resource "opentelekomcloud_identity_user_v3" "user_2" {
      name = "user2"
      enabled = true
      password = "password12345!"
}

resource "opentelekomcloud_identity_group_membership_v3" "membership_1" {
        group = "${opentelekomcloud_identity_group_v3.group_1.id}"
        users = ["${opentelekomcloud_identity_user_v3.user_1.id}",
                "${opentelekomcloud_identity_user_v3.user_2.id}"
                ]
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

