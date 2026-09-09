---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_subnet_private_ip_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-subnet-private-ip-v1"
description: |-
  Queries a private IP address in an OpenTelekomCloud VPC subnet.
---

# opentelekomcloud_vpc_subnet_private_ip_v1

Use this data source to query a VPC subnet private IP address by ID.

## Example Usage

```hcl
data "opentelekomcloud_vpc_subnet_private_ip_v1" "private_ip" {
  id = "d600542a-b231-45ed-af05-e9930cb14f78"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) Specifies the unique private IP address ID.

* `region` - (Optional, String) Specifies the region in which to query the private IP address.
  If omitted, the provider-level region is used.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `subnet_id` - The network ID of the subnet from which the address is assigned.

* `ip_address` - The assigned IPv4 address.

* `status` - The private IP address status. The value is `ACTIVE` or `DOWN`.

* `tenant_id` - The project ID that owns the private IP address.

* `device_owner` - The resource using the private IP address. This is empty when the address is unused.
