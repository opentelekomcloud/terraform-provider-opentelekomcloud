---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_route_table_route_v1"
sidebar_current: "docs-opentelekomcloud-resource-vpc-route-table-route-v1"
description: |-
  Manages a VPC Route Table Route resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC route table you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/route_table/index.html)

# opentelekomcloud_vpc_route_table_route_v1

Manages an individual route within a VPC route table. This resource allows managing
routes independently of the route table lifecycle, including adding routes to the
default route table.

## Example Usage

### Route on Default Route Table

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc-2"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "my_peering"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_table_route_v1" "route" {
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  destination = "172.16.0.0/16"
  type        = "peering"
  nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
  description = "peering route"
}
```

### Route on Custom Route Table

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc-2"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "my_peering"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table" {
  name   = "my_table"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_route_table_route_v1" "route" {
  vpc_id         = opentelekomcloud_vpc_v1.vpc_1.id
  route_table_id = opentelekomcloud_vpc_route_table_v1.table.id
  destination    = "172.16.0.0/16"
  type           = "peering"
  nexthop        = opentelekomcloud_vpc_peering_connection_v2.peering.id
  description    = "peering route on custom table"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the route.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID that the route table belongs to.
  Changing this creates a new resource.

* `destination` - (Required, String, ForceNew) Specifies the destination address in the CIDR notation format,
  for example, 192.168.200.0/24. The destination of each route must be unique and cannot overlap
  with any subnet in the VPC. Changing this creates a new resource.

* `type` - (Required, String) Specifies the route type. Currently, the value can be:
  **ecs**, **eni**, **vip**, **nat**, **peering**, **vpn**, **dc**, **egw**, **er**, **subeni** and **local**

* `nexthop` - (Required, String) Specifies the next hop.
  + If the route type is **ecs**, the value is an ECS instance ID in the VPC.
  + If the route type is **eni**, the value is the extension NIC of an ECS in the VPC.
  + If the route type is **vip**, the value is a virtual IP address.
  + If the route type is **nat**, the value is a NAT gateway ID.
  + If the route type is **peering**, the value is a VPC peering connection ID.
  + If the route type is **vpn**, the value is a VPN gateway ID.
  + If the route type is **dc**, the value is a Direct Connect gateway ID.
  + If the route type is **egw**, the value is a VPC endpoint ID.
  + If the route type is **er**, the value is the ID of an enterprise router.
  + If the route type is **subeni**, the value is the ID of a supplementary network interface.

* `description` - (Optional, String) Specifies the supplementary information about the route.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).

* `route_table_id` - (Optional, String, ForceNew) Specifies the route table ID. If omitted, the
  default route table of the VPC will be used. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `{route_table_id}/{destination}`.
* `route_table_name` - The name of the route table.

## Import

Routes can be imported using the route table ID and destination, separated by a slash, e.g.

-> **NOTE:** The import ID contains the route table UUID followed by a `/` and the CIDR destination
(which itself contains a `/`), e.g. `<route_table_id>/<cidr_destination>`.

```
$ terraform import opentelekomcloud_vpc_route_table_route_v1.route 14c6491a-f90a-41aa-a206-f58bbacdb47d/172.16.0.0/16
```
