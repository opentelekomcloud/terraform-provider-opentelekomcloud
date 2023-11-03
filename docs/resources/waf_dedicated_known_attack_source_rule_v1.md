---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated Known Attack Source rule you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_known_attack_source_rule.html).

# opentelekomcloud_waf_dedicated_known_attack_source_rule_v1

Manages a WAF Dedicated Known Attack Source Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_ka"
}

resource "opentelekomcloud_waf_dedicated_known_attack_source_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  block_time  = 300
  category    = "long_cookie_block"
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `block_time` - (Required, Int) Block duration, in seconds.
  If prefix long is selected for the rule type, the value for `block_time` ranges from `301` to `1800`.
  If prefix short is selected for the rule type, the value for `block_time` ranges from `0` to `300`.

* `category` - (Required, ForceNew, String) Type of the know attack source rule.
  Enumeration values:
    + `long_ip_block`
    + `long_cookie_block`
    + `long_params_block`
    + `short_ip_block`
    + `short_cookie_block`
    + `short_params_block`

* `description` - (Optional, String) Rule description.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Known Attack Source Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_known_attack_source_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
