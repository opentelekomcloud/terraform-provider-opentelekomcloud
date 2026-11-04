---
subcategory: "Application Operations Management (AOM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_aom_alarm_silence_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-aom-alarm-silence-rule-v2"
description: |-
  Manages an AOM alarm silence rule resource within T-Cloud Public (formerly OpenTelekomCloud).
---

# opentelekomcloud_aom_alarm_silence_rule_v2

Manages an AOM alarm silence rule resource within T-Cloud Public (formerly OpenTelekomCloud).

Up-to-date reference of API arguments for AOM alarm you can get at
[documentation portal](https://docs.otc.t-systems.com/application-operations-management-2/api-ref/apis/alarm/index.html)

## Example Usage

```hcl
resource "opentelekomcloud_aom_alarm_silence_rule_v2" "test" {
  name        = "alarm_silence_rule"
  description = "silence Info alarms every day"
  timezone    = "UTC"

  match {
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      value   = ["Info"]
    }
  }

  mute_config {
    type      = "DAILY"
    starts_at = 0
    ends_at   = 86399
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The name of the silence rule.
  A maximum of 100 characters are allowed. The name must not start or end with
  an underscore (_) and only letters, digits, and underscores are allowed.
  Changing this parameter will create a new resource.

* `match` - (Required, List) The match conditions of the rule.
  Up to 10 serial or parallel conditions can be created.
  The [match](#match) structure is documented below.

* `mute_config` - (Required, List) The time when the rule takes effect.
  The [mute_config](#mute_config) structure is documented below.

* `description` - (Optional, String) The description of the rule.

* `timezone` - (Optional, String) The time zone of the rule.

* `user_id` - (Optional, String, ForceNew) The user ID.
  Changing this parameter will create a new resource.

<a name="match"></a>The `match` block supports:

* `conditions` - (Required, List) The list of match conditions of a group.
  The `conditions` structure is documented below.

The `conditions` block supports:

* `key` - (Required, String) The key specified in the metadata for matching.

* `operate` - (Required, String) The match mode.
  Valid values are `EXIST` (exist), `REGEX` (regular expression), and
  `EQUALS` (equal to).

* `value` - (Optional, List) The list of values corresponding to the key.
  If `operate` is set to `EXIST`, this parameter is left blank.

<a name="mute_config"></a>The `mute_config` block supports:

* `type` - (Required, String) The type of the time for the silence rule to
  take effect. Valid values are `FIXED` (fixed time), `DAILY` (certain time
  every day), `WEEKLY` (certain time every week), and `MONTHLY` (certain time
  every month).

* `starts_at` - (Optional, Int) The start time of the silence rule
  (in milliseconds).

* `ends_at` - (Optional, Int) The end time of the silence rule
  (in milliseconds). If the silence rule is permanent, set this parameter
  to `0`.

* `scope` - (Optional, List) The scope of the rule.
  When `type` is set to `WEEKLY` or `MONTHLY`, this parameter is mandatory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the name of the silence rule.

* `create_time` - The creation time (UTC timestamp, in milliseconds).

* `update_time` - The modification time (UTC timestamp, in milliseconds).

## Import

AOM alarm silence rules can be imported using the rule name, e.g.

```
terraform import opentelekomcloud_aom_alarm_silence_rule_v2.test alarm_silence_rule
```
