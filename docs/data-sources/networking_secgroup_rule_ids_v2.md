---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_networking_secgroup_rule_ids_v2

Use this data source to get a list of security group rules ids for a `security_group_id`.

This resource can be useful for getting back a list of security group rules ids for a Security Group.

## Example Usage

The following example shows outputting all security group rules for security group.

```hcl
variable "security_group_id" {}

data "opentelekomcloud_networking_secgroup_rule_ids_v2" "sg_ids" {
  security_group_id = var.security_group_id
}

output "secgroup_rule_ids" {
  value = [for id in data.opentelekomcloud_networking_secgroup_rule_ids_v2.sg_ids.ids : id]
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required) Specifies the security group ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of all the security group rule IDs found. This data source will fail if none are found.
