---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated Geolocation Access Control rule you can get at
`https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_geolocation_access_control_rule.html`.

# opentelekomcloud_waf_dedicated_geo_ip_rule_v1

Manages a WAF Dedicated Geolocation Access Control Rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_gi"
}

resource "opentelekomcloud_waf_dedicated_geo_ip_rule_v1" "rule_1" {
  policy_id   = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  region_code = "BR"
  action      = 0
  name        = "test"
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew, String) The WAF policy ID. Changing this creates a new rule.

* `region_code` - (Required, String) Applicable regions. The value can be the region code. For more geographical location codes, see docs "Appendix - Geographic Location Codes."
  Values:
  + CA: Canada
  + US: USA
  + AU: Australia
  + IN: India
  + JP: Japan
  + UK: United Kingdom
  + FR: France
  + DE: Germany
  + BR: Brazil
  + Ukraine: Ukraine
  + Pakistan: Pakistan
  + Palestine: Palestine
  + Israel: Israel
  + Iraq: Afghanistan
  + Libya: Libya
  + Turkey: Turkey
  + Thailand: Thailand
  + Singapore: Singapore
  + South Africa: South Africa
  + Mexico: Mexico
  + Peru: Peru

* `action` - (Required, Int) Protective action.
  The value can be:
  + 0: WAF blocks the requests that hit the rule.
  + 1: WAF allows the requests that hit the rule.
  + 2: WAF only logs the requests that hit the rule.

* `name` - (Optional, String) Rule name.

* `description` - (Optional, String) Rule description

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the rule.

* `status` - Rule status. The value can be:
  + `0`: The rule is disabled.
  + `1`: The rule is enabled.

* `created_at` - Timestamp the rule is created.

## Import

Dedicated WAF Web Geolocation Access Control rules can be imported using `policy_id/id`, e.g.

```sh
terraform import opentelekomcloud_waf_dedicated_geo_ip_rule_v1.rule_1 ff95e71c8ae74eba9887193ab22c5757/b39f3a5a1b4f447a8030f0b0703f47f5
```
