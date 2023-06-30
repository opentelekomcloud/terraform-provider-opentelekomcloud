---
subcategory: "Web Application Firewall (WAF)"
---

Up-to-date reference of API arguments for WAF policy you can get at
`https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/policies`.

# opentelekomcloud_waf_policy_v1

Manages a WAF policy resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_policy_v1" "policy_1" {
  name = "policy_1"
  options {
    webattack = true
    crawler   = true
  }
  full_detection = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The policy name. The maximum length is 256 characters. Only digits, letters, underscores(_), and hyphens(-) are allowed.

* `action` - (Optional) Specifies the protective action after a rule is matched. The action object structure is documented below.

* `options` - (Optional) Specifies the protection switches. The options object structure is documented below.

* `level` - (Optional) Specifies the protection level.
  * `1`: low
  * `2`: medium
  * `3`: high

* `full_detection` - (Optional) Specifies the detection mode in Precise Protection.
  * `true`: full detection, Full detection finishes all threat detections before blocking requests that meet Precise Protection specified conditions.
  * `false`: instant detection. Instant detection immediately ends threat detection after blocking a request that meets Precise Protection specified conditions.

* `hosts` - (Optional) An array of the domain IDs.

The `action` block supports:

* `category` - (Required) Specifies the protective action.
  * `block`: WAF blocks and logs detected attacks.
  * `log`: WAF logs detected attacks only.

The `options` block supports:

* `webattack` - (Optional) Specifies whether Basic Web Protection is enabled.

* `common` - (Optional) Specifies whether General Check in Basic Web Protection is enabled.

* `crawler` - (Optional) Specifies whether the master crawler detection switch in Basic Web Protection is enabled.

* `crawler_engine` - (Optional) Specifies whether the Search Engine switch in Basic Web Protection is enabled.

* `crawler_scanner` - (Optional) Specifies whether the Scanner switch in Basic Web Protection is enabled.

* `crawler_script` - (Optional) Specifies whether the Script Tool switch in Basic Web Protection is enabled.

* `crawler_other` - (Optional) Specifies whether detection of other crawlers in Basic Web Protection is enabled.

* `webshell` - (Optional) Specifies whether webshell detection in Basic Web Protection is enabled.

* `cc` - (Optional) Specifies whether CC Attack Protection is enabled.

* `custom` - (Optional) Specifies whether Precise Protection is enabled.

* `whiteblackip` - (Optional) Specifies whether Blacklist and Whitelist is enabled.

* `privacy` - (Optional) Specifies whether Data Masking is enabled.

* `ignore` - (Optional) Specifies whether False Alarm Masking is enabled.

* `antitamper` - (Optional) Specifies whether Web Tamper Protection is enabled.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the policy.

* `hosts` - Specifies the domain IDs.

## Import

Policies can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_waf_policy_v1.policy_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
