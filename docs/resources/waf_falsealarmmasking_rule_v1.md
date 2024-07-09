---
subcategory: "Web Application Firewall (WAF)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_falsealarmmasking_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-falsealarmmasking-rule-v1"
description: |-
Manages a WAF False Alarm Masking Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF false alarm masking rules you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/false_alarm_masking_rules)

# opentelekomcloud_waf_falsealarmmasking_rule_v1

Manages a WAF False Alarm Masking Rule resource within OpenTelekomCloud.

!>
This resource is known to be broken due to the API changes and will be fixed in the upcoming releases

## Example Usage

```hcl
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_waf_falsealarmmasking_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
  url       = "/a"
  rule      = "100001"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The WAF policy ID. Changing this creates a new rule.

* `url` - (Required) Specifies a misreported URL excluding a domain name. Changing this creates a new rule.

* `rule` - (Required) Specifies the rule ID, which consists of six digits and cannot be empty. Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the rule.

## Import

False Alarm Masking Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_falsealarmmasking_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
