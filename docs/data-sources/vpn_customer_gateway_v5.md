---
subcategory: "Virtual Private Network (VPN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_enterprise_vpn_customer_gateway_v5"
sidebar_current: "docs-opentelekomcloud-resource-enterprise-vpn-customer-gateway-v5"
description: |-
  Get details about a specific Enterprise VPN Customer Gateway resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for EVPN you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-network/api-ref/api_reference_enterprise_edition_vpn/apis_of_enterprise_edition_vpn/customer_gateway/index.html)

# opentelekomcloud_enterprise_vpn_customer_gateway_v5

Get details about a specific VPN customer gateway resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "id" {}

data "opentelekomcloud_enterprise_vpn_customer_gateway_v5" "test" {
  id = var.id
}
```


## Argument Reference

The following arguments are supported:

* `id` - (Required, String) Specifies the unique ID of the Enterprise VPN customer gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Indicates the customer gateway name.

* `id_value` - Indicates the identifier of a customer gateway. It could be an IP or an FQDN of VPN gateway depending on `id_type`.

* `id_type` - Indicates the identifier type of customer gateway.

* `asn` - Indicates the BGP ASN number of the customer gateway.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.

* `route_mode` - Indicates the routing mode.

* `ip` - Indicates the IP address of the customer gateway.

* `region` - Indicates the region in which resource is created.
