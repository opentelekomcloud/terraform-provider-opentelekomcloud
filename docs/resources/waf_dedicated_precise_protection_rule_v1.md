---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated Precise Protection rule you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_precise_protection_rule.html).

# opentelekomcloud_waf_dedicated_precise_protection_rule_v1

Manages a WAF Dedicated Precise Protection Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_pp"
}

resource "opentelekomcloud_waf_dedicated_precise_protection_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  time        = false
  description = "desc"
  priority    = 50

  conditions {
    category        = "url"
    contents        = ["test"]
    logic_operation = "contain"
  }
  action {
    category = "block"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `time` - (Required, ForceNew, Bool) Time the precise protection rule takes effect. Changing this creates a new rule.
  Values:
  + `false`: The rule takes effect immediately.
  + `true`: The effective time is customized.

* `start` - (Required, ForceNew, Int) Timestamp (ms) when the precise protection rule takes effect. This parameter is returned only when time is true. Changing this creates a new rule.

* `terminal` - (Required, ForceNew, Int) Timestamp (ms) when the precise protection rule expires. This parameter is returned only when time is true. Changing this creates a new rule.

* `description` - (Optional, ForceNew, String) Rule description. Changing this creates a new rule.

* `conditions` - (Optional, ForceNew, List) Match condition List. Changing this creates a new rule.
  The `conditions` block supports:

  + `category` - (Optional, ForceNew, String) Field type. The options are `url`, `user-agent`, `ip`, `params`, `cookie`, `referer`, `header`, `request_line`, `method`, and `request`.

  + `logic_operation` - (Optional, ForceNew, String) Logic for matching the condition. Changing this creates a new rule.
    + If the category is `url`, `user-agent` or `referer` , the optional operations are `contain`, `not_contain`, `equal`, `not_equal`, `prefix`, `not_prefix`, `suffix`, `not_suffix`, `contain_any`, `not_contain_all`, `equal_any`, `not_equal_all`, `equal_any`, `not_equal_all`, `prefix_any`, `not_prefix_all`, `suffix_any`, `not_suffix_all`, `len_greater`, `len_less`, `len_equal` and `len_not_equal`
    + If the category is `ip`, the optional operations are: `equal`, `not_equal`, `equal_any` and `not_equal_all`
    + If the category is `method`, the optional operations are: `equal` and `not_equal`
    + If the category is `request_line` and `request`, the optional operations are: `len_greater`, `len_less`, `len_equal` and `len_not_equal`
    + If the category is `params`, `header`, and `cookie`, the optional operations are: `contain`, `not_contain`, `equal`, `not_equal`, `prefix`, `not_prefix`, `suffix`, `not_suffix`, `contain_any`, `not_contain_all`, `equal_any`, `not_equal_all`, `equal_any`, `not_equal_all`, `prefix_any`, `not_prefix_all`, `suffix_any`, `not_suffix_all`, `len_greater`, `len_less`, `len_equal`, `len_not_equal`, `num_greater`, `num_less`, `num_equal`, `num_not_equal`, `exist` and `not_exist`

  + `contents` - (Optional, ForceNew, List) Content of the conditions. This parameter is mandatory when the suffix of `logic_operation` is not `any` or `all`. This parameter is mandatory when the suffix of `logic_operation` is not `any` or `all`. Changing this creates a new rule.

  + `value_list_id` - (Optional, ForceNew, String) Reference table ID. This parameter is mandatory when the suffix of `logic_operation` is `any` or `all`. The reference table type must be the same as the category type. Changing this creates a new rule.

  + `index` - (Optional, ForceNew, String) Subfield. Changing this creates a new rule.
    + When the field type is `url`, `user-agent`, `ip`, `refer`, `request_line`, `method`, or `request`, index is not required.
    + When the field type is `params`, `header`, or `cookie`, and the subfield is customized, the value of index is the customized subfield.

* `action` - (Required, ForceNew, Set) Protection action to take if the number of requests reaches the upper limit. Changing this creates a new rule.
  The `conditions` block supports:

  + `category` - (Required, ForceNew, String) Action type. Changing this creates a new rule.
    + `block`: WAF blocks attacks.
    + `pass`: WAF allows requests.
    + `log`: WAF only logs detected attacks.

  + `followed_action_id` - (Optional, ForceNew, String) ID of a known attack source rule. This parameter can be configured only when category is set to block. Changing this creates a new rule.

* `priority` - (Optional, ForceNew, Int) Priority of a rule. A small value indicates a high priority. If two rules are assigned with the same priority, the rule added earlier has higher priority. Value range: `0` to `1000`. Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status. The value can be:
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Precise Protection Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_precise_protection_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
