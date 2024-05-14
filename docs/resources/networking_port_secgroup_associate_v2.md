---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_networking_port_secgroup_associate_v2

Manages a V2 port's security groups within OpenTelekomCloud. Useful, when the port was
created not by Terraform (e.g. Manila or LBaaS). It should not be used, when the
port was created directly within Terraform.

When the resource is deleted, Terraform doesn't delete the port, but unsets the
list of user defined security group IDs.  However, if `force` is set to `true`
and the resource is deleted, Terraform will remove all assigned security group
IDs.

## Example Usage

```hcl
data "opentelekomcloud_networking_port_v2" "system_port" {
  fixed_ip = "10.0.0.10"
}

data "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = "secgroup"
}

resource "opentelekomcloud_networking_port_secgroup_associate_v2" "port_1" {
  port_id = data.opentelekomcloud_networking_port_v2.system_port.id
  security_group_ids = [
    data.opentelekomcloud_networking_secgroup_v2.secgroup.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
  A networking client is needed to manage a port. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  resource.

* `port_id` - (Required) An UUID of the port to apply security groups to.

* `security_group_ids` - (Required) A list of security group IDs to apply to
  the port. The security groups must be specified by ID and not name (as
  opposed to how they are configured with the Compute Instance).

* `force` - (Optional) Whether to replace or append the list of security
  groups, specified in the `security_group_ids`. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `all_security_group_ids` - The collection of Security Group IDs on the port
  which have been explicitly and implicitly added.

## Import

Port security group association can be imported using the `id` of the port, e.g.

```
$ terraform import opentelekomcloud_networking_port_secgroup_associate_v2.port_1 eae26a3e-1c33-4cc1-9c31-5ght78rdf12
  lifecycle {
    ignore_changes = [
      force,
      security_group_ids,
    ]
  }
