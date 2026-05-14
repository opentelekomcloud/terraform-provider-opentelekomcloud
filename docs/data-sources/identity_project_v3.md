---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_project_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-project-v3"
description: |-
  Get a IAM project information from OpenTelekomCloud
---

Up-to-date reference of API arguments for IAM project you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/project_management/querying_project_information_based_on_the_specified_criteria.html#en-us-topic-0057845625)

# opentelekomcloud_identity_project_v3

Use this data source to get information about an OpenTelekomCloud project by ID or by filter arguments.

## Example Usage

### Query Project by ID

```hcl
data "opentelekomcloud_identity_project_v3" "project_1" {
  id = "<project_id>"
}
```

### Query Project by Name

```hcl
data "opentelekomcloud_identity_project_v3" "project_1" {
  name = "demo"
}
```

### Query Current Project details

If no filter arguments are provided, the data source gets information about the current project.

```hcl
data "opentelekomcloud_identity_project_v3" "project_1" {
}
```


## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the project. When specified, the project is looked up directly by ID.

* `domain_id` - (Optional) The domain this project belongs to.

* `enabled` - (Optional) Whether the project is enabled or disabled. Valid values are `true` and `false`.

* `is_domain` - (Optional) Whether this project is a domain. Valid values are `true` and `false`.

* `name` - (Optional) The name of the project.

* `parent_id` - (Optional) The parent of this project.

## Attributes Reference

`id` is set to the ID of the found project. In addition, the following attributes are exported:

* `description` - The description of the project.

* `domain_id` - ID of an enterprise account to which a project belongs.

* `enabled` - Whether a project is available.

* `is_domain` - Indicates whether the user calling the API is a tenant.

* `name` - Project name.

* `parent_id` - Parent ID of the project.

* `region` - Indicates the region where the project is present.
