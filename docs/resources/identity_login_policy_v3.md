---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_login_policy_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-login-policy-v3"
description: |-
  Manages a IAM Login Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM provider you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/security_settings/modifying_the_password_policy.html)

# opentelekomcloud_identity_login_policy_v3

Manages the account login authentication policy within OpenTelekomCloud.

`Please use it with care!`
-> You _must_ have security admin privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).

  During action `terraform destroy` it sets values the same as defaults for this resource.
  Defaults is:
  + `custom_info_for_login` = ""
  + `period_with_login_failures` = 60
  + `lockout_duration` = 15
  + `account_validity_period` = 0
  + `login_failed_times` = 3
  + `session_timeout` = 1395
  + `show_recent_login_info` = false

## Example Usage

```hcl
resource "opentelekomcloud_identity_login_policy_v3" "policy_1" {
  custom_info_for_login      = ""
  period_with_login_failures = 60
  lockout_duration           = 15
  account_validity_period    = 0
  login_failed_times         = 3
  session_timeout            = 1395
  show_recent_login_info     = false
}
```

## Argument Reference

The following arguments are supported:

* `account_validity_period` - (Optional, Int) Validity period (days) to disable users if they have not logged in within the period.
  Value range: `0-240`. If this parameter is set to `0`, no users will be disabled. Default: `0`.

* `custom_info_for_login` - (Optional, String) Custom information that will be displayed upon successful login.

* `lockout_duration` - (Optional, Int) Duration (minutes) to lock users out. Value range: `15-30`.

* `login_failed_times` - (Optional, Int) Number of unsuccessful login attempts to lock users out. Value range: `3-10`.

* `period_with_login_failures` - (Optional, Int) Period (minutes) to count the number of unsuccessful login attempts.
  Value range: `15-60`.

* `session_timeout` - (Optional, Int) Session timeout (minutes) that will apply if you or users created using your
  account do not perform any operations within a specific period. Value range: `15-1440`.

* `show_recent_login_info` - (Optional, Bool) Indicates whether to display last login information upon successful login.
  The value can be `true` or `false`. Default: `true`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of account login policy, which is the same as the domain ID.

## Import

Identity login authentication policy can be imported using the account ID or domain ID, e.g.

```bash
$ terraform import opentelekomcloud_identity_login_policy_v3.example <ID>
```
