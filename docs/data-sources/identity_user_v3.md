---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_user_v3

Use this data source to get the ID of an OpenTelekomCloud user.

## Example Usage

```hcl
data "opentelekomcloud_identity_user_v3" "user_1" {
  name = "user_1"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid values are `true` and `false`.
  Default value is `true`.

* `name` - (Optional) The name of the user.

* `region` - (Optional) The region this user belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_project_id` - The default project this user belongs to.

* `password_expires_at` - Password expiration date of the user.
