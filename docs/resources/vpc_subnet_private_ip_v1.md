---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_subnet_private_ip_v1"
sidebar_current: "docs-opentelekomcloud-resource-vpc-subnet-private-ip-v1"
description: |-
  Manages a private IP address in an OpenTelekomCloud VPC subnet.
---

# opentelekomcloud_vpc_subnet_private_ip_v1

Manages a private IP address in a VPC subnet.

## Example Usage

### Automatically assign an available address

```hcl
resource "opentelekomcloud_vpc_subnet_private_ip_v1" "private_ip" {
  subnet_id = opentelekomcloud_vpc_subnet_v1.example.network_id
}
```

### Assign a specific address

```hcl
resource "opentelekomcloud_vpc_subnet_private_ip_v1" "private_ip" {
  subnet_id  = opentelekomcloud_vpc_subnet_v1.example.network_id
  ip_address = "192.168.10.10"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, ForceNew) Specifies the region in which to create the private IP address.
  If omitted, the provider-level region is used. Changing this creates a new resource.

* `subnet_id` - (Required, ForceNew) Specifies the network ID of the subnet from which the address
  is assigned. Use the `network_id` attribute of `opentelekomcloud_vpc_subnet_v1`. Changing this
  creates a new resource.

* `ip_address` - (Optional, ForceNew) Specifies an available IPv4 address in the subnet. If omitted,
  OpenTelekomCloud automatically assigns an address. Changing this creates a new resource.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The unique private IP address ID.

* `status` - The private IP address status. The value is `ACTIVE` or `DOWN`.

* `tenant_id` - The project ID that owns the private IP address.

* `device_owner` - The resource using the private IP address. This is empty when the address is unused.

## Import

Private IP addresses can be imported using their ID:

```shell
terraform import opentelekomcloud_vpc_subnet_private_ip_v1.private_ip <private-ip-id>
```
