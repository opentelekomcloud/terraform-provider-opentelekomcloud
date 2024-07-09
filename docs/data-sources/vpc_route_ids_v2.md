---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_route_ids_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-route-ids-v2"
description: |-
Get a list of route ids from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC route you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/vpc_route/querying_vpc_routes.html#vpc-route-0001)

# opentelekomcloud_vpc_route_ids_v2

Use this data source to get a list of route ids for a vpc_id.

This resource can be useful for getting back a list of route ids for a vpc.

## Example Usage

```hcl
variable "vpc_id" {}

data "opentelekomcloud_vpc_route_ids_v2" "example" {
  vpc_id = var.vpc_id
}

data "opentelekomcloud_vpc_route_v2" "vpc_route" {
  for_each = data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids.ids
  id       = each.value
}

output "route_nexthop" {
  value = [for hop in data.opentelekomcloud_vpc_subnet_v1.subnet : hop.cidr]
}
```

## Argument Reference

* `vpc_id` - (Required) The VPC ID that you want to filter from.

## Attributes Reference

* `ids` - A list of all the route ids found. This data source will fail if none are found.
