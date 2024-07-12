---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_subnet_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-subnet-v1"
description: |-
  Get details about a specific VPC subnet from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC subnet you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/subnet/querying_subnets.html#vpc-subnet01-0003)

# opentelekomcloud_vpc_subnet_v1

Use this data source to get details about a specific VPC subnet.

This data source can prove useful when a module accepts a subnet id as
an input variable and needs to, for example, determine the id of the
VPC that the subnet belongs to.

## Example Usage

```hcl
data "opentelekomcloud_vpc_subnet_v1" "subnet_v1" {
  id = var.subnet_id
}

output "subnet_vpc_id" {
  value = data.opentelekomcloud_vpc_subnet_v1.subnet_v1.vpc_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
subnets in the current tenant. The given filters must match exactly one
subnet whose data will be exported as attributes.

* `id` - (Optional) Specifies a resource ID in UUID format.

* `name` - (Optional) The name of the specific subnet to retrieve.

* `cidr` - (Optional) The network segment of specific subnet to retrieve. The value must be in CIDR format.

* `status` - (Optional) The value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

* `vpc_id` - (Optional) The id of the VPC that the desired subnet belongs to.

* `gateway_ip` - (Optional) The subnet gateway address of specific subnet.

* `primary_dns` - (Optional) The IP address of DNS server 1 on the specific subnet.

* `secondary_dns` - (Optional) The IP address of DNS server 2 on the specific subnet.

* `availability_zone` - (Optional) The availability zone (AZ) to which the subnet should belong.

## Attributes Reference

All the argument attributes are also exported as result attributes.
This data source will complete the data by populating
any fields that are not included in the configuration with the data for
the selected subnet.

* `dns_list` - The IP address list of DNS servers on the subnet.

* `dhcp_enable` - DHCP function for the subnet.

* `subnet_id` - Specifies the OpenStack subnet ID.

* `network_id` - Specifies the OpenStack network ID.
