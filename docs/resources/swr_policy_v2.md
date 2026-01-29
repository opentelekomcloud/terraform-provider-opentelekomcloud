---
subcategory: "Software Repository for Container (SWR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_swr_policy_v2"
sidebar_current: "docs-opentelekomcloud-resource-swr-policy-v2"
description: |-
  Manages an SWR image retention policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SWR image retention policy you can get at
[documentation portal](https://docs.otc.t-systems.com/software-repository-container/api-ref/apis/image_retention_policy_management/index.html#swr-02-0094)

# opentelekomcloud_swr_policy_v2

Manages the SWR image retention policy resource within Open Telekom Cloud.

## Example Usage

```hcl
variable "org_name" {}
variable "repo_name" {}

resource "opentelekomcloud_swr_policy_v2" "policy_1" {
  organization = var.org_name
  repository   = var.repo_name
  algorithm    = "or"
  rules {
    template = "date_rule"
    params = {
      days = "30"
    }
    tag_selector {
      kind    = "label"
      pattern = "v1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `organization` - (Required, String, ForceNew) Specifies the name of the repository organization.

* `repository` - (Required, String, ForceNew) Specifies the name of the image repository.

* `algorithm` - (Required, String) Specifies the image retention policy matching rule. Accepted value: `or`.

* `rules` - (Required, List) Specifies the image retention policy. The [rules](#rules) structure is documented below.

<a name="rules"></a>
The `rules` block supports:

* `template` - (Required, String) Specifies the image retention policy type. Acceped values: `date_rule`, `tag_rule`.

* `params` - (Required, Map) Specifies the image retention policy parameters. If `template` is set to `date_rule`, set `params` to `{ days = "xxx" }`. If template is set to `tag_rule`, set `params` to `{ num = "xxx" }`.

* `tag_selector` - (Required, List) Specifies the exception images. The structure is documented below.
    * `kind` - (Required, String) Specifies the matching rule. Accepted values: `label`, `regexp`.
    * `pattern` - (Required, String) Specifies the matching rule value. If `kind` is set to `label`, set `pattern` to the `<image tag>`, e.g. `"v1"`. If `kind` is set to `regexp`, set `pattern` to a `<regular expression>`, e.g. `"^123$"`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the image retention policy.

## Import

SWR image retention policy can be imported using the organization name `<org_name>`, repository name `<repo_name>`, and policy ID `<id>`, e.g.

```shell
terraform import opentelekomcloud_swr_policy_v2.policy_1 <org_name>/<repo_name>/<id>
```
