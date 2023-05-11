---
subcategory: "Virtual Private Cloud (VPC)"
---

Up-to-date reference of API arguments for VPC floating ip you can get at
`https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/floating_ip_address`.

# opentelekomcloud_networking_floatingip_associate_v2

Associates a floating IP to a port. This is useful for situations
where you have a pre-allocated floating IP or are unable to use the
[`resource/opentelekomcloud_networking_floatingip_v2`](networking_floatingip_v2.md) to create a floating IP.

## Example Usage

### Basic FloatingIP associate

```hcl
resource "opentelekomcloud_networking_port_v2" "port_1" {
  network_id = "a5bbd213-e1d3-49b6-aed1-9df60ea94b9a"
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = "1.2.3.4"
  port_id     = opentelekomcloud_networking_port_v2.port_1.id
}
```

### Associate an instance with `port_id`

```hcl
variable "keypair" {}
variable "image_id" {}
variable "network_name" {}

resource "opentelekomcloud_networking_floatingip_v2" "this" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_compute_instance_v2" "this" {
  name            = "example-instance"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = var.keypair
  security_groups = ["default"]

  network {
    name = var.network_name
  }
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "this" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.this.address
  port_id     = opentelekomcloud_compute_instance_v2.this.network.0.port
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip` - (Required) IP Address of an existing floating IP.

* `port_id` - (Required) ID of an existing port with at least one IP address to
  associate with this floating IP.

## Attributes Reference

The following attributes are exported:

* `floating_ip` - See Argument Reference above.

* `port_id` - See Argument Reference above.

## Import

Floating IP associations can be imported using the `id` of the floating IP, e.g.

```sh
terraform import opentelekomcloud_networking_floatingip_associate_v2.fip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
