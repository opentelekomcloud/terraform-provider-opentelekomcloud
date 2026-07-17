---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_user_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-user-v3"
description: |-
  Get a IAM user information from OpenTelekomCloud
---

Up-to-date reference of API arguments for IAM user you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/user_management/querying_a_user_list.html#en-us-topic-0057845638)

# opentelekomcloud_identity_user_v3

Use this data source to get the ID of an OpenTelekomCloud user.

## Example Usage

### Query by name

```hcl
data "opentelekomcloud_identity_user_v3" "user_1" {
  name = "user_1"
}
```

### Query by ID

```hcl
data "opentelekomcloud_identity_user_v3" "user_1" {
  id = "0a54c86a-0b3c-4762-92cc-bdfb63c572e1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the user. If set, it takes precedence over the arguments below
  and the user is looked up by ID directly.

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid values are `true` and `false`.
  Default value is `true`.

* `name` - (Optional) The name of the user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the user.

* `name` - The name of the user.

* `password_expires_at` - Password expiration date of the user.

* `mfa_device` - Serial number of user MFA device.
  `Security administrator` permissions are needed to set this attribute.
