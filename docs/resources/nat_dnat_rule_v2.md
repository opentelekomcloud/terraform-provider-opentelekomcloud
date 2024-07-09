---
subcategory: "NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_nat_dnat_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-nat-dnat-rule-v2"
description: |-
Manages a NAT DNAT Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for NAT DNAT you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/api_v2.0/dnat_rules)

# opentelekomcloud_nat_dnat_rule_v2

Manages a V2 DNAT rule resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "nat_gw_id" {}
variable "floating_ip_id" {}
variable "private_ip" {}

resource "opentelekomcloud_nat_dnat_rule_v2" "dnat_1" {
  floating_ip_id        = var.floating_ip_id
  nat_gateway_id        = var.nat_gw_id
  private_ip            = var.private_id
  internal_service_port = 993
  protocol              = "tcp"
  external_service_port = 242
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip_id` - (Required) Specifies the ID of the floating IP address.
  Changing this creates a new resource.

* `internal_service_port` - (Required) Specifies port used by ECSs or BMSs
  to provide services for external systems. Changing this creates a new resource.

* `nat_gateway_id` - (Required) ID of the NAT gateway this DNAT rule belongs to.
   Changing this creates a new DNAT rule.

-> You can create a DNAT rule only when status of the NAT gateway is set to `ACTIVE`
and `admin_state_up` of the NAT gateway administrator to `True`.

* `port_id` - (Optional) Specifies the port ID of an ECS or a BMS.
  This parameter and `private_ip` are alternative. Changing this creates a
  new DNAT rule.

->
When the DNAT rule is used in the **VPC** scenario, use `port_id` parameter.

* `private_ip` - (Optional) Specifies the private IP address of a
  user, for example, the IP address of a VPC for dedicated connection.
  This parameter and `port_id` are alternative. Changing this creates a new DNAT rule.

->
When the DNAT rule is used in the **Direct Connect** scenario, use `private_ip` parameter.

* `protocol` - (Required) Specifies the protocol type. Currently,
  `tcp`, `udp`, and `any` are supported. Changing this creates a new DNAT rule.

-> If you create a rule that applies to all port types, set `internal_service_port` to `0`,
`external_service_port` to `0`, and `protocol` to `any`.

* `external_service_port` - (Required) Specifies port used by ECSs or
  BMSs to provide services for external systems. Changing this creates a new DNAT rule.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created_at` - DNAT rule creation time.

* `status` - DNAT rule status.

* `floating_ip_address` - The actual floating IP address.

## Import

DNAT can be imported using the following format:

```sh
terraform import opentelekomcloud_nat_dnat_rule_v2.dnat_1 f4f783a7-b908-4215-b018-724960e5df4a
```
