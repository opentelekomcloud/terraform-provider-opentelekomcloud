---
subcategory: "Web Application Firewall (WAF)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_whiteblackip_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-whiteblackip-rule-v1"
description: |-
Manages a WAF White and Black IP Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF white and black ip rule you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/blacklist_and_whitelist_rules)

# opentelekomcloud_waf_whiteblackip_rule_v1

Manages a WAF WhiteBlackIP Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_waf_whiteblackip_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
  addr      = "192.168.0.0/24"
  white     = 1
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The WAF policy ID. Changing this creates a new rule.

* `addr` - (Required) Specifies the IP address or range. For example, 192.168.0.125 or 192.168.0.0/24.

* `white` - (Optional) Specifies the IP address type. 1: Whitelist, 0: Blacklist. If you do not configure
  the white parameter, the value is Blacklist by default.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

## Import

WhiteBlackIP Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_whiteblackip_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
