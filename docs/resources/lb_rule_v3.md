---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_rule_v3

Manages a Dedicated Load Balancer Rule resource within OpenTelekomCloud.

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

resource "opentelekomcloud_lb_rule_v3" "this" {
  type         = "PATH"
  compare_type = "REGEX"
  value        = "^.+$"
  policy_id    = opentelekomcloud_lb_policy_v3.this.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the Policy. Only administrative users can specify a tenant UUID other than
  their own. Changing this creates a new Policy.

* `policy_id` - (Required) ID of the policy.

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

* `rule_id` - ID of the policy rule.

## Import

Rules can be imported using the `policy_id/rule_id`, e.g.

```shell
terraform import opentelekomcloud_lb_rule_v3.this 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74/1bb93b8b-37a4-4b50-92cc-daa4c89d4e4c
```
