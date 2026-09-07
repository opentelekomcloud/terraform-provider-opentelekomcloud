---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_quotas_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-quotas-v1"
description: |-
  Query VPC network resource quotas from T-Cloud Public (former OpenTelekomCloud).
---

Up-to-date reference of API arguments for VPC Peering Connections you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/vpc_apis_v1_v2/quota/index.html)

# opentelekomcloud_vpc_quotas_v1

Use this data source to query VPC network resource quotas.

## Example Usage

```hcl
data "opentelekomcloud_vpc_quotas_v1" "all" {}
```

### Filter by resource type

```hcl
data "opentelekomcloud_vpc_quotas_v1" "eip" {
  type = "publicIp"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Optional, String) Specifies the resource type. Valid values are `vpc`, `subnet`,
  `securityGroup`, `securityGroupRule`, `publicIp`, `vpn`, `vpcPeer`, `loadbalancer`, `listener`,
  `physicalConnect`, `virtualInterface`, `firewall`, `shareBandwidthIP`, `shareBandwidth`,
  `address_group`, `flow_log`, `vpcContainRoutetable`, and `routetableContainRoutes`.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `region` - The region in which the quotas are queried.

* `quotas` - The matching network resource quotas.
  The [quotas](#quotas) structure is documented below.

<a name="quotas"></a>
The `quotas` block supports:

* `type` - The network resource type.

* `used` - The number of created resources.

* `quota` - The maximum number of resources. A value of `-1` means that the quota is unlimited.

* `min` - The minimum quota value allowed.
