---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_route_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-route-v2"
description: |-
Get details about a specific VPC route from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC route you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/vpc_route/querying_vpc_routes.html#vpc-route-0001)

# opentelekomcloud_vpc_route_v2

Use this data source to get details about a specific VPC route.

## Example Usage

```hcl
variable "route_id" {}

data "opentelekomcloud_vpc_route_v2" "vpc_route" {
  id = var.route_id
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_v1" {
  name       = "test-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = data.opentelekomcloud_vpc_route_v2.vpc_route.vpc_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
routes in the current tenant. The given filters must match exactly one
route whose data will be exported as attributes.

* `id` - (Optional) The id of the specific route to retrieve.

* `vpc_id` - (Optional) The id of the VPC that the desired route belongs to.

* `destination` - (Optional) The route destination address (CIDR).

* `tenant_id` - (Optional) Only the administrator can specify the tenant ID of other tenants.

* `type` - (Optional) Route type for filtering.

## Attribute Reference

All of the argument attributes are also exported as result attributes.

* `nexthop` - The next hop of the route. If the route type is peering, it will provide VPC peering connection ID.
