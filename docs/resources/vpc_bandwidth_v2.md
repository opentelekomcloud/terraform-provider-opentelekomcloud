---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_bandwidth_v2"
sidebar_current: "docs-opentelekomcloud-resource-vpc-bandwidth-v2"
description: |-
Manages a VPC Bandwidth resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC bandwidth you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/bandwidth_v2.0)

# opentelekomcloud_vpc_bandwidth_v2

Provides a resource to create a shared bandwidth within Open Telekom Cloud.

## Example Usage

```hcl
resource "opentelekomcloud_vpc_bandwidth_v2" "band_100mb" {
  name = "shared-100Mbit"
  size = 100
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the bandwidth name.

  The value is a string of 1 to 64 characters that can contain letters, digits, underscores (_), hyphens (-), and periods (.).

* `size` - (Required) Specifies the bandwidth size.
  The value ranges from 5 Mbit/s to 1000 Mbit/s by default.

->
  The specific range may vary depending on the configuration in each region.
  You can see the available bandwidth range on the management console.

## Attributes Reference

In addition, the following attributes are exported:

* `id` - Specifies the bandwidth ID, which uniquely identifies the bandwidth.

* `status` - Specifies the bandwidth status.

## Import

VPC bandwidth can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_bandwidth_v2.band_100mb eb187fc8-e482-43eb-a18a-9da947ef89f6
```
