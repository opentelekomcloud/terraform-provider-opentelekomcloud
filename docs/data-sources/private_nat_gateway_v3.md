---
subcategory: "Private NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_private_nat_gateway_v3"
sidebar_current: "docs-opentelekomcloud-data-source-private-nat-gateway-v3"
description: |-
  Manages a Private NAT Gateway data source within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Private NAT gateway you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/apis_for_private_nat_gateways_v3.0/private_nat_gateways/index.html)

# opentelekomcloud_private_nat_gateway_v3

Manages a V3 Private NAT Gateway data source within OpenTelekomCloud.

## Example Usage

### List all Private NAT gateways

```hcl
data "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {}
```

### Get Private NAT gateway using ID

```hcl
variable gateway_id {}

data "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  id = var.gateway_id
}
```

### Get Private NAT gateway using Name

```hcl
variable gateway_name {}

data "opentelekomcloud_private_nat_gateway_v3" "gateway_1" {
  name = var.gateway_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the private NAT gateway name.

* `id` - (Optional, String) Specifies the private NAT gateway name.

* `spec` - (Optional, String) Specifies the private NAT gateway specifications. The value can be: `Small`, `Medium`, `Large`, `Extra-large`. Default value: `Small`.

* `downlink_vpcs` - (Required, List, ForceNew) Specifies the VPC where the private NAT gateway works. The [downlink_vpcs](#downlink_vpcs) structure is documented below. Default value: `0`.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project that is associated with the private NAT gateway when the private NAT gateway is created.

* `tags` - (Optional, Map, ForceNew) Specifies the tag list in key/value format.

<a name="downlink_vpcs"></a>
The `downlink_vpcs` block supports:

* `virsubnet_id` - (Required, String, ForceNew) Specifies the ID of the enterprise project that is associated with the private NAT gateway when the private NAT gateway is created.

* `ngport_ip_address` - (Optional, String, ForceNew) Specifies the private IP address of the private NAT gateway.


## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `gateways` - Indicates the private NAT gateways. The structure is defined below.

The `gateways` block supports: 

* `id` - Private NAT gateway ID.

* `name` - Indicates the private NAT gateway name.

* `description` - Provides supplementary information about the private NAT gateway.

* `spec` - Indicates the private NAT gateway specifications. The value can be: `Small`, `Medium`, `Large`, `Extra-large`.

* `downlink_vpcs` - Indicates the VPC where the private NAT gateway works. The structure is documented below.
    * `virsubnet_id` - Indicates the ID of the subnet where the private NAT gateway works.
    * `ngport_ip_address` - Indicates the private IP address of the private NAT gateway.
    * `vpc_id` - Indicates the ID of the VPC where the private NAT gateway works.

* `enterprise_project_id` - Indicates the ID of the enterprise project that is associated with the private NAT gateway when the private NAT gateway is created.

* `tags` - Indicates the tag list in key/value format.

* `project_id` - Indicates the project ID.

* `status` - Indicates the private NAT gateway status. The value can be: `ACTIVE` (The private NAT gateway is running properly) or `FROZEN` (The private NAT gateway is frozen).

* `created_at` - Indicates the time when the private NAT gateway was created. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `updated_at` - Indicates the time when the private NAT gateway was updated. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `rule_max` - Indicates Specifies the maximum number of rules. Value range: `0-65535`

* `transit_ip_pool_size_max` - Specifies the maximum number of transit IP addresses in a transit IP address pool. Value range: `0-100`
