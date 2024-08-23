---
subcategory: "NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_nat_dnat_rules_v2"
sidebar_current: "docs-opentelekomcloud-datasource-nat-dnat-rules-v2"
description: |-
  Get details about NAT Gateway DNAT rules resource from OpenTelekomCloud
---

Up-to-date reference of API arguments for NAT Gateway you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/api_v2.0/dnat_rules/querying_dnat_rules.html#nat-api-0012)

# opentelekomcloud_nat_dnat_rules_v2

Use this data source to get the list of DNAT rules within OpenTelekomCloud..

## Example Usage

```hcl
variable "protocol" {}

data "opentelekomcloud_nat_dnat_rules_v2" "rule" {
  protocol = var.protocol
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Optional, String) Specifies the ID of the DNAT rule.

* `gateway_id` - (Optional, String) Specifies the ID of the NAT gateway to which the DNAT rule belongs.

* `protocol` - (Optional, String) Specifies the protocol type of the DNAT rule.
  The value can be one of the following:
  * `tcp`
  * `udp`
  * `any`

* `description` - (Optional, String) Specifies the description of the DNAT rule.

* `port_id` - (Optional, String) Specifies the port ID of the backend instance to which the DNAT rule belongs.

* `private_ip` - (Optional, String) Specifies the private IP address of the backend instance to which the DNAT rule
  belongs.

* `status` - (Optional, String) Specifies the status of the DNAT rule.
  The value can be one of the following:
  * `ACTIVE`: The SNAT rule is available.
  * `EIP_FREEZED`: The EIP is frozen associated with SNAT rule.
  * `INACTIVE`: The SNAT rule is unavailable.

* `internal_service_port` - (Optional, Int) Specifies the port of the backend instance to which the DNAT rule
  belongs.

* `external_service_port` - (Optional, Int) Specifies the port of the EIP associated with the DNAT rule.

* `floating_ip_id` - (Optional, String) Specifies the ID of the EIP associated with the DNAT rule.

* `floating_ip_address` - (Optional, String) Specifies the IP address of the EIP associated with the DNAT rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region where the DNAT rules are located.

* `rules` - The list ot the DNAT rules.
  The [rules](#nat_dnat_rules) structure is documented below.

<a name="nat_dnat_rules"></a>
The `rules` block supports:

* `id` - The ID of the DNAT rule.

* `gateway_id` - The ID of the NAT gateway to which the DNAT rule belongs.

* `protocol` - The protocol type of the DNAT rule.

* `port_id` - The port ID of the backend instance to which the DNAT rule belongs.

* `private_ip` - The private IP address of the backend instance to which the DNAT rule belongs.

* `internal_service_port` - The port of the backend instance to which the DNAT rule belongs.

* `external_service_port` - The port of the EIP associated with the DNAT rule belongs.

* `floating_ip_id` - The ID of the EIP associated with the DNAT rule.

* `floating_ip_address` - The IP address of the EIP associated with the DNAT rule.

* `description` - The description of the DNAT rule.

* `status` - The status of the DNAT rule.

* `created_at` - The creation time of the DNAT rule.
