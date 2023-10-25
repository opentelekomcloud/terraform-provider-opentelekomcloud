---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_vpc_route_table_v1

Provides a resource to create a route table within OpenTelekomCloud.

## Example Usage

### Basic Custom Route Table

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1-1" {
  name       = "vpc-1-subnet-1-1"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1-2" {
  name       = "vpc-1-subnet-1-2"
  cidr       = "192.168.10.0/24"
  gateway_ip = "192.168.10.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc-2"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_2-1" {
  name       = "vpc-2-subnet-2-1"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "my_peering"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table_1" {
  name        = "my_route"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  description = "created by terraform with routes"

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering rule"
  }
}
```

### Associating Subnets with a Route Table

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1-1" {
  name       = "vpc-1-subnet-1-1"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1-2" {
  name       = "vpc-1-subnet-1-2"
  cidr       = "192.168.10.0/24"
  gateway_ip = "192.168.10.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_1.id
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc-2"
  cidr = "172.16.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_2-1" {
  name       = "vpc-2-subnet-2-1"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name        = "my_peering"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table_1" {
  name        = "my_table"
  vpc_id      = opentelekomcloud_vpc_v1.vpc_1.id
  description = "created by terraform with subnets"

  subnets     = [
    opentelekomcloud_vpc_subnet_v1.subnet_1-1.id,
    opentelekomcloud_vpc_subnet_v1.subnet_1-2.id,
  ]

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = opentelekomcloud_vpc_peering_connection_v2.peering.id
    description = "peering rule"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the vpc route table.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID for which a route table is to be added.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the route table name. The value is a string of no more than
  64 characters that can contain letters, digits, underscores (_), hyphens (-), and periods (.).

* `description` - (Optional, String) Specifies the supplementary information about the route table.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).

* `subnets` - (Optional, List) Specifies an array of one or more subnets associating with the route table.

  -> **NOTE:** The custom route table associated with a subnet affects only the outbound traffic.
  The default route table determines the inbound traffic.

* `route` - (Optional, List) Specifies the route object list. The [route object](#route_object)
  is documented below.

<a name="route_object"></a>
The `route` block supports:

* `destination` - (Required, String) Specifies the destination address in the CIDR notation format,
  for example, 192.168.200.0/24. The destination of each route must be unique and cannot overlap
  with any subnet in the VPC.

* `type` - (Required, String) Specifies the route type. Currently, the value can be:
  **ecs**, **eni**, **vip**, **nat**, **peering**, **vpn**, **dc** and **cc**.

* `nexthop` - (Required, String) Specifies the next hop.
  + If the route type is **ecs**, the value is an ECS instance ID in the VPC.
  + If the route type is **eni**, the value is the extension NIC of an ECS in the VPC.
  + If the route type is **vip**, the value is a virtual IP address.
  + If the route type is **nat**, the value is a VPN gateway ID.
  + If the route type is **peering**, the value is a VPC peering connection ID.
  + If the route type is **vpn**, the value is a VPN gateway ID.
  + If the route type is **dc**, the value is a Direct Connect gateway ID.
  + If the route type is **cc**, the value is a Cloud Connection ID.

* `description` - (Optional, String) Specifies the supplementary information about the route.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `created_at` - Specifies the time (UTC) when the route table is created.
* `updated_at` - Specifies the time (UTC) when the route table is updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

vpc route tables can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_vpc_route_table_v1.my_table 14c6491a-f90a-41aa-a206-f58bbacdb47d
```
