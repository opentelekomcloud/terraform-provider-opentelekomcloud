---
subcategory: "Virtual Private Network (VPN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_enterprise_vpn_gateway_v5"
sidebar_current: "docs-opentelekomcloud-resource-enterprise-vpn-gateway-v5"
description: |-
  Get details about a specific Enterprise VPN Gateway Service resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for EVPN you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-network/api-ref/api_reference_enterprise_edition_vpn/apis_of_enterprise_edition_vpn/vpn_gateway/index.html)


# opentelekomcloud_enterprise_vpn_gateway_v5

Use this data source to get details about a specific Enterprise VPN gateway resource within OpenTelekomCloud.

## Example Usage

### Basic Usage

```hcl
variable "id" {}

data "opentelekomcloud_enterprise_vpn_gateway_v5" "gw_1" {
  id = var.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) Specifies the unique ID of the Enterprise VPN gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Indicates the name of the VPN gateway.

* `availability_zones` - Indicates the list of availability zone IDs.

* `flavor` - Indicates the flavor of the VPN gateway.

* `attachment_type` - Indicates the attachment type. 

* `network_type` - Indicates the network type. 

* `vpc_id` - Indicates the ID of the VPC to which the VPN gateway is connected.

* `local_subnets` - Indicates the list of local subnets.

* `connect_subnet` - Indicates the Network ID of the VPC subnet used by the VPN gateway.

* `er_id` - Indicates the enterprise router ID to attach with to VPN gateway.

* `ha_mode` - Indicates the HA mode of VPN gateway.

* `access_vpc_id` - Indicates the access VPC ID.

* `access_subnet_id` - Indicates the access subnet ID.

* `access_private_ip_1` - Indicates the private IP 1 in private network type VPN gateway.

* `access_private_ip_2` - Indicates the private IP 2 in private network type VPN gateway.

* `asn` - Indicates the ASN number of BGP.

* `status` - Indicates the status of VPN gateway.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.

* `used_connection_group` - Indicates the number of used connection groups.

* `used_connection_number` - Indicates the number of used connections.

* `region` - Indicates the region in which resource is created.

* `eip1` - Indicates the master 1 IP in active-active VPN gateway or the master IP in active-standby VPN gateway.
  The [object](#GatewayGetResponseEip) structure is documented below.

* `eip2` - Indicates the master 2 IP in active-active VPN gateway or the slave IP in active-standby VPN gateway.
  The [object](#GatewayGetResponseEip) structure is documented below.

<a name="GatewayGetResponseEip"></a>
The `eip1` or `eip2` block supports:

* `id` - Indicates the public IP ID.

* `bandwidth_name` - Indicates the bandwidth name.

* `bandwidth_size` - Indicates the Bandwidth size in Mbit/s. 

* `charge_mode` - Indicates the charge mode of the bandwidth.

* `type` - Indicates the EIP type.

* `bandwidth_id` - The bandwidth ID.

* `ip_address` - The public IP address.

* `ip_version` - Specifies the EIP version.
