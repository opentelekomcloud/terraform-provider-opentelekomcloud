---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_route_table_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-route-table-v1"
description: |-
Get details about a specific VPC route table from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC route table you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/route_table/querying_route_tables.html#vpc-apiroutetab-0001)

# opentelekomcloud_vpc_route_table_v1

Provides details about a specific VPC route table.

## Example Usage

```hcl
variable "vpc_id" {}

# get the default route table
data "opentelekomcloud_vpc_route_table_v1" "default" {
  vpc_id = var.vpc_id
}

# get a custom route table
data "opentelekomcloud_vpc_route_table_v1" "custom" {
  vpc_id = var.vpc_id
  name   = "my"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String) Specifies the VPC ID where the route table resides.

* `name` - (Optional, String) Specifies the name of the route table.

* `id` - (Optional, String) Specifies the ID of the route table.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `default` - Whether the route table is default or not.

* `description` - The supplementary information about the route table.

* `subnets` - An array of one or more subnets associating with the route table.

* `region` - The region in which belongs the vpc route table.

* `route` - The route object list. The [route object](#route_object) is documented below.

<a name="route_object"></a>
The `route` block supports:

* `type` - The route type.
* `destination` - The destination address in the CIDR notation format
* `nexthop` - The next hop.
* `description` - The description about the route.
