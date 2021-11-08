---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_policy_v3

Manages a Dedicated Load Balancer Policy resource within OpenTelekomCloud.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V3 ELB client.
  If omitted, the `region` argument of the provider is used.
  Changing this creates a new Policy.

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


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the policy.

* `status` - Specifies the provisioning status of the forwarding policy.

## Import

Load Balancer Policy can be imported using the Policy ID, e.g.:

```shell
terraform import opentelekomcloud_lb_policy_v3.this 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
