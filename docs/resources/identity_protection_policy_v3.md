---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_protection_policy_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-protection-policy-v3"
description: |-
Manages a IAM Protection Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM provider you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/security_settings/modifying_the_operation_protection_policy.html)

# opentelekomcloud_identity_protection_policy_v3

Manages the IAM operation protection policy within OpenTelekomCloud.

`Please use it with care!`
-> You _must_ have security admin privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).

  During action `terraform destroy` it sets values the same as defaults for this resource.
  Defaults is:
  +  `enable_operation_protection_policy` = false

## Example Usage

```hcl
resource "opentelekomcloud_identity_protection_policy_v3" "policy_1" {
  enable_operation_protection_policy = false
}
```

## Argument Reference

The following arguments are supported:

* `enable_operation_protection_policy` - (Optional, Bool) Indicates whether operation protection has been enabled.
  The value can be `true` or `false`. Default: `false`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of account protection policy, which is the same as the domain ID.

## Import

Identity operation protection policy can be imported using the account ID or domain ID, e.g.

```bash
$ terraform import opentelekomcloud_identity_protection_policy_v3.example <ID>
```
