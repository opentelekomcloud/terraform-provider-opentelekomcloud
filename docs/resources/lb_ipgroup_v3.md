---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_ipgroup_v3"
sidebar_current: "docs-opentelekomcloud-resource-lb-ipgroup-v3"
description: |-
Manages a LB IpGroup resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DLB ip group you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/ip_address_group)

# opentelekomcloud_lb_ipgroup_v3

Manages a Dedicated Load Balancer IP address group resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "group description"

  ip_list {
    ip          = "192.168.50.10"
    description = "one"
  }
  ip_list {
    ip          = "192.168.100.10"
    description = "two"
  }
  ip_list {
    ip          = "192.168.150.10"
    description = "three"
  }
}
```

## Example empty ip list

```hcl
resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "group description"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the IP address group name.

* `description` - (Optional, String) Provides supplementary information about the IP address group.

* `project_id` - (Optional, String) Specifies the project ID of the IP address group.

* `ip_list` - (Optional, List) Specifies the IP addresses or CIDR blocks in the IP address group.
    Any IP address can be used if this block isn't specified.
  * `ip` - (Required, String) Specifies the IP addresses in the IP address group.
    IPv6 is unsupported. The value cannot be an IPv6 address.
  * `description` - (Optional, String) Provides remarks about the IP address group.

## Attributes Reference

In addition, the following attributes are exported:

* `listeners` - Lists the IDs of listeners with which the IP address group is associated.

* `updated_at` - Indicates the update time.

* `created_at` - Indicates the creation time.

## Import

Ip groups can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_lb_ipgroup_v3.group_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
