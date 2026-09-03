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

* `enterprise_project_id` - (Optional) Specifies the enterprise project associated
  with the shared bandwidth. If omitted, the provider-level enterprise project or
  the default enterprise project (`0`) is used. Changing this creates a new resource.

* `public_border_group` - (Optional) Specifies whether the shared bandwidth is
  located at the central site or an edge site. Changing this creates a new resource.

## Attributes Reference

In addition, the following attributes are exported:

* `id` - Specifies the bandwidth ID, which uniquely identifies the bandwidth.

* `status` - Specifies the bandwidth status.

* `share_type` - Indicates whether the bandwidth is shared or dedicated.

* `bandwidth_type` - Specifies the bandwidth type.

* `charge_mode` - Specifies the bandwidth charging mode.

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

## Import

VPC bandwidth can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_bandwidth_v2.band_100mb eb187fc8-e482-43eb-a18a-9da947ef89f6
```
