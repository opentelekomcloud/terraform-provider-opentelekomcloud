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

* `name` - (Required) The name of the user. The user name consists of 5 to 32
  characters. It can contain only uppercase letters, lowercase letters,
  digits, spaces, and special characters (-_) and cannot start with a digit.

* `description` - (Optional) The description of the Identity User.

* `default_project_id` - (Optional) The default project this user belongs to.

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid
  values are `true` and `false`.

* `password` - (Optional) The password for the user. It must contain at least
  two of the following character types: uppercase letters, lowercase letters,
  digits, and special characters.

* `email` - (Optional) The email associated with user.

* `send_welcome_email` - (Optional) Whether to send a `Welcome Email` or not.
  Possible values are `true` and `false`.

-> Welcome Email will be sent when email is set/changed and `send_welcome_email` is set to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domain_id` - See Argument Reference above.

## Import

Users can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
