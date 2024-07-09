---
subcategory: "Dedicated Web Application Firewall (WAFD)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_dedicated_anti_leakage_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-dedicated-anti-leakage-rule-v1"
description: |-
Manages a WAF Dedicated Anti Leakage Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF dedicated Information Leakage Protection rule you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_an_information_leakage_protection_rule.html).

# opentelekomcloud_waf_dedicated_anti_leakage_rule_v1

Manages a WAF Dedicated Information Leakage Protection Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_al"
}

resource "opentelekomcloud_waf_dedicated_anti_leakage_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  url         = "/attack"
  category    = "sensitive"
  contents    = ["id_card"]
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `category` - (Required, String) Sensitive information type in the information leakage prevention rule.
  Values:
  + `sensitive`: The rule masks sensitive user information, such as ID code, phone numbers, and email addresses.
  + `code`: The rule blocks response pages of specified HTTP response code.

* `url` - (Required, String) URL to which the rule applies, for example, `/admin`

* `description` - (Optional, String) Rule description.

* `contents` - (Optional, List) Content corresponding to the sensitive information type.
  Multiple options can be set.
  + When category is set to `code`, the pages that contain the following HTTP response codes will be blocked: `400`, `401`, `402`, `403`, `404`, `405`, `500`, `501`, `502`, `503`, `504` and `507`.
  + When category is set to `sensitive`, parameters `phone`, `id_card`, and `email` can be set.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status. The value can be:
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Web Information Leakage Protection rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_anti_leakage_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
