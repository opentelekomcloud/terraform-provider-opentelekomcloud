---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_networking_network_v2

Manages a V2 Neutron network resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name               = "port_1"
  network_id         = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up     = "true"
  security_group_ids = [opentelekomcloud_compute_secgroup_v2.secgroup_1.id]

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.10"
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = [opentelekomcloud_compute_secgroup_v2.secgroup_1.name]

  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the network. Changing this updates the name of
  the existing network.

* `shared` - (Optional)  Specifies whether the network resource can be accessed
  by any tenant or not. Changing this updates the sharing capabilities of the
  existing network. Shared SNAT only available in eu-de region.

* `tenant_id` - (Optional) The owner of the network. Required if admin wants to
  create a network for another tenant. Changing this creates a new network.

* `admin_state_up` - (Optional) The administrative state of the network.
  Acceptable values are "true" and "false". Changing this value updates the
  state of the existing network.

* `value_specs` - (Optional) Map of additional options.

* `segments` - (Optional) An array of one or more provider segment objects.

The `segments` block supports:

* `physical_network` - The physical network where this network is implemented.

* `segmentation_id` - An isolated segment on the physical network.

* `network_type` - The type of physical network.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `shared` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

## Import

Networks can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_networking_network_v2.network_1 d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
