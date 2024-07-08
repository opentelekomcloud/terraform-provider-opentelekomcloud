---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_pool_v3"
sidebar_current: "docs-opentelekomcloud-resource-lb-pool-v3"
description: |-
Manages a LB Pool resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DLB pool you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/backend_server_group)

# opentelekomcloud_lb_pool_v3

Manages a Dedicated LB pool resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = var.router_id
  network_ids = [var.network_id]

  availability_zones = [var.availability_zone]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  name            = "pool_1"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "TCP"

  session_persistence {
    type                = "SOURCE_IP"
    persistence_timeout = "30"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specifies the backend server group name.

* `description` - (Optional) Provides supplementary information about the backend server group.

* `protocol` - (Required) Specifies the protocol used by the backend server group to receive requests.
  `TCP`, `UDP`, `HTTP`, `HTTPS`, and `QUIC` are supported.

  * For `UDP` listeners, the protocol of the backend server group must be `UDP`.
  * For `TCP` listeners, the protocol of the backend server group must be `TCP`.
  * For `HTTP` listeners, the protocol of the backend server group must be `HTTP`.
  * For `HTTPS` listeners, the protocol of the backend server group must be `HTTPS`.

* `lb_algorithm` - (Required) Specifies the load balancing algorithm used by the load balancer to route requests to backend servers.

  The value can be one of the following:
  * `ROUND_ROBIN`: weighted round-robin
  * `LEAST_CONNECTIONS`: weighted least connections
  * `SOURCE_IP`: source IP hash

  When the value is `SOURCE_IP`, the weights of backend servers are invalid.

* `listener_id` - (Optional) Specifies the ID of the listener associated with the backend server group.

* `loadbalancer_id` - (Optional) Specifies the ID of the associated load balancer.

-> Specify either `listener_id` or `loadbalancer_id`, or **both** of them.

* `project_id` - (Optional) Specifies the project ID.

* `session_persistence` - (Optional) Specifies whether to enable sticky sessions.

The `session_persistence` block supports:

* `type` - (Required) Specifies the sticky session type. The value can be `SOURCE_IP`, `HTTP_COOKIE`, or `APP_COOKIE`.

  * If the protocol of the backend server group is `TCP`, `UDP`, and `QUIC`, only `SOURCE_IP` takes effect.

  * For dedicated load balancers, if the protocol of the backend server group is `HTTP` or `HTTPS`, the value can only be `HTTP_COOKIE`.

  * For shared load balancers, if the protocol of the backend server group is `HTTP` or `HTTPS`, the value can be `HTTP_COOKIE` or `APP_COOKIE`.

* `cookie_name` - (Optional) Specifies the cookie name. This parameter will take effect only when type is set to `APP_COOKIE`.
  The value can contain only letters, digits, hyphens (-), underscores (_), and periods (.).

* `persistence_timeout` - (Optional) Specifies the stickiness duration, in minutes.
  This parameter will not take effect when type is set to `APP_COOKIE`.
  * If the protocol of the backend server group is TCP or UDP,
  the value ranges from `1` to `60`, and the default value is `1`.
  * If the protocol of the backend server group is HTTP or HTTPS, the value ranges from `1` to `1440`,
  and the default value is `1440`.

* `member_deletion_protection` - (Optional) Specifies whether to enable removal protection for the pool members.
  `true`: Enable removal protection.
  `false` (default): Disable removal protection.

* `vpc_id` - (Optional) Specifies the ID of the VPC where the backend server group works.

* `type` - (Optional) Specifies the type of the backend server group.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies the backend server group ID.

* `ip_version` - Specifies the IP version supported by the backend server group.

## Import

Pools can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_lb_pool_v3.pool 7b80e108-1636-44e5-aece-986b0052b7dd
```
