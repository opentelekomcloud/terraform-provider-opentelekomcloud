---
subcategory: "Dedicated Web Application Firewall (WAFD)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_dedicated_alarm_masking_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-dedicated-alarm-masking-rule-v1"
description: |-
  Manages a WAF Dedicated False Alarm Masking Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF dedicated Global Protection Whitelist (formerly False Alarm Masking) rule you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_global_protection_whitelist_formerly_false_alarm_masking_rule.html).

# opentelekomcloud_waf_dedicated_alarm_masking_rule_v1

Manages a WAF Dedicated Global Protection Whitelist (formerly False Alarm Masking) Rule resource within OpenTelekomCloud.

## Example Usage

### Basic example

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_am"
}

resource "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  domains     = ["www.example.com"]
  rule        = "xss"
  description = "description"

  conditions {
    category        = "url"
    contents        = ["/login"]
    logic_operation = "equal"
  }
}
```

### Advanced settings with empty contents

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_am"
}

resource "opentelekomcloud_waf_dedicated_alarm_masking_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  domains     = ["www.example.com"]
  rule        = "all"
  description = "description"

  conditions {
    category        = "url"
    contents        = ["/login"]
    logic_operation = "equal"
  }
  advanced_settings {
    index = "cookie"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `domains` - (Required, ForceNew, List) Domain names to be protected. Changing this creates a new rule.

* `conditions` - (Optional, ForceNew, List) Condition list. Changing this creates a new rule.
  The `conditions` block supports:

  + `category` - (Required, ForceNew, String) Field type. The value can be `url`, `ip`, `params`, `cookie`, or `header`.

  + `logic_operation` - (Required, ForceNew, String) The matching logic varies depending on the field type.
    + if the field type is `ip`, the logic can be `equal` or `not_equal`.
    + If the field type is `url`, `params`, `cookie`, or `header`, the logic can be `equal`, `not_equal`, `contain`, `not_contain`, `prefix`, `not_prefix`, `suffix`, `not_suffix`.

  + `contents` - (Optional, ForceNew, List) Content. The array length is limited to 1.
    The content format varies depending on the field type.
    + For example, if the field type is `ip`, the value must be an `IP address` or `IP address range`.
    + If the field type is `url`, the value must be in the `standard URL format`.
    + IF the field type is `params`, `cookie`, or `header`, the content format is not limited.

  + `index` - (Optional, ForceNew, String) Subfield. When `category` is set to `params`, `cookie`, or `header`, set this parameter based on site requirements. This parameter is mandatory.

* `advanced_settings` - (Optional, ForceNew, List) To ignore attacks of a specific field, specify the field in the Advanced settings area.
  After you add the rule, WAF will stop blocking attacks of the specified field.
  This parameter is not included if all modules are bypassed. Changing this creates a new rule.
  The `advanced_settings` block supports:
  + `contents` - (Optional, ForceNew, List) Subfield of the specified field type. The default value is all.

  + `index` - (Optional, ForceNew, String) Field type.
    The following field types are supported: `Params`, `Cookie`, `Header`, `Body`, and `Multipart`.
    When you select `Params`, `Cookie`, or `Header`, you can set this parameter to `all` or configure subfields as required.

* `rule` - (Required, ForceNew, String) Items to be masked. Changing this creates a new rule.
  You can provide multiple items and separate them with semicolons (;).
  + If you want to disable a specific built-in rule for a domain name, the value of this parameter is the rule ID.
    When requests are blocked against a certain built-in rule while you do not want this rule to block requests later,
    you can query the rule in the Events page on the console and find its rule ID in the Hit Rule column.
    Then, you can disk the rule by its ID (including 6 digits).
  + If you want to mask a type of basic web protection rules, set this parameter to the name of the type of basic
    web protection rules.
    + `xss`: XSS attacks
    + `webshell`: Web shells
    + `vuln`: Other types of attacks
    + `sqli`: SQL injection attack
    + `robot`: Malicious crawlers
    + `rfi`: Remote file inclusion
    + `lfi`: Local file inclusion
    + `cmdi`: Command injection attack
  + To bypass the basic web protection, set this parameter to `all`.
  + To bypass all WAF protection, set this parameter to `bypass`.

* `description` - (Optional, ForceNew, String) Rule description. Changing this creates a new rule.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status. The value can be:
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Web Global Protection Whitelist (formerly False Alarm Masking) rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_alarm_masking_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
