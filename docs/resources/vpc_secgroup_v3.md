---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_secgroup_v3"
sidebar_current: "docs-opentelekomcloud-resource-vpc-secgroup-v3"
description: |-
  Manages a VPC Security Group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC security group v3 you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/api_v3/security_group/index.html)

# opentelekomcloud_vpc_secgroup_v3

Manages a V3 VPC security group resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_vpc_secgroup_v3" "secgroup_1" {
  name        = "secgroup_1"
  description = "My VPC security group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) A unique name for the security group.

* `description` - (Optional, String) Security group description.

* `enterprise_project_id` - (Optional, String, ForceNew) ID of the enterprise project to which the security group belongs. Default: `0`.

* `tags` - (Optional, Map, ForceNew) Specifies the tags.

* `delete_default_rules` - (Optional, Boolean, ForceNew) Specify whether to delete all default rules when security group is created. Default: `false`.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Security Group ID.

* `project_id` - Indicates the project ID.

* `created_at` - Indicates the time when the security group was created. It is a UTC time in yyyy-MM-ddTHH:mm:ssZ format.

* `updated_at` - Indicates the time when the security group was updated. It is a UTC time in yyyy-MM-ddTHH:mm:ssZ format.

## Import

VPC Security Group V3 can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_secgroup_v3.secgroup_1 <id>
```

## Notes

It is required to ignore changes as below:

```hcl
resource "opentelekomcloud_vpc_secgroup_v3" "group_1" {
  # ...

  lifecycle {
    ignore_changes = [
      delete_default_rules,
    ]
  }
}
```
