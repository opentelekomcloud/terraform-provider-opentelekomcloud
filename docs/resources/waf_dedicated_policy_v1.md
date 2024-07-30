---
subcategory: "Dedicated Web Application Firewall (WAFD)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_dedicated_policy_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-dedicated-policy-v1"
description: |-
  Manages a WAF Dedicated Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF policy you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/policy_management/index.html).

# opentelekomcloud_waf_dedicated_policy_v1

Manages a WAF dedicated policy resource within OpenTelekomCloud.

-> **Note:** For this resource region must be set in environment variable `OS_REGION_NAME` or in `clouds.yaml`


## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name            = "policy_1"
  level           = 3
  protection_mode = "block"
  full_detection  = true

  options {
    crawler    = false
    web_attack = false
    cc         = true
    web_shell  = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The policy name.

* `protection_mode` - (Optional, String) Specifies the protective action after a rule is matched.
  Values are:
  + `block`: WAF blocks and logs detected attacks.
  + `log`: WAF logs detected attacks only.

* `level` - (Optional, Int) Specifies the protection level.
  Values are:
  + `1`: low
  + `2`: medium
  + `3`: high

* `options` - (Optional, List)  Specifies the protection switches.
  The `options` block supports:
  + `web_attack` - (Optional, Bool) Specifies whether Basic Web Protection is enabled.
  + `common` - (Optional, Bool) Specifies whether General Check in Basic Web Protection is enabled.
  + `crawler` - (Optional, Bool) Specifies whether the master crawler detection switch in Basic Web Protection is enabled.
  + `anti_crawler` - (Optional, Bool) JavaScript anti-crawler function.
  + `crawler_engine` - (Optional, Bool) Specifies whether the Search Engine switch in Basic Web Protection is enabled.
  + `crawler_scanner` - (Optional, Bool) Specifies whether the Scanner switch in Basic Web Protection is enabled.
  + `crawler_script` - (Optional, Bool) Specifies whether the Script Tool switch in Basic Web Protection is enabled.
  + `crawler_other` - (Optional, Bool) Specifies whether detection of other crawlers in Basic Web Protection is enabled.
  + `web_shell` - (Optional, Bool) Specifies whether webshell detection in Basic Web Protection is enabled.
  + `cc` - (Optional, Bool) Specifies whether CC Attack Protection is enabled.
  + `custom` - (Optional, Bool) Specifies whether Precise Protection is enabled.
  + `blacklist` - (Optional, Bool) Specifies whether Blacklist and Whitelist is enabled.
  + `geolocation_access_control` - (Optional, Bool) Whether geolocation access control is enabled.
  + `ignore` - (Optional, Bool) Whether false alarm masking is enabled.
  + `privacy` - (Optional, Bool) Specifies whether Data Masking is enabled.
  + `ignore` - (Optional, Bool) Specifies whether False Alarm Masking is enabled.
  + `anti_tamper` - (Optional, Bool) Specifies whether Web Tamper Protection is enabled.
  + `anti_leakage` - (Optional, Bool) Whether the information leakage prevention is enabled.
  + `followed_action` - (Optional, Bool) Whether the Known Attack Source protection is enabled.

* `full_detection` - (Optional, Bool) Specifies the detection mode in Precise Protection.
  * `true`: full detection, Full detection finishes all threat detections before blocking requests that meet Precise Protection specified conditions.
  * `false`: instant detection. Instant detection immediately ends threat detection after blocking a request that meets Precise Protection specified conditions.

* `deep_inspection` - (Optional, Bool) The deep inspection in basic web protection.

* `header_inspection` - (Optional, Bool) The header inspection in basic web protection.

* `shiro_decryption_check` - (Optional, Bool) The shiro decryption check in basic web protection.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the policy.

* `domains` - Specifies the domain IDs.

* `created_at` - Time the policy is created. The value is a 13-digit timestamp, in ms.

## Import

WAF dedicated policies can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_policy_v1.policy_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
