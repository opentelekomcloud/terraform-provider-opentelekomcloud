---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_bandwidth_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-bandwidth-v2"
description: |-
  Get details about a specific shared bandwidth from OpenTelekomCloud
---

# opentelekomcloud_vpc_bandwidth_v2

Provides details about a specific shared bandwidth.

## Example Usage

```hcl
variable "bandwidth_name" {}

data "opentelekomcloud_vpc_bandwidth_v2" "bandwidth_1" {
  name = var.bandwidth_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
bandwidth in the current tenant. The following arguments are supported:

* `id` - (Optional) The ID of the shared bandwidth to retrieve.

* `name` - (Optional) The name of the shared bandwidth to retrieve.

* `size` - (Optional) The size of the shared bandwidth to retrieve.

* `share_type` - (Optional) Indicates whether the bandwidth is shared or dedicated.

* `enterprise_project_id` - (Optional) Specifies the enterprise project associated
  with the bandwidth.

* `public_border_group` - (Optional) Specifies whether the bandwidth is located at
  the central site or an edge site.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the Shared Bandwidth.

* `name` -  See Argument Reference above.

* `size` - See Argument Reference above.

* `share_type` - Indicates whether the bandwidth is a shared or dedicated one.

* `bandwidth_type` - Indicates the bandwidth type.

* `charge_mode` - Specifies that the bandwidth is billed by bandwidth. The value can be traffic.

* `status` - Indicates the bandwidth status.

* `billing_info` - Specifies yearly/monthly billing information, when applicable.

* `tenant_id` - Specifies the project ID.

* `created_at` - Specifies when the bandwidth was created.

* `updated_at` - Specifies when the bandwidth was last updated.

* `publicip_info` - Lists EIPs associated with the bandwidth:
  * `id` - EIP ID.
  * `address` - IPv4 address.
  * `ipv6_address` - IPv6 address, when available.
  * `ip_version` - IP address version.
  * `type` - EIP type.
