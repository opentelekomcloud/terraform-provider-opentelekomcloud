---
subcategory: "Web Application Firewall (WAF)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_alarm_notification_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-alarm-notification-v1"
description: |-
  Manages a WAF Alarm Notification resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for WAF alarm notification you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/alarm_notification)

# opentelekomcloud_waf_alarm_notification_v1

Manages a WAF alarm notification resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name = "topic_alarm"
}

resource "opentelekomcloud_waf_alarm_notification_v1" "notification_1" {
  enabled        = true
  topic_urn      = opentelekomcloud_smn_topic_v2.topic_1.id
  send_frequency = 30
  times          = 200
  threats        = ["cc", "cmdi"]
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required) Specifies whether to send an alarm notification. The options are `true` and `false`.

* `topic_urn` - (Required) Specifies the SMN topic to which an alarm is sent.

-> The selected topic must be a topic whose subscription information has been configured.

* `send_frequency` - (Required) Specifies the minimum interval between two alarms in minutes.
  The options are `5`, `15`, `30`, and `60`.

* `times` - (Required) Specifies the alarm threshold. Alarm notifications are sent when the
  number of attacks is greater than or equal to the threshold within the configured period.
  This value is greater than or equal to `1`.

* `threat` - (Required) Specifies the list of event types. Possible values are:
  * `all` refers to all types of events.
  * `cc` refers to CC attack.
  * `cmdi` refers to command injection.
  * `custom` refers to Precise Protection events.
  * `illegal` refers to invalid requests.
  * `sqli` refers to SQL injection.
  * `lfi` refers to local file inclusion.
  * `robot` refers to malicious crawlers.
  * `antitamper` refers to Web Tamper Protection events.
  * `rfi` refers to remote file inclusion.
  * `vuln` refers to other types of attacks.
  * `xss` refers to XSS attack.
  * `whiteblackip` refers to Blacklist and Whitelist events.
  * `webshell` refers to webshells.
