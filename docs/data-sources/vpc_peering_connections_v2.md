---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_peering_connections_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-peering-connections-v2"
description: |-
  List VPC peering connections from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC peering connections you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/vpc_peering_connection/querying_vpc_peering_connections.html#vpc-peering-0001)

# opentelekomcloud_vpc_peering_connections_v2

Use this data source to list VPC peering connections matching the specified criteria.

## Example Usage

```hcl
data "opentelekomcloud_vpc_peering_connections_v2" "peerings" {
  vpc_id = opentelekomcloud_vpc_v1.vpc.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the VPC peering connection to filter by.

* `status` - (Optional) The status of the VPC peering connection to filter by.

* `vpc_id` - (Optional) The ID of the requester VPC to filter by.

* `peer_vpc_id` - (Optional) The ID of the accepter/peer VPC to filter by.

* `peer_tenant_id` - (Optional) The tenant ID of the accepter/peer VPC to filter by.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region` - The region of the VPC peering connections.

* `peering_connections` - A list of VPC peering connections. Each element contains the following attributes:

  * `id` - The ID of the VPC peering connection.

  * `name` - The name of the VPC peering connection.

  * `description` - The description of the VPC peering connection.

  * `status` - The status of the VPC peering connection.

  * `vpc_id` - The ID of the requester VPC.

  * `vpc_tenant_id` - The project ID the requester VPC belongs to.

  * `peer_vpc_id` - The ID of the accepter/peer VPC.

  * `peer_tenant_id` - The project ID the accepter VPC belongs to.
