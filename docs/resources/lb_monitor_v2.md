---
subcategory: "Elastic Load Balancer (ELB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_monitor_v2"
sidebar_current: "docs-opentelekomcloud-resource-lb-monitor-v2"
description: |-
Manages a ELB Monitor resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for ELB monitor you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v2.0/health_check)

# opentelekomcloud_lb_monitor_v2

Manages an Enhanced LB monitor resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_monitor_v2" "monitor_1" {
  pool_id     = opentelekomcloud_lb_pool_v2.pool_1.id
  type        = "HTTP"
  delay       = 20
  timeout     = 10
  max_retries = 5
  url_path    = "/"
}
```

## Argument Reference

The following arguments are supported:

* `pool_id` - (Required) The id of the pool that this monitor will be assigned to.

* `name` - (Optional) The Name of the Monitor.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the monitor. Only administrative users can specify a tenant UUID
  other than their own. Changing this creates a new monitor.

* `type` - (Required) The type of probe, which is `TCP`, `UDP_CONNECT`, or `HTTP`,
  that is sent by the load balancer to verify the member state. Changing this
  creates a new monitor.

* `delay` - (Required) The time, in seconds, between sending probes to members.

* `timeout` - (Required) Maximum number of seconds for a monitor to wait for a
  ping reply before it times out. The value must be less than the delay value.

* `max_retries` - (Required) Number of permissible ping failures before
  changing the member's status to INACTIVE. Must be a number between 1 and 10.

* `admin_state_up` - (Optional) The administrative state of the monitor.
  A valid value is `true` (`UP`) or `false` (`DOWN`).

* `http_method` - (Optional) Required for HTTP types. The HTTP method used
  for requests by the monitor. If this attribute is not specified, it
  defaults to `GET`. The value can be `GET`, `HEAD`, `POST`, `PUT`, `DELETE`,
  `TRACE`, `OPTIONS`, `CONNECT`, and `PATCH`.

-> These parameters `domain_name`, `url_path`, `expected_codes` and `monitor_port`
  are valid when the value of `type` is set to `HTTP`.

* `domain_name` - (Optional) The `domain_name` of the HTTP request during the health check.

* `url_path` - (Optional) Required for HTTP types. URI path that will be
  accessed if monitor type is `HTTP`.

* `expected_codes` - (Optional) Required for `HTTP` types. Expected HTTP codes
  for a passing HTTP monitor. You can either specify a single status like
  `"200"`, or a list like `"200,202"`.

* `monitor_port` - (Optional) Specifies the health check port. The port number
  ranges from 1 to 65535. The value is left blank by default, indicating that
  the port of the backend server is used as the health check port.


## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the monitor.

* `tenant_id` - See Argument Reference above.

* `type` - See Argument Reference above.

* `delay` - See Argument Reference above.

* `timeout` - See Argument Reference above.

* `max_retries` - See Argument Reference above.

* `url_path` - See Argument Reference above.

* `domain_name` - See Argument Reference above.

* `http_method` - See Argument Reference above.

* `expected_codes` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `monitor_port` - See Argument Reference above.
