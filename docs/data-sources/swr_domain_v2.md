---
subcategory: "Software Repository for Container (SWR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_swr_domain_v2"
sidebar_current: "docs-opentelekomcloud-datasource-swr-domain-v2"
description: |-
  Get details of an SWR Domain resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SWR domain you can get at
[documentation portal](https://docs.otc.t-systems.com/software-repository-container/api-ref/api)

# opentelekomcloud_swr_domain_v2

Get details of SWR image sharing domain resource within Open Telekom Cloud.

## Example Usage

```hcl
variable "org_name" {}
variable "repo_name" {}
variable "access_domain" {}

data opentelekomcloud_swr_domain_v2 domain_1 {
  organization  = var.org_name
  repository    = var.repo_name
  access_domain = var.access_domain
}
```

## Argument Reference

The following arguments are supported:

* `organization` - (Required) The name of the repository organization.

* `repository` - (Required) The name of the repository.

* `access_domain` - (Required) The name of the domain for image sharing.

-> `access_domain` should be an existing OTC domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `permission` - Permission to be granted.

* `deadline` - End date of image sharing (UTC). 

* `description` - Specifies SWR domain description.

* `creator_id` - Username ID of the domain creator.

* `creator_name` - Username of the domain creator.

* `created` - Indicates the creation time.

* `updated` - Indicates the domain when was last updated.

* `status` - Indicates the domain is valid (`true`) or expired (`false`).
