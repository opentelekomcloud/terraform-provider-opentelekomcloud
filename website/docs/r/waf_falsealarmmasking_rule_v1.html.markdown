---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_falsealarmmasking_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-falsealarmmasking-rule-v1"
description: |-
  Manages a V1 WAF False Alarm Masking Rule resource within OpenTelekomCloud.
---

# opentelekomcloud_waf_falsealarmmasking_rule_v1

Manages a WAF False Alarm Masking Rule resource within OpenTelekomCloud.

## Example Usage

```hcl

resource "opentelekomcloud_waf_policy_v1" "policy_1" {
	name = "policy_1"
}

resource "opentelekomcloud_waf_falsealarmmasking_rule_v1" "rule_1" {
	policy_id = "${opentelekomcloud_waf_policy_v1.policy_1.id}"
	url = "/a"
	rule = "100001"
}

```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The WAF policy ID. Changing this creates a new rule.

* `url` - (Required) Specifies a misreported URL excluding a domain name. Changing this creates a new rule.

* `rule` - (Required) Specifies the rule ID, which consists of six digits and cannot be empty. Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

## Import

False Alarm Masking Rules can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_waf_falsealarmmasking_rule_v1.rule_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
