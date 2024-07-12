---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_peering_connection_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-peering-connection-v2"
description: |-
  Get details about a specific VPC peering connection from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC EIP you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/vpc_peering_connection/querying_vpc_peering_connections.html#vpc-peering-0001)

# opentelekomcloud_vpc_peering_connection_v2

Use this data source to get details about a specific VPC peering connection.

## Example Usage

```hcl
data "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  vpc_id      = opentelekomcloud_vpc_v1.vpc.id
  peer_vpc_id = opentelekomcloud_vpc_v1.peer_vpc.id
}


resource "opentelekomcloud_vpc_route_v2" "vpc_route" {
  type        = "peering"
  nexthop     = data.opentelekomcloud_vpc_peering_connection_v2.peering.id
  destination = "192.168.0.0/16"
  vpc_id      = opentelekomcloud_vpc_v1.vpc.id
}
```


## Argument Reference

The arguments of this data source act as filters for querying the available VPC peering connection.
The given filters must match exactly one VPC peering connection whose data will be exported as attributes.

* `id` - (Optional) The ID of the specific VPC Peering Connection to retrieve.

* `status` - (Optional) The status of the specific VPC Peering Connection to retrieve.

* `vpc_id` - (Optional) The ID of the requester VPC of the specific VPC Peering Connection to retrieve.

* `peer_vpc_id` - (Optional)  The ID of the accepter/peer VPC of the specific VPC Peering Connection to retrieve.

* `peer_tenant_id` - (Optional) The Tenant ID of the accepter/peer VPC of the specific VPC Peering Connection to retrieve.

* `name` - (Optional) The name of the specific VPC Peering Connection to retrieve.


## Attributes Reference

All of the argument attributes are exported as result attributes.
