---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_nat_snat_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-nat-snat-rule-v2"
description: |-
  Manages a V2 snat rule resource within OpenTelekomCloud Nat.
---

# opentelekomcloud\_nat\_snat\_rule_v2

Manages a V2 snat rule resource within OpenTelekomCloud Nat

## Example Usage

```hcl
resource "opentelekomcloud_nat_snat_rule_v2" "snat_1" {
  nat_gateway_id = "3c0dffda-7c76-452b-9dcc-5bce7ae56b17"
  network_id = "dc8632e2-d9ff-41b1-aa0c-d455557314a0"
  floating_ip_id = "0a166fc5-a904-42fb-b1ef-cf18afeeddca"
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) ID of the nat gateway this snat rule belongs to.
    Changing this creates a new snat rule.

* `network_id` - (Optional) ID of the network this snat rule connects to. This parameter
	and cidr are alternative. Changing this creates a new snat rule.

* `source_type` - (Optional) 0: Either network_id or cidr can be specified in a VPC. 1:
	Only cidr can be specified over a dedicated network. Changing this creates a new snat rule.

* `cidr` - (Optional) Specifies CIDR, which can be in the format of a network segment or
	a host IP address. This parameter and network_id are alternative. If the value of
	source_type is 0, the CIDR block must be a subset of the VPC subnet CIDR block. If
	the value of source_type is 1, the CIDR block must be a CIDR block of Direct Connect
	and cannot conflict with the VPC CIDR blocks. Changing this creates a new snat rule.

* `floating_ip_id` - (Required) ID of the floating ip this snat rule connets to.
    Changing this creates a new snat rule.

## Attributes Reference

The following attributes are exported:

* `nat_gateway_id` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `floating_ip_id` - See Argument Reference above.
* `source_type` - See Argument Reference above.
* `cidr` - See Argument Reference above.
