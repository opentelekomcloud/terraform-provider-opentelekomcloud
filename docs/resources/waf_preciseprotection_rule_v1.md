---
subcategory: "Web Application Firewall (WAF)"
---

Up-to-date reference of API arguments for WAF precise protection rule you can get at
`https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/precise_protection_rules`.

# opentelekomcloud_waf_preciseprotection_rule_v1

Manages a WAF Precise Protection Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_waf_preciseprotection_rule_v1" "rule_1" {
  policy_id = opentelekomcloud_waf_policy_v1.policy_1.id
  name      = "rule_1"

  conditions {
    category = "path"
    contents = ["/login"]
    logic    = "contain"
  }
  conditions {
    category = "ip"
    contents = ["192.168.1.1"]
    logic    = "equal"
  }

  action_category = "block"
  priority        = 10
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The WAF policy ID. Changing this creates a new rule.

* `name` - (Required) Specifies the name of a precise protection rule. Changing this creates a new rule.

* `time` - (Optional) Specifies the effect time of the precise protection rule. Changing this creates a new rule.
  * `false` - The rule takes effect immediately.
  * `true` - The rule takes effect at the scheduled time.

* `start` - (Optional) Specifies the time when the precise protection rule takes effect. If time is set to true,
  either the start time or the end time must be set. Changing this creates a new rule.

* `end` - (Optional) Specifies the time when the precise protection rule expires. If time is set to true,
  either the start time or the end time must be set. Changing this creates a new rule.

* `conditions` - (Required) Specifies the condition parameters. Changing this creates a new rule.
  The conditions object structure is documented below.

* `action` - (Required) Specifies the protective action after the precise protection rule is matched.
  Changing this creates a new rule. The action object structure is documented below.

* `priority` - (Optional) Specifies the priority of a rule being executed. Smaller values correspond to higher priorities.
  If two rules are assigned with the same priority, the rule added earlier has higher priority, the rule added earlier
  has higher priority. The value ranges from 0 to 65535. Changing this creates a new rule.

The `conditions` block supports:

* `category` - (Required) Specifies the condition type. The value can be path, user-agent, ip, params, cookie, referer, or header.

* `index` - (Optional) If `category` is set to cookie, index indicates cookie name, if set to params, index indicates param name,
  if set to header, index indicates an option in the header.

* `logic` - (Required) contain, not_contain, equal, not_equal, prefix, not_prefix, suffix, and not_suffix indicate Include, Exclude, Equal to, Not equal to, Prefix is, Prefix is not, Suffix is, and Suffix is not respectively. If category is set to ip, logic can only be equal or not_equal.

* `contents` - (Required) Specifies a list of content matching the condition. Currently, only one value is accepted.

The `action` block supports:

* `category` - (Required) Specifies the protective action. The value can be block or pass.


## Attributes Reference

The following attributes are exported:

* `id` - ID of the rule.

## Import

Precise Protection Rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_preciseprotection_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
