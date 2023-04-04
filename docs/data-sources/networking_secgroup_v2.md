---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_networking_secgroup_v2

Use this data source to get the ID of an available OpenTelekomCloud security group.

## Example Usage

```hcl
data "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Example Filter by regex

```hcl
data "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name_regex  = "^secgroup_1.+"
}
```

## Argument Reference

* `secgroup_id` - (Optional) The ID of the security group.

* `name` - (Optional) The name of the security group.

* `name_regex` - (Optional) A regex string to apply to the security group list.
  This allows more advanced filtering not supported from the OpenTelekomCloud API.
  This filtering is done locally on what OpenTelekomCloud returns.

* `tenant_id` - (Optional) The owner of the security group.

## Attributes Reference

`id` is set to the ID of the found security group. In addition, the following attributes are exported:

* `name` - See Argument Reference above.

* `description`- The description of the security group.
