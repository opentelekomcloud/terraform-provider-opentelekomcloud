---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_floatingip_associate_v2"
sidebar_current: "docs-opentelekomcloud-resource-networking-floatingip-associate-v2"
description: |-
  Associates a Floating IP to a Port
---

# opentelekomcloud\_networking\_floatingip\_associate_v2

Associates a floating IP to a port. This is useful for situations
where you have a pre-allocated floating IP or are unable to use the
`opentelekomcloud_networking_floatingip_v2` resource to create a floating IP.

## Example Usage

```hcl
resource "opentelekomcloud_networking_port_v2" "port_1" {
  network_id = "a5bbd213-e1d3-49b6-aed1-9df60ea94b9a"
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = "1.2.3.4"
  port_id = "${opentelekomcloud_networking_port_v2.port_1.id}"
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

```
$ terraform import opentelekomcloud_networking_floatingip_associate_v2.fip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
