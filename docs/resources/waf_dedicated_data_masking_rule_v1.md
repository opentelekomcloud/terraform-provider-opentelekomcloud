---
subcategory: "Dedicated Web Application Firewall (WAFD)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_dedicated_data_masking_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-dedicated-data-masking-rule-v1"
description: |-
Manages a WAF Dedicated Data Masking Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF dedicated Data Masking rule you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_data_masking_rule.html).

# opentelekomcloud_waf_dedicated_data_masking_rule_v1

Manages a WAF Dedicated Data Masking Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_dm"
}

resource "opentelekomcloud_waf_dedicated_data_masking_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name        = "data_masking"
  url         = "/login"
  category    = "params"
  description = "description"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `url` - (Required, String) URL protected by the data masking rule. The value must be in the standard URL format, for example, /admin.

* `name` - (Required, String) Name of the masked field.

* `category` - (Required, String) Masked field.
  Values:
  + `params`
  + `cookie`
  + `header`
  + `form`

* `description` - (Optional, String) Rule description.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` -  Rule status. The value can be `0` or `1`.
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Data Masking Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_data_masking_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
