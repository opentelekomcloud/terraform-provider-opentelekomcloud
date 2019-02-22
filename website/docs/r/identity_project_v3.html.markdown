---
layout: "opentelekomcloud"
page_title: "OpentelekomCloud: opentelekomcloud_identity_project_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-project-v3"
description: |-
  Manages a Project resource within OpentelekomCloud Keystone.
---

# opentelekomcloud\_identity\_project_v3

Manages a Project resource within OpentelekomCloud Identity And Access 
Management service.

Note: You _must_ have security admin privileges in your OpentelekomCloud 
cloud to use this resource. please refer to [User Management Model](
https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html)

## Example Usage

```hcl
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name = "eu-de_project1"
  description = "This is a test project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project. it must start with 
    ID of an existing region_ and be less than or equal to 64 characters.
    Example: eu-de_project1.

* `description` - (Optional) A description of the project.

* `domain_id` - (Optional) The domain this project belongs to. Changing this
    creates a new Project.

* `parent_id` - (Optional) The parent of this project. Changing this creates
    a new Project.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `parent_id` - See Argument Reference above.

## Import

Projects can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_identity_project_v3.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
