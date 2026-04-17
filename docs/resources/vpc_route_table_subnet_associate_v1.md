---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_route_table_subnet_associate_v1"
sidebar_current: "docs-opentelekomcloud-resource-vpc-route-table-subnet-associate-v1"
description: |-
  Manages a VPC Route Table Subnet Association resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC route table you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/route_table/index.html)

# opentelekomcloud_vpc_route_table_subnet_associate_v1

Manages the association between a subnet and a VPC route table.

A subnet is always associated with exactly one route table. Associating a subnet with a
new route table automatically moves it from the previous one. Destroying this resource
disassociates the subnet, which returns it to the VPC's default route table.

~> **NOTE:** When using this resource, the `opentelekomcloud_vpc_route_table_v1` resource must include
`lifecycle { ignore_changes = [subnets] }` to avoid conflicts. Both resources manage subnet associations
on the same route table, and without `ignore_changes`, Terraform will detect a perpetual diff on the
route table's `subnets` attribute.

## Example Usage

```hcl
resource "opentelekomcloud_vpc_v1" "vpc" {
  name = "my-vpc"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  name       = "my-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
}

resource "opentelekomcloud_vpc_route_table_v1" "table" {
  name   = "my-table"
  vpc_id = opentelekomcloud_vpc_v1.vpc.id

  lifecycle {
    ignore_changes = [subnets]
  }
}

resource "opentelekomcloud_vpc_route_table_subnet_associate_v1" "assoc" {
  route_table_id = opentelekomcloud_vpc_route_table_v1.table.id
  subnet_id      = opentelekomcloud_vpc_subnet_v1.subnet.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the association.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `route_table_id` - (Required, String, ForceNew) Specifies the route table ID to associate the
  subnet with. Changing this creates a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID to associate with the route
  table. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `{route_table_id}/{subnet_id}`.
* `vpc_id` - The VPC ID that the route table belongs to.

## Import

Route table subnet associations can be imported using the route table ID and subnet ID, separated by a slash, e.g.

```
$ terraform import opentelekomcloud_vpc_route_table_subnet_associate_v1.assoc 14c6491a-f90a-41aa-a206-f58bbacdb47d/a1b2c3d4-e5f6-7890-abcd-ef1234567890
```
