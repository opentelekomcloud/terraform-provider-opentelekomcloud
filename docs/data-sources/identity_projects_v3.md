---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_projects_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-projects-v3"
description: |-
  Get a IAM projects information from OpenTelekomCloud
---

Up-to-date reference of API arguments for IAM projects you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/project_management/querying_project_information_based_on_the_specified_criteria.html#en-us-topic-0057845625)

# opentelekomcloud_identity_projects_v3

Use this data source to get the list of all OpenTelekomCloud projects.

## Example Usage

```hcl
data "opentelekomcloud_identity_projects_v3" "all" {}
```


## Argument Reference

Data resource lists all available project therefore no arguments are provided.

## Attributes Reference

* `id` - Indicates the domain of queried projects.

* `region` - Indicates the region of queried projects.

* `projects` - List of projects details. The object structure of each Project is documented below.

The `projects` block supports:

* `region` - Indicates the region where the project is present.

* `name` - Indicated the name of the project.

* `description` - The description of the project.

* `domain_id` - The domain this project belongs to.

* `project_id` - The ID of the project.

* `parent_id` - The parent of this project.

* `enabled` - Describes whether the project is available

* `is_domain` - Indicates whether the user calling the API is a tenant.
