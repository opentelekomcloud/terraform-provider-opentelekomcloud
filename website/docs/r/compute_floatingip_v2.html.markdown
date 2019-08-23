---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_floatingip_v2"
sidebar_current: "docs-opentelekomcloud-resource-compute-floatingip-v2"
description: |-
  Manages a V2 floating IP resource within OpenTelekomCloud Nova (compute).
---

# opentelekomcloud\_compute\_floatingip_v2

Manages a V2 floating IP resource within OpenTelekomCloud Nova (compute)
that can be used for compute instances.
These are similar to Neutron (networking) floating IP resources,
but only networking floating IPs can be used with load balancers.

Floating IPs created with this module will have a bandwidth of 1000Mbit/s,
for manually specifying the bandwidth please use the
[`opentelekomcloud_vpc_eip_v1`](vpc_eip_v1.html) module.


## Example Usage

```hcl
resource "opentelekomcloud_compute_floatingip_v2" "floatip_1" {
}
```

## Argument Reference

The following arguments are supported:

* `pool` - (Optional) The name of the pool from which to obtain the floating
    IP. Default value is admin_external_net. Changing this creates a new floating IP.

## Attributes Reference

The following attributes are exported:

* `pool` - See Argument Reference above.
* `address` - The actual floating IP address itself.
* `fixed_ip` - The fixed IP address corresponding to the floating IP.
* `instance_id` - UUID of the compute instance associated with the floating IP.

## Import

Floating IPs can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_compute_floatingip_v2.floatip_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
