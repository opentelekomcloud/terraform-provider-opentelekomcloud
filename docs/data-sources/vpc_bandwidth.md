---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_bandwidth"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-bandwidth"
description: |-
Get details about a specific shared bandwidth from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPC shared bandwidth you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/bandwidth/querying_bandwidths.html#vpc-bandwidth-0002)

# opentelekomcloud_vpc_bandwidth

Provides details about a specific shared bandwidth.

~>
Deprecated, use `opentelekomcloud_vpc_bandwidth_v2` data source instead

## Example Usage

```hcl
variable "bandwidth_name" {}

data "opentelekomcloud_vpc_bandwidth" "bandwidth_1" {
  name = var.bandwidth_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
bandwidth in the current tenant. The following arguments are supported:

* `region` - (Optional) The region in which to obtain the bandwidth. If omitted, the provider-level region will be used.

* `name` - (Required) The name of the Shared Bandwidth to retrieve.

* `size` - (Optional) The size of the Shared Bandwidth to retrieve. The value ranges from 5 Mbit/s to 2000 Mbit/s.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the Shared Bandwidth.

* `name` -  See Argument Reference above.

* `size` - See Argument Reference above.

* `share_type` - Indicates whether the bandwidth is a shared or dedicated one.

* `bandwidth_type` - Indicates the bandwidth type.

* `charge_mode` - Specifies that the bandwidth is billed by bandwidth. The value can be traffic.

* `status` - Indicates the bandwidth status.
