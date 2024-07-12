---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_router_v2"
sidebar_current: "docs-opentelekomcloud-resource-networking-router-v2"
description: |-
  Manages a VPC Router resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC router you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/router)

# opentelekomcloud_networking_router_v2

Manages a V2 router resource within OpenTelekomCloud. The router is the top-level resource for the VPC within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "my_router"
  external_gateway = "f67f0d72-0ddf-11e4-9d95-e1f29f417e2f"
}
```

## Enable SNAT
```hcl
data "opentelekomcloud_networking_network_v2" "admin_external_net" {
  name = "admin_external_net"
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "my_router"
  admin_state_up   = true
  distributed      = false
  external_gateway = data.opentelekomcloud_networking_network_v2.admin_external_net.id
  enable_snat      = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) A unique name for the router. Changing this
  updates the `name` of an existing router.

* `admin_state_up` - (Optional) Administrative up/down status for the router
  (must be `true` or `false` if provided). Changing this updates the
  `admin_state_up` of an existing router.

* `distributed` - (Optional) Indicates whether or not to create a
  distributed router. The default policy setting in Neutron restricts
  usage of this property to administrative users only.

* `external_gateway` - (Optional) The network UUID of an external gateway for
  the router. A router with an external gateway is required if any compute
  instances or load balancers will be using floating IPs. Changing this
  updates the `external_gateway` of an existing router.

* `enable_snat` - (Optional) Enable Source NAT for the router. Valid values are
  `true` or `false`. An `external_gateway` has to be set in order to set this
  property. Changing this updates the `enable_snat` of the router.

* `tenant_id` - (Optional) The owner of the floating IP. Required if admin wants
  to create a router for another tenant. Changing this creates a new router.

* `value_specs` - (Optional) Map of additional driver-specific options.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the router.

* `name` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `external_gateway` - See Argument Reference above.

* `enable_snat` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `value_specs` - See Argument Reference above.
