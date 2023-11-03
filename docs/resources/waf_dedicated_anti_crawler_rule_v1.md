---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated Anti Crawler rule you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_javascript_anti-crawler_rule.html).

# opentelekomcloud_waf_dedicated_anti_crawler_rule_v1

Manages a WAF Dedicated Anti Crawler Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_cc"
}

resource "opentelekomcloud_waf_dedicated_anti_crawler_rule_v1" "rule_1" {
  policy_id       = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  name            = "anticrawler_1"
  url             = "/patent/id"
  logic           = 3
  protection_mode = "anticrawler_except_url"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `url` - (Required, String) URL to which the rule applies.

* `logic` - (Required, Int) Rule matching logic.
  Values are:
  + 1: Include
  + 2: Not include
  + 3: Equal
  + 4: Not equal
  + 5: Prefix is
  + 6: Prefix is not
  + 7: Suffix is
  + 8: Suffix is not

* `name` - (Required, String) Rule name.

* `protection_mode` - (Required, ForceNew, String) JavaScript anti-crawler rule type.
  Values are:
  + `anticrawler_specific_url`: used to protect a specific path specified by the rule.
  + `anticrawler_except_url`: used to protect all paths except the one specified by the rule
  Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status. The value can be `0` or `1`.
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Anti Crawler Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_anti_crawler_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
