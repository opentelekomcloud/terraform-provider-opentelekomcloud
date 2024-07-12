---
subcategory: "Dedicated Web Application Firewall (WAFD)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_dedicated_blacklist_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-dedicated-blacklist-rule-v1"
description: |-
  Manages a WAF Dedicated Blacklist Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF dedicated Blacklist rule you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_blacklist_or_whitelist_rule.html).

# opentelekomcloud_waf_dedicated_blacklist_rule_v1

Manages a WAF Dedicated Blacklist Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_black"
}

resource "opentelekomcloud_waf_dedicated_blacklist_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name        = "my_blacklist"
  ip_address  = "192.168.1.0/24"
  action      = 0
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule. Changing this creates a new rule.

* `name` - (Required, ForceNew, String) Rule name. Changing this creates a new rule.

* `ip_address` - (Required, ForceNew, String) IP addresses or an IP address range to be added to the blacklist or whitelist. Changing this creates a new rule.
  For example, `192.x.x.3` or `10.x.x.0/24`

* `action` - (Required, ForceNew, Int) Protective action. Changing this creates a new rule.
  The value can be:
    + `0`: WAF blocks the requests that hit the rule.
    + `1`: WAF allows the requests that hit the rule.
    + `2`: WAF only logs the requests that hit the rule.

* `followed_action_id` - (Optional, ForceNew, String) ID of a known attack source rule. Changing this creates a new rule.
  This parameter can be configured only when `action` is set to `0`.

* `description` - (Optional, ForceNew, String) Rule description. Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status. The value can be:
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Blacklist Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_blacklist_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
