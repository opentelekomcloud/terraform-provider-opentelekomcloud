---
layout: "opentelekomcloud"
page_title: "OpentelekomCloud: opentelekomcloud_identity_role_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-role-v3"
description: |-
  Manages a V3 Role resource within OpentelekomCloud Keystone.
---

# opentelekomcloud\_identity\_role_v3

Manages a V3 Role resource within OpentelekomCloud Keystone.

Note: You _must_ have admin privileges in your OpentelekomCloud cloud to use
this resource.

## Example Usage

```hcl
resource "opentelekomcloud_identity_role_v3" "role_1" {
  name = "role_1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the role.

* `domain_id` - (Optional) The domain the role belongs to.

* `region` - (Optional) The region in which to obtain the IAM client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new Role.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `region` - See Argument Reference above.

## Import

Roles can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_identity_role_v3.role_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
