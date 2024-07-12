---
subcategory: "Web Application Firewall (WAF)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_datamasking_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-datamasking-rule-v1"
description: |-
  Manages a WAF Datamasking rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF datamasking rule you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/data_masking_rules)

# opentelekomcloud_waf_datamasking_rule_v1

Manages a WAF Data Masking Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_waf_datamasking_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
  url       = "/login"
  category  = "params"
  index     = "password"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The WAF policy ID. Changing this creates a new rule.

* `url` - (Required) Specifies the URL to which the data masking rule applies.

* `category` - (Required) Specifies the masked field. The options are params and header.

* `index` - (Required) Specifies the masked subfield.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the rule.

## Import

Data Masking Rules can be imported using the `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_datamasking_rule_v1.rule_1 b39f3a5a1b4f447a8030f0b0703f47f5/7117d38e4c8f4624a505bd96b97d024c
```
