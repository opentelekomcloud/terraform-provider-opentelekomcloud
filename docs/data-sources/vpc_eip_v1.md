---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_eip_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-eip-v1"
description: |-
  Get details about a specific VPC EIP from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC EIP you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/eip/querying_eips.html#vpc-eip-0003)

# opentelekomcloud_vpc_eip_v1

Use this data source to get details about a specific VPC elastic IP.

## Example Usage

```hcl
data "opentelekomcloud_vpc_eip_v1" "eip_v1" {
  id = var.elastic_ip
}

output "eip_vpc_id" {
  value = data.opentelekomcloud_vpc_eip_v1.eip_v1.id
}
```

## Search by name regex

```hcl
resource "opentelekomcloud_vpc_eip_v1" "eip" {
  publicip {
    type = "5_bgp"
    name = "my_eip"
  }
  bandwidth {
    name        = "acc-band"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

data "opentelekomcloud_vpc_eip_v1" "by_regex" {
  name_regex = "^my_.+"
}

output "eip_vpc_id" {
  value = data.opentelekomcloud_vpc_eip_v1.by_regex.name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
elastic IP in the current tenant. The given filters must match exactly one
elastic IP whose data will be exported as attributes.

* `id` - (Optional) Specifies a resource ID in UUID format.

* `name_regex` - (Optional) A regex string to apply to the eip list. This allows more advanced filtering.

* `status` - (Optional) The status of the specific elastic IP to retrieve.

* `public_ip_address` - (Optional) The public IP address of the elastic IP.

* `private_ip_address` - (Optional) The private IP address bound to the elastic IP.

* `port_id` - (Optional) The port ID.

-> `private_ip_address` and `port_id` are returned only when a port/private IP address is
associated with the elastic IP.

* `bandwidth_id` - (Optional) The bandwidth ID of specific elastic IP.

* `tags` - (Optional) Tags key/value pairs to filter the elastic IPs.

## Attributes Reference

All the argument attributes are also exported as result attributes.

* `ip_version` - The IP version of elastic IP.

* `bandwidth_share_type` - Specifies the EIP bandwidth type.

* `bandwidth_size` - Specifies the bandwidth (Mbit/s).

* `create_time` - Specifies the time (UTC) when the elastic IP is assigned.

* `tenant_id` - Specifies the project ID.

* `type` - Specifies the elastic IP type.

* `name` - Specifies the elastic IP Name.
