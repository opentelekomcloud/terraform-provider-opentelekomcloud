---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_secgroup_v2"
sidebar_current: "docs-opentelekomcloud-resource-networking-secgroup-v2"
description: |-
  Manages a VPC Security Group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC security group you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/security_group)

# opentelekomcloud_networking_secgroup_v2

Manages a V2 neutron security group resource within OpenTelekomCloud.
Unlike Nova security groups, neutron separates the group from the rules
and also allows an admin to target a specific tenant_id.

## Example Usage

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
}
```

## Default Security Group Rules

In most cases, OpenTelekomCloud will create some egress security group rules for each
new security group. These security group rules will not be managed by
Terraform, so if you prefer to have *all* aspects of your infrastructure
managed by Terraform, set `delete_default_rules` to `true` and then create
separate security group rules such as the following:

```hcl
resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_v4" {
  direction         = "egress"
  ethertype         = "IPv4"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_v6" {
  direction         = "egress"
  ethertype         = "IPv6"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
}
```

-> **Note:** This behavior may differ depending on the configuration of
the OpenTelekomCloud cloud. The above illustrates the current default Neutron
behavior. Some OpenTelekomCloud clouds might provide additional rules and some might
not provide any rules at all (in which case the `delete_default_rules` setting
is moot).

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the security group.

* `description` - (Optional) A unique name for the security group.

* `tenant_id` - (Optional) The owner of the security group. Required if admin
  wants to create a port for another tenant. Changing this creates a new
  security group.

* `delete_default_rules` - (Optional) Whether or not to delete the default
  egress security rules. This is `false` by default.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

## Import

Security Groups can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_networking_secgroup_v2.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
