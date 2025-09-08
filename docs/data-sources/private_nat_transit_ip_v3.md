---
subcategory: "Private NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_private_nat_transit_ip_v3"
sidebar_current: "docs-opentelekomcloud-data-source-private-nat-transit-ip-v3"
description: |-
  Manages a Private NAT Transit IP data source within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Private NAT Transit IP you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/apis_for_private_nat_gateways_v3.0/transit_ip_addresses/index.html)

# opentelekomcloud_private_nat_transit_ip_v3

Manages a V3 Private NAT Transit IP data source within OpenTelekomCloud.

## Example Usage

### List all Transit IPs

```hcl
data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {}
```

### List all Transit IPs in a subnet

```hcl
variable "network_id" {}

data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  virsubnet_id = var.network_id
}
```

### Get Transit IP by ID

```hcl
variable "id" {}

data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  id = var.id
}
```

### Get Transit IP by IP address

```hcl
variable "ip_address" {}

data "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  ip_address = var.ip_address
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) Specifies the private NAT transit IP ID.

* `virsubnet_id` - (Optional, String) Specifies the subnet ID of the current project.

* `ip_address` - (Optional, String) Specifies the transit IP address.


## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `transit_ips` - The list of private NAT transit IPs. The structure is defined below.

The `transit_ips` block supports: 

* `id` - Private NAT Transit IP ID.

* `virsubnet_id` - Indicates the subnet ID of the current project.

* `ip_address` - Indicates the transit IP address.

* `enterprise_project_id` - Indicates the ID of the enterprise project that is associated with the transit IP address when the transit IP address is being assigned.

* `tags` - Indicates the tag list in key/value format.

* `project_id` - Indicates the project ID.

* `network_interface_id` - Indicates the network interface ID of the transit IP address.

* `created_at` - Indicates the time when the private NAT Transit IP was created. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `updated_at` - Indicates the time when the private NAT Transit IP was updated. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `gateway_id` - Indicates the ID of the private NAT gateway associated with the transit IP address.

* `status` - Indicates the private NAT transit IP status. The value can be: `ACTIVE` (The Transit IP is running properly) or `FROZEN` (The Transit IP is frozen).
