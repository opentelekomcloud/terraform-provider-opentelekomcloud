---
subcategory: "Virtual Private Cloud (VPC)"
---

Up-to-date reference of API arguments for VPC route you can get at
`https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/vpc_route`.

# opentelekomcloud_vpc_route_v2

Provides a resource to create a route within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_vpc_route_v2" "vpc_route" {
  type        = "peering"
  nexthop     = var.nexthop
  destination = "192.168.0.0/16"
  vpc_id      = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `destination` - (Required) Specifies the destination IP address or CIDR block. Changing this creates a new Route.

* `nexthop` - (Required) Specifies the next hop. If the route type is peering, enter the VPC peering connection ID. Changing this creates a new Route.

* `type` - (Required) Specifies the route type. Currently, the value can only be `peering`. Changing this creates a new Route.

* `vpc_id` - (Required) Specifies the VPC for which a route is to be added. Changing this creates a new Route.

* `tenant_id` - (Optional) Specifies the tenant ID. Only the administrator can specify the tenant ID of other tenant. Changing this creates a new Route.

## Attributes Reference

The following attributes are exported:

* `id` - The route ID.

* `destination` - The destination address in the CIDR notation format, for example, `192.168.200.0/24`.

* `nexthop` - The next hop. If the route type is `peering`, enter the VPC peering connection ID.

* `type` - The route type. Currently, the value can only be peering.

* `vpc_id` - The VPC ID of the route.

* `tenant_id` - The project ID.

## Import

VPC route can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_route_v2.vpc_route 2c7fs9f3-712b-18d1-940c-b50384177ee1
```
