---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_role_assignment_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-role-assignment-v3"
description: |-
  Manages a IAM Role Assignment resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM role assignment you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/permission_management)

# opentelekomcloud_identity_role_assignment_v3

Manages a V3 Role assignment within group on OpenTelekomCloud IAM Service.

-> **Note:** You _must_ have admin privileges in your OpenTelekomCloud cloud to use this resource.

## Example Usage

### Assign Role On Project Level

```hcl
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name = "eu-de_project_1"
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "group_1"
}

data "opentelekomcloud_identity_role_v3" "role_1" {
  name = "system_all_4" #ECS admin
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id   = opentelekomcloud_identity_group_v3.group_1.id
  project_id = opentelekomcloud_identity_project_v3.project_1.id
  role_id    = data.opentelekomcloud_identity_role_v3.role_1.id
}
```

### Assign Role On Domain Level

```hcl
variable "domain_id" {
  default     = "01aafcf63744d988ebef2b1e04c5c34"
  description = "this is the domain id"
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "group_1"
}

data "opentelekomcloud_identity_role_v3" "role_1" {
  name = "secu_admin" #security admin
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id  = opentelekomcloud_identity_group_v3.group_1.id
  domain_id = var.domain_id
  role_id   = data.opentelekomcloud_identity_role_v3.role_1.id
}
```

### Assign Role for All Projects (existing and future)

```hcl
variable "domain_id" {
  default     = "01aafcf63744d988ebef2b1e04c5c34"
  description = "this is the domain id"
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
  name = "group_1"
}

data "opentelekomcloud_identity_role_v3" "role_1" {
  name = "secu_admin" #security admin
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id     = opentelekomcloud_identity_group_v3.group_1.id
  domain_id    = var.domain_id
  role_id      = data.opentelekomcloud_identity_role_v3.role_1.id
  all_projects = true
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional; Required if `project_id` is empty) The domain to assign the role in.

* `group_id` - (Required) The group to assign the role to.

* `project_id` - (Optional; Required if `domain_id` is empty) The project to assign the role in.

* `role_id` - (Required) The role to assign.

* `all_projects` - (Optional) Whether to assign role for all existing and future projects.
  `domain_id` has to be specified if `all_projects` is set to `true`.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.

* `project_id` - See Argument Reference above.

* `group_id` - See Argument Reference above.

* `role_id` - See Argument Reference above.
