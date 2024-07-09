---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_group_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-group-v3"
description: |-
Get a IAM group information from OpenTelekomCloud
---

Up-to-date reference of API arguments for IAM group you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/user_group_management/listing_user_groups.html#en-us-topic-0057845602)

# opentelekomcloud_identity_group_v3
Use this data source to get the ID of an OpenTelekomCloud group.

-> **Note:** This usually requires admin privileges.

## Example Usage

```hcl
data "opentelekomcloud_identity_group_v3" "admins" {
  name = "admins"
}
```

## Argument Reference

* `name` - (Required) The name of the group.

* `domain_id` - (Optional) The domain the group belongs to.


## Attributes Reference

`id` is set to the ID of the found group. In addition, the following attributes are exported:

* `name` - See Argument Reference above.

* `domain_id` - See Argument Reference above.
