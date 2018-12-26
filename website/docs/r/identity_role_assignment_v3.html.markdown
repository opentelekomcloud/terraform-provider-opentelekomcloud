---
layout: "opentelekomcloud"
page_title: "OpentelekomCloud: opentelekomcloud_identity_role_assignment_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-role-assignment-v3"
description: |-
  Manages a V3 Role assignment within OpentelekomCloud Keystone.
---

# opentelekomcloud\_identity\_role\_assignment_v3

Manages a V3 Role assignment within OpentelekomCloud Keystone.

Note: You _must_ have admin privileges in your OpentelekomCloud cloud to use
this resource.

## Example Usage

```hcl
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "opentelekomcloud_identity_user_v3" "user_1" {
  name = "user_1"
  default_project_id = "${opentelekomcloud_identity_project_v3.project_1.id}"
}

resource "opentelekomcloud_identity_role_v3" "role_1" {
  name = "role_1"
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
  user_id = "${opentelekomcloud_identity_user_v3.user_1.id}"
  project_id = "${opentelekomcloud_identity_project_v3.project_1.id}"
  role_id = "${opentelekomcloud_identity_role_v3.role_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional; Required if `project_id` is empty) The domain to assign the role in.

* `group_id` - (Optional; Required if `user_id` is empty) The group to assign the role to.

* `project_id` - (Optional; Required if `domain_id` is empty) The project to assign the role in.

* `user_id` - (Optional; Required if `group_id` is empty) The user to assign the role to.

* `role_id` - (Required) The role to assign.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `group_id` - See Argument Reference above.
* `user_id` - See Argument Reference above.
* `role_id` - See Argument Reference above.
