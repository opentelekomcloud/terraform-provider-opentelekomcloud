---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekcomcloud_subnet_ids_v1"
sidebar_current: "docs-opentelekomcloud-datasource-subnet-ids-v1"
description: |-
  Provides a list of subnet Ids for a VPC
---

# Data Source: opentelekomcloud_vpc_subnet_ids_v1

`opentelekomcloud_vpc_subnet_ids_v1` provides a list of subnet ids for a vpc_id

This resource can be useful for getting back a list of subnet ids for a vpc.

## Example Usage

The following example shows outputing all cidr blocks for every subnet id in a vpc.

 ```hcl
data "opentelekomcloud_vpc_subnet_ids_v1" "subnet_ids" {
  vpc_id = "${var.vpc_id}" 
}

data "opentelekomcloud_vpc_subnet_v1" "subnet" {
  count = "${length(data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids.ids)}"
  id    = "${data.opentelekomcloud_vpc_subnet_ids_v1.subnet_ids.ids[count.index]}"
 }

output "subnet_cidr_blocks" {
  value = "${data.opentelekomcloud_vpc_subnet_v1.subnet.*.cidr}"
}
 ```

## Argument Reference

The following arguments are supported:

* `vpc_id` (Required) - Specifies the VPC ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of all the subnet ids found. This data source will fail if none are found.