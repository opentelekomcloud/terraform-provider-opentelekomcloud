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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `password_expires_at` - Password expiration date of the user.

* `mfa_device` - Serial number of user MFA device.
  `Security administrator` permissions are needed to set this attribute.
