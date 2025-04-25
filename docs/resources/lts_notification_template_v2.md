---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_notification_template_v2"
sidebar_current: "docs-opentelekomcloud-resource-lts-notification-template-v2"
description: |-
  Manages a LTS notification template resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log group you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/message_template_management/index.html)

# opentelekomcloud_lts_notification_template_v2

Manages an LTS notification template resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lts_notification_template_v2" "test" {
  name        = "my_template"
  description = "test"
  source      = "LTS"
  locale      = "en-us"

  templates {
    sub_type = "sms"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The name of the notification template.
  Changing this parameter will create a new resource.

* `source` - (Required, String) The source of the notification template.
  Currently, this parameter is mandatory to `LTS`.

* `locale` - (Required, String) Language.
  Currently, the value can be `en-us`.

* `templates` - (Required, List) The list of notification template body.
  The [templates](#NotificationTemplate) structure is documented below.

* `description` - (Optional, String) The description of the notification template.

<a name="NotificationTemplate"></a>
The `templates` block supports:

* `sub_type` - (Required, String) The type of the sub-template.
  Only the following five types are supported: `sms`, `email`.

* `content` - (Required, String) The content of the sub-template.
  In the sub-template body, only the following variables are supported for the variables following the `$` symbol.
  The supported variables vary according to the alarm type (keyword alarm and SQL alarm).

    + Common variables:
      * Alarm severity: **${event_severity}**.
      * Occurrence time: **${starts_at}**.
      * Alarm source: **$event.metadata.resource_provider**.
      * Resource type: **$event.metadata.resource_type**.
      * Resource ID: **${resources}**.
      * Expression: **$event.annotations.condition_expression**.
      * current value: **$event.annotations.current_value**.
      * Statistical period: **$event.annotations.frequency**.

    + Keywords alarm specific variable:
      * query time: **$event.annotations.results[0].time**.
      * Run the **$event.annotations.results[0].raw_results** command to query LTSs.

    + SQL alarm specific variable:
      * LTS group/stream name: **$event.annotations.results[0].resource_id**.
      * Query statement: **$event.annotations.results[0].sql**.
      * Query time: **$event.annotations.results[0].time**.
      * Query URL: **$event.annotations.results[0].url**.
      * Run the **$event.annotations.results[0].raw_results** command to query LTSs.

  -> semicolon(;) after variable must be added. Otherwise, the template will fail to be replaced.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals `name`.

* `region` - Shows the region in the template resource created.

## Import

The LTS notification template can be imported using the `id` which equals `name`, e.g.

```bash
$ terraform import opentelekomcloud_lts_notification_template_v2.test <id>
```
