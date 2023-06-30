---
subcategory: "Web Application Firewall (WAF)"
---

Up-to-date reference of API arguments for WAF CC attack protection rule you can get at
`https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/cc_attack_protection_rules`.

# opentelekomcloud_waf_ccattackprotection_rule_v1

Manages a WAF CC Attack Protection Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_waf_ccattackprotection_rule_v1" "rule_1" {
  policy_id    = opentelekomcloud_waf_policy_v1.policy_1.id
  url          = "/abc1"
  limit_num    = 10
  limit_period = 60
  lock_time    = 10
  tag_type     = "cookie"
  tag_index    = "sessionid"

  action_category    = "block"
  block_content_type = "application/json"
  block_content      = "{\"error\":\"forbidden\"}"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The WAF policy ID. Changing this creates a new rule.

* `url` - (Required) Specifies a misreported URL excluding a domain name. Changing this creates a new rule.

* `limit_num` - (Required) Specifies the number of requests allowed from a web visitor in a rate limiting period. Changing this creates a new rule.

* `limit_period` - (Required) Specifies the rate limiting period. Changing this creates a new rule.

* `lock_time` - (Optional) Specifies the lock duration. The value ranges from 0 seconds to 2^32 seconds. Changing this creates a new rule.

* `tag_type` - (Required) Specifies the rate limit mode. Changing this creates a new rule. Valid Options are:
  * `ip` - A web visitor is identified by the IP address.
  * `cookie` - A web visitor is identified by the cookie key value.
  * `other` - A web visitor is identified by the Referer field(user-defined request source).

* `tag_index` - (Optional) If `tag_type` is set to `cookie`, this parameter indicates cookie name. Changing this creates a new rule.

* `tag_category` - (Optional) Specifies the category. The value is `referer`. Changing this creates a new rule.

* `tag_contents` - (Optional) Specifies the category content. Changing this creates a new rule.

* `action_category` - (Required) Specifies the action. Changing this creates a new rule. Valid Options are:
  * `block` - block the requests.
  * `captcha` - Verification code. The user needs to enter the correct verification code after blocking to restore the correct access page.

* `block_content_type` - (Optional) Specifies the type of the returned page. The options are `application/json`, `text/html`, and `text/xml`. Changing this creates a new rule.

* `block_content` - (Optional) Specifies the content of the returned page. Changing this creates a new rule.


## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `default` - Specifies whether the rule is the default CC attack protection rule.

## Import

CC Attack Protection Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_ccattackprotection_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
