---
subcategory: "NAT"
---

# opentelekomcloud_nat_snat_rule_v2

Manages a V2 snat rule resource within OpenTelekomCloud Nat.

## Example Usage

```hcl
variable "network_id" {}
variable "vpc_id" {}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name                = "nat_1"
  description         = "test for terraform"
  spec                = "1"
  internal_network_id = var.network_id
  router_id           = var.vpc_id
}

resource "opentelekomcloud_nat_snat_rule_v2" "snat_1" {
  nat_gateway_id = opentelekomcloud_nat_gateway_v2.nat_1.id
  floating_ip_id = opentelekomcloud_networking_floatingip_v2.fip_1.id
  cidr           = "192.168.0.0/24"
  source_type    = 0
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) ID of the nat gateway this snat rule belongs to.
  Changing this creates a new snat rule.

* `network_id` - (Optional) ID of the network this snat rule connects to. This parameter
  and `cidr` are alternative. Changing this creates a new snat rule.

* `source_type` - (Optional) `0`: Either `network_id` or cidr can be specified in a VPC. `1`:
  Only `cidr` can be specified over a dedicated network. Changing this creates a new snat rule.

* `cidr` - (Optional) Specifies CIDR, which can be in the format of a network segment or
  a host IP address. This parameter and `network_id` are alternative. If the value of
  `source_type` is `0`, the CIDR block must be a subset of the VPC subnet CIDR block. If
  the value of `source_type` is `1`, the CIDR block must be a CIDR block of Direct Connect
  and cannot conflict with the VPC CIDR blocks. Changing this creates a new snat rule.

* `floating_ip_id` - (Required) ID of the floating ip this snat rule connects to.
  Changing this creates a new snat rule.

## Attributes Reference

The following attributes are exported:

* `nat_gateway_id` - See Argument Reference above.

* `network_id` - See Argument Reference above.

* `floating_ip_id` - See Argument Reference above.

* `source_type` - See Argument Reference above.

* `cidr` - See Argument Reference above.
