---
subcategory: "Virtual Private Cloud (VPC)"
---

#

[API Reference](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/route_table/querying_route_tables.html)

# opentelekomcloud_vpc_route_tables_v1

Returns list of route tables.

## Example Usage

```hcl
variable "vpc_id" {}
variable "route_table_id" {}
variable "subnet_network_id" {}

# get all route tables
data "opentelekomcloud_vpc_route_tables_v1" "all_route_tables" {
}

# get route tables for specific vpc
data "opentelekomcloud_vpc_route_tables_v1" "vpc_route_tables" {
  vpc_id = var.vpc_id
}

# get a list that includes single specific route table
data "opentelekomcloud_vpc_route_table_v1" "single_route_table" {
  id = var.route_table_id
}

# get a list of route table associated with a specific subnet
data "opentelekomcloud_vpc_route_table_v1" "subnet_route_table" {
  subnet_id = var.subnet_network_id
}
```

## Argument Reference

The following arguments are supported:

- `vpc_id` - (Optional, String) Specifies the VPC ID where the route tables reside.

- `subnet_id` - (Optional, String) Specifies the id of the subnet. **Note**: the corresponding subnet resource attribute is `network_id`.

- `id` - (Optional, String) Specifies the ID of the route table.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `routetables` - list of [routetable object](#routtable_object) (documented belos)

<a name="routetable_object"></a>
The `routetable` object has the following attributes:

- `default` - Whether the route table is default or not.

- `description` - Route table description.

- `subnets` - An array of subnets associating with the route table.

- `tenant_id` - Project id to which route table belongs.

- `vpc_id` - VPC Id to which route table belongs.

- `routes` - List of non-local routes in the route table (`local` routes are considered system internal and can't be managed via API, though are visible in web UI). Structure of the [route object](#route_object) is documented below.

<a name="route_object"></a>
The `route` object has the following attributes:

- `type` - The route type. Check [API reference](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/route_table/creating_a_route_table.html) for supperted types.

- `destination` - The destination address in the CIDR notation format

- `nexthop` - The next hop. Value depends on the route type.

- `description` - Route description.
