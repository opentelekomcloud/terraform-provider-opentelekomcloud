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

* `name` - (Optional) The name of the Shared Bandwidth to retrieve.

* `size` - (Optional) The size of the Shared Bandwidth to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the Shared Bandwidth.

* `name` -  See Argument Reference above.

* `size` - See Argument Reference above.

* `share_type` - Indicates whether the bandwidth is a shared or dedicated one.

* `bandwidth_type` - Indicates the bandwidth type.

* `charge_mode` - Specifies that the bandwidth is billed by bandwidth. The value can be traffic.

* `status` - Indicates the bandwidth status.
