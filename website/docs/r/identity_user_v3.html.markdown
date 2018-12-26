---
layout: "opentelekomcloud"
page_title: "OpentelekomCloud: opentelekomcloud_identity_user_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-user-v3"
description: |-
  Manages a User resource within OpentelekomCloud Keystone.
---

# opentelekomcloud\_identity\_user_v3

Manages a User resource within OpentelekomCloud Keystone.

Note: You _must_ have admin privileges in your OpentelekomCloud cloud to use
this resource.

## Example Usage

```hcl
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "opentelekomcloud_identity_user_v3" "user_1" {
  default_project_id = "${opentelekomcloud_identity_project_v3.project_1.id}"
  name = "user_1"
  description = "A user"

  password = "password123"

}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the user.

* `default_project_id` - (Optional) The default project this user belongs to.

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid
  values are `true` and `false`.

* `name` - (Optional) The name of the user.

* `password` - (Optional) The password for the user.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new User.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.

## Import

Users can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
