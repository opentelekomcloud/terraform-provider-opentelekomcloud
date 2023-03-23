---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_policy_v3

Manages a Dedicated Load Balancer Policy resource within OpenTelekomCloud.

## Example Usage Basic

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "this" {
  router_id   = var.router_id
  network_ids = [var.network_id]

  availability_zones = [var.az]
}

resource "opentelekomcloud_lb_listener_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  protocol        = "HTTP"
  protocol_port   = 8080
}

resource "opentelekomcloud_lb_pool_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
}

resource "opentelekomcloud_lb_policy_v3" "this" {
  action           = "REDIRECT_TO_POOL"
  listener_id      = opentelekomcloud_lb_listener_v3.this.id
  redirect_pool_id = opentelekomcloud_lb_pool_v3.this.id
  position         = 37
}
```

## Fixed Response Example

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "this" {
  router_id   = var.router_id
  network_ids = [var.network_id]

  availability_zones = [var.az]
}

resource "opentelekomcloud_lb_listener_v3" "this" {
  loadbalancer_id     = opentelekomcloud_lb_loadbalancer_v3.this.id
  protocol            = "HTTP"
  protocol_port       = 8080
  advanced_forwarding = true
}

resource "opentelekomcloud_lb_pool_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
}

resource "opentelekomcloud_lb_policy_v3" "this" {
  action      = "FIXED_RESPONSE"
  listener_id = opentelekomcloud_lb_listener_v3.this.id
  position    = 37
  priority    = 10

  fixed_response_config {
    status_code  = "200"
    content_type = "text/plain"
    message_body = "Fixed Response"
  }
}
```

## Redirect To Url Example

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "this" {
  router_id   = var.router_id
  network_ids = [var.network_id]

  availability_zones = [var.az]
}

resource "opentelekomcloud_lb_listener_v3" "this" {
  loadbalancer_id     = opentelekomcloud_lb_loadbalancer_v3.this.id
  protocol            = "HTTP"
  protocol_port       = 8080
  advanced_forwarding = true
}

resource "opentelekomcloud_lb_pool_v3" "this" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.this.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
}

resource "opentelekomcloud_lb_policy_v3" "this" {
  action      = "REDIRECT_TO_URL"
  listener_id = opentelekomcloud_lb_listener_v3.this.id
  position    = 37
  priority    = 10

  redirect_url = "https://www.google.com:443"

  redirect_url_config {
    status_code = "301"
    query       = "name=my_name"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the Policy. Only administrative users can specify a tenant UUID other than
  their own. Changing this creates a new Policy.

* `name` - (Optional) Specifies the forwarding policy name.

* `description` - (Optional) Provides supplementary information about the forwarding policy.

* `action` - (Required) The Policy action - can either be `REDIRECT_TO_POOL`,
  or `REDIRECT_TO_LISTENER`. Changing this creates a new Policy.

* `listener_id` - (Required) The Listener on which the Policy will be associated with.
  Changing this creates a new Policy.

* `position` - (Optional) The position of this policy on the listener. Positions start at `1`.
  Changing this creates a new Policy.

* `redirect_pool_id` - (Optional) Requests matching this policy will be redirected to the pool with this ID.
  Only valid if `action` is `REDIRECT_TO_POOL`.

* `redirect_listener_id` - (Optional) Requests matching this policy will be redirected to the listener with this ID.
  Only valid if `action` is `REDIRECT_TO_LISTENER`.

* `rules` - (Optional) Lists the forwarding rules in the forwarding policy.
  * `type` - (Required) Specifies the match content. The value can be one of the following: `HOST_NAME`, `PATH`.
  * `compare_type` - (Required) - Specifies how requests are matched with the domain name or URL.
    The values can be: `EQUAL_TO`, `REGEX`, `STARTS_WITH`.

  ->If `type` is set to `HOST_NAME`, this parameter can only be set to `EQUAL_TO` (exact match).
  If `type` is set to `PATH`, this parameter can be set to `REGEX` (regular expression match),
  `STARTS_WITH` (prefix match), or `EQUAL_TO` (exact match).

  * `value` - (Required) Specifies the value of the match item. For example, if a domain name is
    used for matching, value is the domain name.

  ->If type is set to `HOST_NAME`, the value can contain letters, digits, hyphens `-`, and periods `.`
  and must start with a letter or digit. If you want to use a wildcard domain name, enter an asterisk `*`
  as the leftmost label of the domain name.
  If type is set to `PATH` and `compare_type` to `STARTS_WITH` or `EQUAL_TO`, the value must start with
  a slash `/` and can contain only letters, digits, and special characters `_~';@^-%#&$.*+?,=!:|/()[]{}`.

* `priority` - (Optional) Specifies the forwarding policy priority.
  A smaller value indicates a higher priority. The value must be unique for forwarding policies of the same listener.
  This parameter will take effect only when `advanced_forwarding` is set to `true`.
  If this parameter is passed and `advanced_forwarding` is set to `false`, an error will be returned.
  This parameter is unsupported for shared load balancers and not available in `eu-nl`.

* `fixed_response_config` - (Optional) Specifies the configuration of the page that will be returned.
  This parameter will take effect when `advanced_forwarding` is set to `true`.
  If this parameter is passed and `advanced_forwarding` is set to `false`, an error will be returned.
  Not available in `eu-nl`.
  * `status_code` - (Required) Specifies the fixed HTTP status code configured in the forwarding rule.
    The value can be any integer in the range of `200-299`, `400-499`, or `500-599`.
  * `content_type` - (Optional) - Specifies the format of the response body.
  * `message_body` - (Optional) - Specifies the content of the response message body.

* `redirect_url` - (Optional) Specifies the URL to which requests are forwarded.

* `redirect_url_config` - (Optional) Specifies the URL to which requests are forwarded.
  For dedicated load balancers, This parameter will take effect when `advanced_forwarding` is set to `true`.
  If it is passed when `advanced_forwarding` is set to `false`, an error will be returned. Not available in `eu-nl`.
  * `protocol` - (Optional) - Specifies the protocol for redirection. The value can be `HTTP`, `HTTPS`,
    or `${protocol}`.
    The default value is `${protocol}`, indicating that the protocol of the request will be used.
  * `host` - (Optional) - Specifies the host name that requests are redirected to.
    The value can contain only letters, digits, hyphens (-), and periods (.) and must start with a letter or digit.
    The default value is `${host}`, indicating that the host of the request will be used.
  * `port` - (Optional) - Specifies the port that requests are redirected to. The default value is `${port}`,
    indicating that the port of the request will be used.
  * `path` - (Optional) - Specifies the path that requests are redirected to.
    The default value is `${path}`, indicating that the path of the request will be used.
    The value can contain only letters, digits, and special characters `_~';@^- %#&$.*+?,=!:|/()[]{}`
    and must start with a slash (`/`).
  * `query` - (Optional) - Specifies the query string set in the URL for redirection.
    The default value is `${query}`, indicating that the query string of the request will be used.
  * `status_code` - (Required) - Specifies the status code returned after the requests are redirected.
    The value can be `301`, `302`, `303`, `307`, or `308`.

* `redirect_pools_config` - (Optional) Specifies the configuration of the backend server group that the requests
  are forwarded to. This parameter is valid only when action is set to `REDIRECT_TO_POOL`.
  * `pool_id` - (Required) - Specifies the ID of the backend server group.
  * `weight` - (Required) - Specifies the weight of the backend server group. The value ranges from 0 to 100.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the policy.

* `status` - Specifies the provisioning status of the forwarding policy.

## Import

Load Balancer Policy can be imported using the Policy ID, e.g.:

```shell
terraform import opentelekomcloud_lb_policy_v3.this 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
