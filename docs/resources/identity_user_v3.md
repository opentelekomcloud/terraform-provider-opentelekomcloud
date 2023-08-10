---
subcategory: "Identity and Access Management (IAM)"
---

Up-to-date reference of API arguments for IAM user you can get at
`https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/user_management`.

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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the user. The user name consists of 5 to 32 characters. It can
  contain only uppercase letters, lowercase letters, digits, spaces, and special characters (-_) and cannot start with a
  digit.

* `description` - (Optional) Specifies the description of the user.

* `email` - (Optional) Specifies the email address with a maximum of 255 characters.

* `phone` - (Optional) Specifies the mobile number with a maximum of 32 digits. This parameter must be used
  together with `country_code`.

* `country_code` - (Optional) Specifies the country code. The country code of the Chinese mainland is 0086. This
  parameter must be used together with `phone`.

* `password` - (Optional) Specifies the password for the user with 6 to 32 characters. It must contain at least
  two of the following character types: uppercase letters, lowercase letters, digits, and special characters.

* `pwd_reset` - (Optional) Specifies whether the password should be reset. By default, the password is asked
  to reset at the first login.

* `enabled` - (Optional) Specifies whether the user is enabled or disabled. Valid values are `true` and `false`.

* `access_type` - (Optional) Specifies the access type of the user. Available values are:
  + **default**: support both programmatic and management console access.
  + **programmatic**: only support programmatic access.
  + **console**: only support management console access.

* `send_welcome_email` - (Optional) Whether to send a `Welcome Email` or not.
  Possible values are `true` and `false`.

-> Welcome Email will be sent when email is set/changed and `send_welcome_email` is set to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `password_strength` - Indicates the password strength.

* `create_time` - The time when the IAM user was created.

* `last_login` - The time when the IAM user last login.

## Import

Users can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```

Due to the security reasons, `password` can not be imported. It can be ignored as shown below.

```hcl
resource "huaweicloud_identity_user" "user_1" {
  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
