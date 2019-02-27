---
layout: "opentelekomcloud"
page_title: "OpentelekomCloud: opentelekomcloud_identity_user_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-user-v3"
description: |-
  Get information on an OpentelekomCloud User.
---

# opentelekomcloud\_identity\_user_v3

Use this data source to get the ID of an OpentelekomCloud user.

## Example Usage

```hcl
data "opentelekomcloud_identity_user_v3" "user_1" {
  name = "user_1"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the user.

* `default_project_id` - (Optional) The default project this user belongs to.

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid
  values are `true` and `false`.

* `idp_id` - (Optional) The identity provider ID of the user.

* `name` - (Optional) The name of the user.

* `password_expires_at` - (Optional) Query for expired passwords. See the [OpentelekomCloud API docs](https://docs.otc.t-systems.com/en-us/api/iam/en-us_topic_0057845638.html) for more information on the query format.


## Attributes Reference

The following attributes are exported:

* `default_project_id` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `enabled` - See Argument Reference above.
* `idp_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `password_expires_at` - See Argument Reference above.
