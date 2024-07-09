---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_subnet_ids_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-subnet-ids-v1"
description: |-
Get a list of VPC subnet ids from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC subnet you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/subnet/querying_subnets.html#vpc-subnet01-0003)

# opentelekomcloud_vpc_subnet_ids_v1

Use this data source to get a list of subnet ids for a vpc_id

This resource can be useful for getting back a list of subnet ids for a VPC.

## Example Usage

The following example shows outputting all cidr blocks for every subnet id in a VPC.

```hcl
data "opentelekomcloud_vpc_subnet_ids_v1" "subnet_ids" {
  vpc_id = var.vpc_id
}

data "opentelekomcloud_vpc_subnet_v1" "subnet" {
  for_each = data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids.ids
  id       = each.value
}

output "subnet_cidr_blocks" {
  value = [for s in data.opentelekomcloud_vpc_subnet_v1.subnet : s.cidr]
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) Specifies the VPC ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of all the subnet ids found. This data source will fail if none are found.
