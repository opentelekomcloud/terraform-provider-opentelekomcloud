---
subcategory: "NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_nat_gateway_v2"
sidebar_current: "docs-opentelekomcloud-datasource-nat-gateway-v2"
description: |-
  Get details about NAT Gateway resource from OpenTelekomCloud
---

Up-to-date reference of API arguments for NAT Gateway you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/api_v2.0/nat_gateway_service/querying_nat_gateways.html#nat-api-0002)

# opentelekomcloud_nat_gateway_v2

Use this data source to get the info about an existing V2 NAT Gateway resource within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_nat_gateway_v2" "this" {
  name = "tf_nat"
  spec = "1"
}
```

## Argument Reference

The following arguments are supported:

* `nat_id` - (Optional) The ID of the NAT Gateway.

* `name` - (Optional) The name of the NAT Gateway.

* `description` - (Optional) The description of the NAT Gateway.

* `spec` - (Optional) The specification of the NAT Gateway, valid values are `"1"`, `"2"`, `"3"`, `"4"`.

* `tenant_id` - (Optional) The target tenant ID in which to allocate the NAT
  Gateway.

* `router_id` - (Optional) ID of the router (or VPC) this NAT Gateway belongs to.

* `internal_network_id` - (Optional) ID of the network this NAT Gateway connects to.

* `status` - (Optional) Specifies the NAT gateway status.

* `admin_state_up` - (Optional) Specifies whether the NAT gateway is up or down. Possible values are:
  * `true` refers to NAT gateway being up.
  * `false` refers to NAT gateway being down.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - ID of NAT gateway.

* `region` - Region of NAT gateway.
