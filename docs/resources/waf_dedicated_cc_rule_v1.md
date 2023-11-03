---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated CC rule you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_cc_attack_protection_rule.html).

# opentelekomcloud_waf_dedicated_cc_rule_v1

Manages a WAF Dedicated CC Attack Protection Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_cc"
}

resource "opentelekomcloud_waf_dedicated_cc_rule_v1" "rule_1" {
  policy_id    = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  mode         = 0
  url          = "/abc1"
  limit_num    = 10
  limit_period = 60
  lock_time    = 10
  tag_type     = "cookie"
  tag_index    = "sessionid"

  action {
    category     = "block"
    content_type = "application/json"
    content      = "{\"error\":\"forbidden\"}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `mode` - (Required, ForceNew, Int) Protection mode of the CC attack protection rule. Changing this creates a new rule. Valid Options are:
  * `0` - Standard. Only the protected paths of domain names can be specified.
  * `1` - The path, IP address, cookie, header, and params fields can all be set.

* `url` - (Required, ForceNew, String) Path to be protected in the CC attack protection rule. Changing this creates a new rule.

* `conditions` - (Optional, ForceNew, List) Rate limit conditions of the CC protection rule. Changing this creates a new rule.
    The `conditions` block supports:

  + `category` - (Required, ForceNew, String) Field type. The value can be `url`, `ip`, `params`, `cookie`, or `header`.

  + `logic_operation` - (Required, ForceNew, String) Logic for matching the condition.
    + If the category is `url`, the optional operations are `contain`, `not_contain`, `equal`, `not_equal`, `prefix`, `not_prefix`, `suffix`, `not_suffix`, `contain_any`, `not_contain_all`, `equal_any`, `not_equal_all`, `equal_any`, `not_equal_all`, `prefix_any`, `not_prefix_all`, `suffix_any`, `not_suffix_all`, `len_greater`, `len_less`, `len_equal` and `len_not_equal`
    + If the category is `ip`, the optional operations are: `equal`, `not_equal`, `equal_any` and `not_equal_all`
    + If the category is `params`, `cookie` and `header`, the optional operations are: `contain`, `not_contain`, `equal`, `not_equal`, `prefix`, `not_prefix`, `suffix`, `not_suffix`, `contain_any`, `not_contain_all`, `equal_any`, `not_equal_all`, `equal_any`, `not_equal_all`, `prefix_any`, `not_prefix_all`, `suffix_any`, `not_suffix_all`, `len_greater`, `len_less`, `len_equal`, `len_not_equal`, `num_greater`, `num_less`, `num_equal`, `num_not_equal`, `exist` and `not_exist`
      Changing this creates a new rule.

  + `contents` - (Optional, ForceNew, List) Content of the conditions. This parameter is mandatory when the suffix of `logic_operation` is not `any` or `all`. Changing this creates a new rule.

  + `value_list_id` - (Optional, ForceNew, String) Reference table ID. This parameter is mandatory when the suffix of `logic_operation` is `any` or `all`. The reference table type must be the same as the category type. Changing this creates a new rule.

  + `index` - (Optional, ForceNew, String) Subfield. When `category` is set to `params`, `cookie`, or `header`, set this parameter based on site requirements. This parameter is mandatory. Changing this creates a new rule.

* `action` - (Required, ForceNew, Set) Protection action to take if the number of requests reaches the upper limit. Changing this creates a new rule.
  The `conditions` block supports:

  + `category` - (Required, ForceNew, String) Action type. Changing this creates a new rule.
    + `captcha`: Verification code. WAF requires visitors to enter a correct verification code to continue their access to requested page on your website.
    + `block`: WAF blocks the requests. When tag_type is set to other, the value can only be blocked.
    + `log`: WAF logs the event only.
    + `dynamic_block`: In the previous rate limit period, if the request frequency exceeds the value of Rate Limit Frequency, the request is blocked. In the next rate limit period, if the request frequency exceeds the value of Permit Frequency, the request is still blocked.

      -> **Note:**: The `dynamic_block` protection action can be set only when the advanced protection mode is enabled for the CC protection rule.

  + `content_type` - (Optional, ForceNew, String) User identifier. The value is fixed at referer. Changing this creates a new rule.

  + `content` - (Optional, ForceNew, String) Protection page content. Changing this creates a new rule.

* `tag_type` - (Required, ForceNew, String) Rate limit mode. Changing this creates a new rule. Valid Options are:
  * `ip` - IP-based rate limiting. Website visitors are identified by IP address.
  * `cookie` - User-based rate limiting. Website visitors are identified by the cookie key value.
  * `header` - User-based rate limiting. Website visitors are identified by the header field.
  * `other` - Website visitors are identified by the Referer field (user-defined request source).

* `tag_index` - (Optional, ForceNew, String) User identifier. Changing this creates a new rule.
  If `tag_type` is set to `cookie`, this parameter indicates cookie name.
  If `tag_type` is set to `header`, this parameter indicates header name.

* `tag_category` - (Optional, ForceNew, String) Specifies the category. The value is `referer`. Changing this creates a new rule.

* `tag_contents` - (Optional, ForceNew, String) Specifies the category content. Changing this creates a new rule.

* `limit_num` - (Required, ForceNew, Int) Rate limit frequency based on the number of requests. The value ranges from `1` to `2,147,483,647`. Changing this creates a new rule.

* `limit_period` - (Required, ForceNew, Int) Rate limit period, in seconds. The value ranges from `1` to `3,600`. Changing this creates a new rule.

* `unlock_num` - (Optional, ForceNew, Int) Allowable frequency based on the number of requests. The value ranges from `0` to `2,147,483,647`. This parameter is required only when the protection `action` type is `dynamic_block`. Changing this creates a new rule.

* `lock_time` - (Optional, ForceNew, String) Block duration, in seconds. The value ranges from `0` to `65,535`. Specifies the period within which access is blocked. An error page is displayed in this period. Changing this creates a new rule.

* `description` - (Optional, ForceNew, String) Rule description. Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF CC Attack Protection Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_cc_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
