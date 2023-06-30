---
subcategory: "Dedicated Load Balancer (DLB)"
---

Up-to-date reference of API arguments for DLB monitor you can get at
`https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/health_check`.

# opentelekomcloud_lb_monitor_v3

Manages a Dedicated LB monitor (health check) resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = [var.availability_zone]
}


resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  protocol        = "HTTP"
  lb_algorithm    = "ROUND_ROBIN"
}

resource "opentelekomcloud_lb_monitor_v3" "monitor" {
  pool_id      = opentelekomcloud_lb_pool_v3.pool.id
  type         = "HTTP"
  delay        = 3
  timeout      = 30
  monitor_port = 8080

  max_retries      = 5
  max_retries_down = 1
}
```

## Argument Reference

The following arguments are supported:

* `admin_state_up` - (Optional) Specifies the administrative status of the health check.
  `true` indicates that the health check is enabled, and `false` indicates that the health check is disabled.

  Default: `true`

* `pool_id` - (Required) Specifies the ID of the backend server group for which the health check is configured.
  Changing this creates a new monitor.

* `type` - (Required) Specifies the health check protocol.

  The value can be `TCP`, `UDP_CONNECT`, `HTTP`, `HTTPS`, or `PING`.

* `delay` - (Required) Specifies the interval between health checks, in seconds.

  The value of this parameter ranges from 1 to 50.

* `timeout` - (Required) Specifies the maximum time required for waiting for a response from the health check, in
  seconds.

  The value of this parameter ranges from 1 to 50.

  It is recommended that you set the value less than that of parameter `delay`.

* `max_retries` - (Required) Specifies the number of consecutive health checks when the health check result of a backend
  server changes from `OFFLINE` to `ONLINE`.

  The value ranges from 1 to 10.

* `max_retries_down` - (Optional) Specifies the number of consecutive health checks when the health check result of a
  backend server changes from `ONLINE` to `OFFLINE`.

  The value ranges from 1 to 10.

  Default value is 3

* `monitor_port` - (Optional) Specifies the port used for the health check. If this parameter is left blank, the port of
  the backend server group will be used by default.

* `name` - (Optional) Specifies the health check name.

* `project_id` - (Optional) Specifies the project ID. Changing this creates a new monitor.

* `domain_name` - (Optional) Specifies the domain name that HTTP requests are sent to during the health check.

  This parameter is available only when type is set to `HTTP`.

  The value is left blank by default, indicating that the virtual IP address of the load balancer is used as the
  destination address of HTTP requests.

  The value can contain only digits, letters, hyphens (-), and periods (.) and must start with a digit or letter.

* `url_path` - (Optional) Specifies the HTTP request path for the health check.

  The value must start with a slash (/), and the default value is `/`.

  This parameter is available only when `type` is set to `HTTP`.

* `expected_codes` - (Optional) Specifies the expected HTTP status code. This parameter will take effect only
  when `type` is set to HTTP.

  The value options are as follows:
  * A specific value, for example, `200`
  * A list of values that are separated with commas (,), for example, `200, 202`
  * A value range, for example, `200-204`

  Default: `200`

* `http_method` - (Optional) Specifies the HTTP method.

  The value can be `GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `TRACE`, `OPTIONS`, `CONNECT`, or `PATCH`.

  This parameter will take effect only when type is set to `HTTP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies the health check (monitor) ID.

## Import

Load Balancer Monitor can be imported using the monitor ID, e.g.:

```shell
terraform import opentelekomcloud_lb_monitor_v3.monitor b4ef7345-cf1a-41ca-8baa-941466a66853
```
