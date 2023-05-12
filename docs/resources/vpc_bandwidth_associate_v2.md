---
subcategory: "Virtual Private Cloud (VPC)"
---

Up-to-date reference of API arguments for VPC bandwidth association you can get at
`https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/bandwidth_v2.0`.

# opentelekomcloud_vpc_bandwidth_associate_v2

Provides a resource to associate floating IP with a shared bandwidth within Open Telekom Cloud.

## Example Usage

```hcl
resource "opentelekomcloud_networking_floatingip_v2" "ip1" {}
resource "opentelekomcloud_networking_floatingip_v2" "ip2" {}

resource "opentelekomcloud_vpc_bandwidth_v2" "band20m" {
  name = "bandwidth-20MBit"
  size = 20
}

resource "opentelekomcloud_vpc_bandwidth_associate_v2" "associate" {
  bandwidth = opentelekomcloud_vpc_bandwidth_v2.band20m.id
  floating_ips = [
    opentelekomcloud_networking_floatingip_v2.ip1.id,
    opentelekomcloud_networking_floatingip_v2.ip2.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) Specifies ID of the bandwidth to be assigned.

* `floating_ips` - (Required) Specifies IDs of floating IPs to be added to the bandwidth.

->
After an EIP is removed from a shared bandwidth, a dedicated bandwidth will be allocated to the EIP, and you will be
billed for the dedicated bandwidth.

* `backup_charge_mode` - (Optional) Specifies whether the dedicated bandwidth used by the EIP that has been removed from
  a shared bandwidth is billed by traffic or by bandwidth.

  The value can be `bandwidth` or `traffic`.

  Default value is `bandwidth`.

* `backup_size` - (Optional) Specifies the size (Mbit/s) of the dedicated bandwidth used by the EIP that has been
  removed from a shared bandwidth.

  Default value is `1`.

## Import

VPC bandwidth association can be imported using the bandwidth `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_bandwidth_associate_v2.associate eb187fc8-e482-43eb-a18a-9da947ef89f6
```
