---
subcategory: "Software Repository for Container (SWR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_swr_repository_v2"
sidebar_current: "docs-opentelekomcloud-resource-swr-repository-v2"
description: |-
  Manages an SWR Repository resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SWR repository you can get at
[documentation portal](https://docs.otc.t-systems.com/software-repository-container/api-ref/api)

# opentelekomcloud_swr_repository_v2

Manages the SWR repository resource within Open Telekom Cloud.

## Example Usage

```hcl
resource opentelekomcloud_swr_organization_v2 org_1 {
  name = "organization_1"
}

resource opentelekomcloud_swr_repository_v2 repo_1 {
  organization = opentelekomcloud_swr_organization_v2.org_1.name
  name         = "repository_1"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}
```

## Argument Reference

The following arguments are supported:

* `organization` - The name of the repository organization.

* `name` - Enter 1 to 128 characters, starting and ending with a lowercase letter or digit.
  Only lowercase letters, digits, periods (`.`), slashes (`/`), underscores (`_`), and hyphens (`-`) are allowed.
  Periods, underscores, and hyphens cannot be placed next to each other.
  A maximum of two consecutive underscores are allowed.

* `is_public` - Whether the repository is public.
  When the value is `true`, it indicates the repository is public.
  When the value is `false`, it indicates the repository is private.

* `description` - (Optional) Repository description.

* `category` - (Optional) Repository category. The value can be `app_server`, `linux`, `framework_app`, `database`,
  `lang`, `other`, `windows`, `arm`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `repository_id` - Numeric ID of the repository.

* `path` - Image address for docker pull.

* `internal_path` - Intra-cluster image address for docker pull.

* `num_images` - Number of image tags in a repository.

* `size` - Repository size.

## Import

Repositories can be imported with `organization/repository`, e.g.

```shell
terraform import opentelekomcloud_swr_repository_v2.repo_1 organization_1/repository_1
```
