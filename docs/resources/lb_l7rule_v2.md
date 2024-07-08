---
subcategory: "Elastic Load Balancer (ELB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_l7rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-lb-l7rule-v2"
description: |-
Manages a ELB L7 Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for ELB l7 rule you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v2.0/forwarding_rule)

# opentelekomcloud_lb_l7rule_v2

Manages a V2 L7 Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "SUBNET_ID"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "opentelekomcloud_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  = opentelekomcloud_lb_listener_v2.listener_1.id
  redirect_url = "http://www.example.com"
}

resource "opentelekomcloud_lb_l7rule_v2" "l7rule_1" {
  l7policy_id  = opentelekomcloud_lb_l7policy_v2.l7policy_1.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  If omitted, the `region` argument of the provider is used.
  Changing this creates a new L7 Rule.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the L7 Rule.  Only administrative users can specify a tenant UUID
  other than their own. Changing this creates a new L7 Rule.

* `description` - (Optional) Human-readable description for the L7 Rule.

* `type` - (Required) The L7 Rule type - can either be HOST_NAME or PATH. Changing this creates a new L7 Rule.

* `compare_type` - (Required) The comparison type for the L7 rule - can either be
  STARTS_WITH, EQUAL_TO or REGEX

* `l7policy_id` - (Required) The ID of the L7 Policy to query. Changing this creates a new
  L7 Rule.

* `value` - (Required) The value to use for the comparison. For example, the file type to
  compare.

* `key` - (Optional) The key to use for the comparison. For example, the name of the cookie to
  evaluate. Valid when `type` is set to COOKIE or HEADER. Changing this creates a new L7 Rule.

* `admin_state_up` - (Optional) The administrative state of the L7 Rule.
  The value can only be true (UP).

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the L7 Rule.

* `region` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `type` - See Argument Reference above.

* `compare_type` - See Argument Reference above.

* `l7policy_id` - See Argument Reference above.

* `value` - See Argument Reference above.

* `key` - See Argument Reference above.

* `invert` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `listener_id` - The ID of the Listener owning this resource.

## Import

Load Balancer L7 Rule can be imported using the L7 Policy ID and L7 Rule ID
separated by a slash, e.g.:

```sh
terraform import opentelekomcloud_lb_l7rule_v2.l7rule_1 e0bd694a-abbe-450e-b329-0931fd1cc5eb/4086b0c9-b18c-4d1c-b6b8-4c56c3ad2a9e
```
