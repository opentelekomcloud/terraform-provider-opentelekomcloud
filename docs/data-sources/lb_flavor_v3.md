---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_flavor_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lb-flavor-v3"
description: |-
Get ELBv3 flavor from OpenTelekomCloud
---

Up-to-date reference of API arguments for ELBv3 flavor you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/load_balancer_flavor/querying_flavors.html#listflavors)

# opentelekomcloud_lb_flavor_v3

Use this data source to get the info about an existing ELBv3 flavor.

## Example Usage

```hcl
data "opentelekomcloud_lb_flavor_v3" "l7_s2_small" {
  name = "L7_flavor.elb.s2.small"
}
```

## Argument Reference

* `id` - (Optional) Specifies the flavor ID.

* `name` - (Optional) Specifies the flavor name.

## Attributes Reference

In addition, the following attributes are exported:

* `type` - Specifies the flavor type.

* `shared` - Specifies whether the flavor is available to all users.

* `max_connections` - Specifies the maximum concurrent connections.

* `cps` - Specifies the number of new connections per second.

* `qps` - Specifies the number of requests per second at Layer 7.

* `bandwidth` - Specifies the inbound and outbound bandwidth in the unit of Kbit/s.
