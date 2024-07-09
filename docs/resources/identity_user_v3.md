---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_user_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-user-v3"
description: |-
Manages a IAM User resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM user you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/user_management)

# opentelekomcloud_identity_user_v3

Manages a User resource within OpenTelekomCloud IAM service.

-> You need to have admin privileges in your OpenTelekomCloud cloud to use
this resource.

## Example Usage

```hcl
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name     = "user_1"
  password = "password123!"
}
```

## Example with login protection

```hcl
resource "opentelekomcloud_identity_user_v3" "user_1" {
  name               = "user_protected"
  password           = "password123@!"
  enabled            = true
  email              = "test@acme.org"
  send_welcome_email = true

  login_protection {
    enabled             = true
    verification_method = "email"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the user. The user name consists of 5 to 32 characters. It can
  contain only uppercase letters, lowercase letters, digits, spaces, and special characters (-_) and cannot start with a
  digit.

* `description` - (Optional, String) Specifies the description of the user.

* `email` - (Optional, String) Specifies the email address with a maximum of 255 characters.

* `phone` - (Optional, String) Specifies the mobile number with a maximum of 32 digits. This parameter must be used
  together with `country_code`.

* `country_code` - (Optional, String) Specifies the country code. This parameter must be used together with `phone`.

* `password` - (Optional, String) Specifies the password for the user with 6 to 32 characters. It must contain at least
  two of the following character types: uppercase letters, lowercase letters, digits, and special characters.

* `pwd_reset` - (Optional, Bool) Specifies whether the password should be reset. By default, the password is asked
  to reset at the first login.

* `enabled` - (Optional, Bool) Specifies whether the user is enabled or disabled. Valid values are `true` and `false`.

* `access_type` - (Optional, String) Specifies the access type of the user. Available values are:
  + **default**: support both programmatic and management console access.
  + **programmatic**: only support programmatic access.
  + **console**: only support management console access.

* `send_welcome_email` - (Optional, Bool) Whether to send a `Welcome Email` or not.
  Possible values are `true` and `false`.

-> Welcome Email will be sent when email is set/changed and `send_welcome_email` is set to `true`.

* `login_protection` - (Optional, List) Login protection configuration.
  The `login_protection` block supports:
  + `enabled` - (Required, Bool) Indicates whether login protection has been enabled for the user. The value can be `true` or `false`.
  + `verification_method` - (Required, String) Login authentication method of the user. Options: `sms`, `email`, and `vmfa`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `password_strength` - Indicates the password strength.

* `create_time` - The time when the IAM user was created.

* `last_login` - The time when the IAM user last login.

* `domain_id` - The domain user belongs to.

* `xuser_type` - Type of the user in the external system.

* `xuser_id` - ID of the user in the external system.

## Import

Users can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```

Due to the security reasons, `password` can not be imported. It can be ignored as shown below.

```hcl
resource "opentelekomcloud_identity_user_v3" "user_1" {
  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
