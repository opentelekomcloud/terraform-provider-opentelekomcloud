---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_secgroup_v2"
sidebar_current: "docs-opentelekomcloud-datasource-networking-secgroup-v2"
description: |-
  Get information on an OpenTelekomCloud Security Group.
---

# opentelekomcloud\_networking\_secgroup\_v2

Use this data source to get the ID of an available OpenTelekomCloud security group.

## Example Usage

```hcl
data "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Argument Reference

* `secgroup_id` - (Optional) The ID of the security group.

* `name` - (Optional) The name of the security group.

* `tenant_id` - (Optional) The owner of the security group.

## Attributes Reference

`id` is set to the ID of the found security group. In addition, the following
attributes are exported:

* `name` - See Argument Reference above.
* `description`- The description of the security group.
