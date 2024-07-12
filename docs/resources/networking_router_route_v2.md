---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_router_route_v2"
sidebar_current: "docs-opentelekomcloud-resource-networking-router-route-v2"
description: |-
  Manages a VPC Router Route resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC router route you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/router)

# opentelekomcloud_networking_router_route_v2

Creates a routing entry on a OpenTelekomCloud V2 router.

## Example Usage

```hcl
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_networking_router_route_v2" "router_route_1" {
  depends_on       = ["opentelekomcloud_networking_router_interface_v2.int_1"]
  router_id        = opentelekomcloud_networking_router_v2.router_1.id
  destination_cidr = "10.0.1.0/24"
  next_hop         = "192.168.199.254"
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required) ID of the router this routing entry belongs to. Changing
  this creates a new routing entry.

* `destination_cidr` - (Required) CIDR block to match on the packet’s destination IP. Changing
  this creates a new routing entry.

* `next_hop` - (Required) IP address of the next hop gateway.  Changing
  this creates a new routing entry.

## Attributes Reference

The following attributes are exported:

* `router_id` - See Argument Reference above.

* `destination_cidr` - See Argument Reference above.

* `next_hop` - See Argument Reference above.

-> **Note:** The `next_hop` IP address must be directly reachable from the router at the `opentelekomcloud_networking_router_route_v2`
  resource creation time.  You can ensure that by explicitly specifying a dependency on the `opentelekomcloud_networking_router_interface_v2`
  resource that connects the next hop to the router, as in the example above.
