---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_user_group_membership_v3

Manages a User Group Membership resource within OpenTelekomCloud IAM service.

-> **Note:** You _must_ have admin privileges in your OpenTelekomCloud cloud to use this resource.

## Example Usage

```hcl
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "user-1"
  password = "password123@#"
  enabled  = true
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "group-1"
}

resource "opentelekomcloud_identity_group_v3" "group_2" {
  name = "group-2"
}

resource "opentelekomcloud_identity_user_group_membership_v3" "membership_1" {
  user = opentelekomcloud_identity_user_v3.user_1.id
  groups = [
    opentelekomcloud_identity_group_v3.group_1.id,
    opentelekomcloud_identity_group_v3.group_2.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) ID of a user.

* `groups` - (Required) IDs of the groups for the user to be assigned to.

## Attributes Reference

The following attributes are exported:

* `user` - See Argument Reference above.

* `groups` - See Argument Reference above.
