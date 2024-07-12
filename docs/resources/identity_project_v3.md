---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_project_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-project-v3"
description: |-
  Manages a IAM Project resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM project you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/project_management)

# opentelekomcloud_identity_project_v3

Manages a Project resource within OpenTelekomCloud Identity And Access
Management service.

-> **Note:** You _must_ have security admin privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).

## Example Usage

```hcl
resource "opentelekomcloud_identity_project_v3" "project_1" {
  name        = "eu-de_project1"
  description = "This is a test project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project. it must start with
  ID of an existing region and be less than or equal to 64 characters.
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

```sh
terraform import opentelekomcloud_identity_project_v3.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
