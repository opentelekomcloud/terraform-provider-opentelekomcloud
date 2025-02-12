---
subcategory: "Virtual Private Network (VPN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_enterprise_vpn_connection_v5"
sidebar_current: "docs-opentelekomcloud-resource-enterprise-vpn-connection-v5"
description: |-
  Get details about a specific Enterprise VPN connection Service resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for EVPN connection you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-network/api-ref/api_reference_enterprise_edition_vpn/apis_of_enterprise_edition_vpn/vpn_connection/index.html)

# opentelekomcloud_enterprise_vpn_connection_v5

Get details about a specific VPN connection resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "id" {}

data "opentelekomcloud_enterprise_vpn_connection_v5" "conn" {
  id = var.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) Specifies the unique ID of the Enterprise VPN Connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the VPN connection.

* `gateway_id` - The VPN gateway ID.

* `gateway_ip` - The VPN gateway IP ID.

* `vpn_type` - The connection type.

* `customer_gateway_id` - The customer gateway ID.

* `peer_subnets` - The CIDR list of customer subnets. 

* `tunnel_local_address` - The local tunnel address.

* `tunnel_peer_address` - The peer tunnel address.

* `enable_nqa` -  Whether to enable NQA check.

* `ikepolicy` - The IKE policy configurations.
The [ikepolicy](#Connection_CreateRequestIkePolicy) structure is documented below.

* `ipsecpolicy` - The IPsec policy configurations.
The [ipsecpolicy](#Connection_CreateRequestIpsecPolicy) structure is documented below.

* `policy_rules` - The policy rules.
The [policy_rules](#Connection_PolicyRule) structure is documented below.

* `tags` - Specifies the tags of the VPN connection.

* `ha_role` - Specifies the mode of the VPN connection.

* `status` - The status of the VPN connection.

* `created_at` - The create time.

* `updated_at` - The update time.

* `region` - Specifies the region in which resource is created.

<a name="Connection_CreateRequestIkePolicy"></a>
The `ikepolicy` block supports:

* `authentication_algorithm` - The authentication algorithm. 

* `encryption_algorithm` - The encryption algorithm.

* `ike_version` - The IKE negotiation version.

* `lifetime_seconds` - The life cycle of SA in seconds.

* `local_id_type` - The local ID type.

* `local_id` - The local ID.

* `peer_id_type` - The peer ID type.

* `peer_id` - The peer ID.

* `phase_one_negotiation_mode` - The negotiation mode.

* `authentication_method` - The authentication method during IKE negotiation.

* `dh_group` - Specifies the DH group used for key exchange in phase 1.

* `dpd` -  Specifies the dead peer detection (DPD) object.
  The [dpd](#Connection_DPD) structure is documented below.

<a name="Connection_DPD"></a>
The `dpd` block supports:

* `timeout` - Specifies the interval for retransmitting DPD packets.

* `interval` - Specifies the DPD idle timeout period.

* `msg` - Specifies the format of DPD packets.

<a name="Connection_CreateRequestIpsecPolicy"></a>
The `ipsecpolicy` block supports:

* `authentication_algorithm` - The authentication algorithm.

* `encryption_algorithm` - The encryption algorithm.

* `pfs` - The DH key group used by PFS.

* `lifetime_seconds` - The lifecycle time of Ipsec tunnel in seconds.

* `transform_protocol` - The transform protocol.

* `encapsulation_mode` - The encapsulation mode.

<a name="Connection_PolicyRule"></a>
The `policy_rules` block supports:

* `rule_index` - The rule index.

* `destination` -  The list of destination CIDRs.

* `source` - The source CIDR.
