---
subcategory: "Software Repository for Container (SWR)"
---

Up-to-date reference of API arguments for SWR domain you can get at
`https://docs.otc.t-systems.com/software-repository-container/api-ref/api`.

# opentelekomcloud_swr_domain_v2

Manages the SWR image sharing domain resource within Open Telekom Cloud.

## Example Usage

```hcl
variable "access_domain" {}

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

resource opentelekomcloud_swr_domain_v2 domain_1 {
  organization  = opentelekomcloud_swr_organization_v2.org_1.name
  repository    = opentelekomcloud_swr_organization_v2.repo_1.name
  access_domain = var.access_domain
  permission    = "read"
  deadline      = "forever"
}
```

## Argument Reference

The following arguments are supported:

* `organization` - (Required) The name of the repository organization.

* `repository` - (Required) The name of the repository.

* `access_domain` - (Required) The name of the domain for image sharing.

-> `access_domain` should be an existing OTC domain.

* `permission` - (Required) Permission to be granted. Currently, only the `read` permission is supported.

* `deadline` - (Required) End date of image sharing (UTC). When the value is set to `forever`,
  the image will be permanently available for the domain. The validity period is calculated by day.
  The shared images expire at `00:00:00` on the day after the end date.

* `description` - (Optional) Specifies SWR domain description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `creator_id` - Username ID of the domain creator.

* `creator_name` - Username of the domain creator.

* `created` - Indicates the creation time.

* `updated` - Indicates the domain when was last updated.

* `status` - Indicates the domain is valid (`true`) or expired (`false`).
