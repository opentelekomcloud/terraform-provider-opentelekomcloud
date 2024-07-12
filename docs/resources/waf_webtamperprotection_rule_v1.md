---
subcategory: "Web Application Firewall (WAF)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_webtamperprotection_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-webtamperprotection-rule-v1"
description: |-
  Manages a WAF Web Tamper Protection Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF web tamper protection rule you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/web_tamper_protection_rules)

# opentelekomcloud_waf_webtamperprotection_rule_v1

Manages a WAF Web Tamper Protection Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_waf_webtamperprotection_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
  hostname  = "www.abc.com"
  url       = "/a"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The WAF policy ID. Changing this creates a new rule.

* `hostname` - (Required) Specifies the domain name. Changing this creates a new rule.

* `url` - (Required) Specifies the URL protected by the web tamper protection rule, excluding a domain name. Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the rule.

## Import

Web Tamper Protection Rules can be imported using the `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_webtamperprotection_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/7117d38e4c8f4624a505-bd96b97d024c
```
