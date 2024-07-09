---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-v1"
description: |-
Get details about a specific VPC from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/virtual_private_cloud/querying_vpcs.html#vpc-api01-0003)

# opentelekomcloud_vpc_v1

Use this data source to get details about a specific VPC.

This data source can prove useful when a module accepts a VPC id as an input variable and needs to, for example,
determine the CIDR block of that VPC.

## Example Usage

```hcl
variable "vpc_name" {}

data "opentelekomcloud_vpc_v1" "vpc" {
  name   = var.vpc_name
  shared = true
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPCs in the current region.
The given filters must match exactly one VPC whose data will be exported as attributes.

* `id` - (Optional) The id of the specific VPC to retrieve.

* `status` - (Optional) The current status of the desired VPC.
  Can be either `CREATING`, `OK`, `DOWN`, `PENDING_UPDATE`, `PENDING_DELETE`, or `ERROR`.

* `name` - (Optional) A unique name for the VPC. The name must be unique for a tenant.
  The value is a string of no more than 64 characters and can contain digits, letters, underscores (_), and hyphens (-).

* `cidr` - (Optional) The cidr block of the desired VPC.

* `shared` - (Optional) Enable SNAT (In order to let instances without an EIP access the internet).

## Attributes Reference

The following attributes are exported:

* `id` - ID of the VPC.

* `name` -  See Argument Reference above.

* `status` - See Argument Reference above.

* `cidr` - See Argument Reference above.

* `routes` - The list of route information with `destination` and `nexthop` fields.

* `shared` - Specifies whether the cross-tenant sharing is supported.
