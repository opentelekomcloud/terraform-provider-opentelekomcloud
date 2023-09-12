---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated Web Tamper rule you can get at
`https://docs-beta.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_web_tamper_protection_rule.html`.

# opentelekomcloud_waf_dedicated_web_tamper_rule_v1

Manages a WAF Dedicated Web Tamper Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_at"
}

resource "opentelekomcloud_waf_dedicated_web_tamper_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  hostname    = "www.domain.com"
  url         = "/login"
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `hostname` - (Required, ForceNew, String) Protected website.

* `url` - (Required, ForceNew, String) URL protected by the web tamper protection rule.
  The value must be in the standard URL format, for example, `/admin`

* `description` - (Optional, ForceNew, String) Rule description.

* `update_cache` - (Optional, Bool) To update the cache for a web tamper protection Rule.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status. The value can be:
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Web Tamper Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_web_tamper_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
