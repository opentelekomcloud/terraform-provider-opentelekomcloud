---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_keywords_alarm_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-lts-keywords-alarm-rule-v2"
description: |-
  Manages a LTS keywords alarm rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log group you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/keyword_alarm_rules/index.html)


# opentelekomcloud_lts_keywords_alarm_rule_v2

Manages an LTS keywords alarm rule resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}
variable "topic_name" {}
variable "topic_urn" {}

resource "opentelekomcloud_lts_keywords_alarm_rule_v2" "test" {
  name        = "name"
  description = "created by terraform"
  severity    = "CRITICAL"

  notification_frequency = 5

  keywords_requests {
    keyword                = "key_words"
    condition              = ">"
    number                 = 100
    log_group_id           = var.log_group_id
    log_stream_id          = var.log_stream_id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type = "HOURLY"
  }

  notification_rule {
    language  = "en-us"
    timezone  = "xx/xx"
    user_name = "test"

    topics {
      name      = var.topic_name
      topic_urn = var.topic_urn
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the keywords alarm rule.
  The value can contain no more than 64 characters.
  Changing this parameter will create a new resource.

* `keywords_requests` - (Required, List) Specifies the keywords requests.
  The [KeywordsRequests](#KeywordsAlarmRule_KeywordsRequests) structure is documented below.

* `frequency` - (Required, List) Specifies the alarm frequency configurations.
  The [Frequency](#KeywordsAlarmRule_Frequency) structure is documented below.

* `severity` - (Required, String) Specifies the alarm level.
  The value can be: **INFO**, **MINOR**, **MAJOR** and **CRITICAL**.

* `description` - (Optional, String) Specifies the description of the keywords alarm rule.

* `send_notifications` - (Optional, Bool) Specifies whether to send notifications.

* `notification_frequency` - (Required, Int) Specifies notification frequency, in minutes.

* `notification_rule` - (Optional, List) Specifies the notification rule.
  The [NotificationRule](#KeywordsAlarmRule_NotificationRule) structure is documented below.

* `trigger_condition_count` - (Optional, Int) Specifies the count to trigger the alarm.
  Defaults to `1`.

* `trigger_condition_frequency` - (Optional, Int) Specifies the frequency to trigger the alarm.
  Defaults to `1`.

* `send_recovery_notifications` - (Optional, Bool) Specifies whether to send recovery notifications.

* `recovery_frequency` - (Optional, Int) Specifies the frequency to recover the alarm.
  Defaults to `3`.

<a name="KeywordsAlarmRule_KeywordsRequests"></a>
The `KeywordsRequests` block supports:

* `keyword` - (Required, String) Specifies the keywords.

* `condition` - (Required, String) Specifies the keywords request condition.
  The value can be: **>=**, **<=**, **<** and **>**.

* `number` - (Required, Int) Specifies the line number.

* `log_stream_id` - (Required, String) Specifies the log stream id.

* `log_group_id` - (Required, String) Specifies the log group id.

* `search_time_range_unit` - (Required, String) Specifies the unit of search time range.
  The value can be: **minute** and **hour**.

* `search_time_range` - (Required, Int) Specifies the search time range.
  + When the `search_time_range_unit` is **minute**, the value ranges from `1` to `60`.
  + When the `search_time_range_unit` is **hour**, the value ranges from `1` to `24`.

<a name="KeywordsAlarmRule_Frequency"></a>
The `Frequency` block supports:

* `type` - (Required, String) Specifies the frequency type.
  The value can be: **CRON**, **HOURLY**, **DAILY**, **WEEKLY** and **FIXED_RATE**.

* `cron_expression` - (Optional, String) Specifies the cron expression.
  This parameter is used when `type` is set to **CRON**.

* `hour_of_day` - (Optional, Int) Specifies the hour of day.
  This parameter is used when `type` is set to **DAILY** or **WEEKLY**.
  The value ranges from `0` to `23`.

* `day_of_week` - (Optional, Int) Specifies the day of week.
  This parameter is used when `type` is set to **WEEKLY**.
  The value ranges from `1` to `7`. `1` means Sunday.

* `fixed_rate_unit` - (Optional, String) Specifies the unit of fixed rate.
  The value can be: **minute** and **hour**.

* `fixed_rate` - (Optional, Int) Specifies the unit fixed rate.
  This parameter is used when `type` is set to **FIXED_RATE**.
  + When the `fixed_rate_unit` is **minute**, the value ranges from `1` to `60`.
  + When the `fixed_rate_unit` is **hour**, the value ranges from `1` to `24`

<a name="KeywordsAlarmRule_NotificationRule"></a>
The `NotificationRule` block supports:

* `template_name` - (Optional, String) Specifies the notification template name.

* `user_name` - (Required, String) Specifies the username.

* `topics` - (Required, List) Specifies the SMN topics.
  The [Topic](#KeywordsAlarmRule_Topic) structure is documented below.

* `timezone` - (Optional, String) Specifies the timezone.

* `language` - (Optional, String) Specifies the notification language.
  The value can be **en-us**.

<a name="KeywordsAlarmRule_Topic"></a>
The `NotificationRuleTopic` block supports:

* `name` - (Required, String) Specifies the topic name.
  Changing this parameter will create a new resource.

* `topic_urn` - (Required, String) Specifies the topic URN.
  Changing this parameter will create a new resource.

* `display_name` - (Optional, String) Specifies the display name.
  This will be shown as the sender of the message.
  Changing this parameter will create a new resource.

* `push_policy` - (Optional, Int) Specifies the push policy.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domain_id` - The domain ID.

* `created_at` - The creation time of the alarm rule.

* `updated_at` - The last update time of the alarm rule.

* `region` - Shows the region in the rule resource created.

* `status` - Status of the rule.

## Import

The keywords alarm rule can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_lts_keywords_alarm_rule_v2.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `notification_rule`.

```hcl
resource "opentelekomcloud_lts_keywords_alarm_rule_v2" "test" {

  lifecycle {
    ignore_changes = [
      notification_rule,
    ]
  }
}
```
