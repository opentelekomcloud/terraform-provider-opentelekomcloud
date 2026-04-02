---
subcategory: "Software Repository for Container (SWR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_swr_policy_v2"
sidebar_current: "docs-opentelekomcloud-data-source-swr-policy-v2"
description: |-
  Get details of an SWR image retention policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SWR image retention policy you can get at
[documentation portal](https://docs.otc.t-systems.com/software-repository-container/api-ref/apis/image_retention_policy_management/index.html#swr-02-0094)

# opentelekomcloud_swr_policy_v2

Get details of an SWR image retention policy resource within Open Telekom Cloud.

## Example Usage

```hcl
variable "org_name" {}
variable "repo_name" {}
variable "policy_id" {}

data "opentelekomcloud_swr_policy_v2" "policy_1" {
  organization = var.org_name
  repository   = var.repo_name
  policy_id    = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `organization` - (Required, String) Specifies the name of the repository organization.

* `repository` - (Required, String) Specifies the name of the image repository.

* `policy_id` - (Required, String) Specifies the ID of the SWR policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `algorithm` - The image retention policy matching rule.

* `rules` - The image retention policy. The [rules](#rules) structure is documented below.

<a name="rules"></a>
The `rules` block supports:

* `template` - The image retention policy type.
* `params` - The image retention policy parameters.
* `tag_selector` - The exception images. The structure is documented below.
    * `kind` - The matching rule.
    * `pattern` - The matching rule value.
