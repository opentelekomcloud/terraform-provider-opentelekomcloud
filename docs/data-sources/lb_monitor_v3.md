---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_monitor_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lb-monitor-v3"
description: |-
  Looks up a Dedicated Load Balancer v3 health monitor.
---

# opentelekomcloud_lb_monitor_v3

Looks up one Dedicated Load Balancer health monitor by ID or by supported API filters.
When using filters, the query must match exactly one monitor across all result pages.

## Example Usage

```hcl
data "opentelekomcloud_lb_monitor_v3" "by_id" {
  id = var.monitor_id
}

data "opentelekomcloud_lb_monitor_v3" "disabled" {
  name           = "application-health"
  admin_state_up = false
}
```

## Argument Reference

All arguments are optional filters and are also exported as attributes:

* `id` - The monitor ID. When configured, a direct lookup is performed and other filters are ignored.
* `name` - The health monitor name.
* `type` - The health check protocol: `TCP`, `UDP_CONNECT`, `HTTP`, or `HTTPS`.
* `admin_state_up` - Whether health checks are enabled. Explicit `false` is sent to the API;
  leaving this argument unset does not apply an administrative-state filter.
* `delay` - The interval between health checks, in seconds, from 1 to 50.
* `timeout` - The response timeout, in seconds, from 1 to 50.
* `max_retries` - The consecutive successful checks required for recovery, from 1 to 10.
* `max_retries_down` - The consecutive failed checks required for failure, from 1 to 10.
* `monitor_port` - The health check port, from 1 to 65535.
* `domain_name` - The HTTP or HTTPS health check domain name.
* `url_path` - The health check request path, starting with `/`.
* `expected_codes` - The expected HTTP status code, list, or range.

Only explicitly configured arguments are used as filters. Computed values from a previous read do not
become filters or switch a filtered query into an ID lookup. Missing IDs, no matching monitors, multiple
matching monitors, and API failures produce errors.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `pool_id` - The associated backend server group ID, if returned by the API.
* `project_id` - The project ID.
* `http_method` - The returned HTTP method. This legacy response field is retained for consistency with
  the resource, but it is not exposed as a new query filter because the reviewed OTC API marks it unsupported.
* `created_at` - The creation time, in UTC format `yyyy-MM-dd'T'HH:mm:ss'Z'`.
* `updated_at` - The last update time, in the same UTC format.

The timestamp fields can be empty for older dedicated load balancers and for shared load balancers.
The stable data-source ID is the health monitor ID.

HTTP-version, content-verification, and enterprise-project filters marked unsupported by OTC are not
exposed. Undocumented legacy SDK filters are not added. The SDK handles pagination internally.

See the [OTC health check API reference](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/health_check/index.html).
