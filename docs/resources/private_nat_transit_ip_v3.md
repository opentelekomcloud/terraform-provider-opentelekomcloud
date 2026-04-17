---
subcategory: "Private NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_private_nat_transit_ip_v3"
sidebar_current: "docs-opentelekomcloud-resource-private-nat-transit-ip-v3"
description: |-
  Manages a Private NAT Transit IP resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Private NAT Transit IP you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/apis_for_private_nat_gateways_v3.0/transit_ip_addresses/index.html)

# opentelekomcloud_private_nat_transit_ip_v3

Manages a V3 Private NAT Transit IP resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "network_id" {}

resource "opentelekomcloud_private_nat_transit_ip_v3" "transit_ip_1" {
  virsubnet_id = var.network_id
  tags = {
    kuh = "muh"
  }
}
```

## Argument Reference

The following arguments are supported:

* `virsubnet_id` - (Required, String, ForceNew) Specifies the subnet ID of the current project.

* `ip_address` - (Optional, String, ForceNew) Specifies the transit IP address.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project that is associated with the transit IP address when the transit IP address is being assigned. Default: `0`.

* `tags` - (Optional, Map, ForceNew) Specifies the tag list in key/value format.


## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Private NAT Transit IP ID.

* `project_id` - Indicates the project ID.

* `network_interface_id` - Indicates the network interface ID of the transit IP address.

* `created_at` - Indicates the time when the private NAT Transit IP was created. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `updated_at` - Indicates the time when the private NAT Transit IP was updated. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `gateway_id` - Indicates the ID of the private NAT gateway associated with the transit IP address.

* `status` - Indicates the private NAT transit IP status. The value can be: `ACTIVE` (The Transit IP is running properly) or `FROZEN` (The Transit IP is frozen).


## Import

Private NAT Transit IP V3 resource can be imported using the Transit IP ID, `id`.

```sh
terraform import opentelekomcloud_private_nat_transit_ip_v3.ip_1 <id>
```
