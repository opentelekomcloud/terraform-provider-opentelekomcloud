---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_whiteblackip_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-whiteblackip-rule-v1"
description: |-
  Manages a V1 WAF WhiteBlackIP Rule resource within OpenTelekomCloud.
---

# opentelekomcloud_waf_whiteblackip_rule_v1

Manages a WAF WhiteBlackIP Rule resource within OpenTelekomCloud.

## Example Usage

```hcl

resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_whiteblackip_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	addr = "192.168.0.0/24"
	white = 1
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

WhiteBlackIP Rules can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_waf_whiteblackip_rule_v1.rule_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
