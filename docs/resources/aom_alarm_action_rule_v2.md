---
subcategory: "Application Operations Management (AOM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_aom_alarm_action_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-aom-alarm-action-rule-v2"
description: |-
  Manages an AOM alarm action rule resource within T-Cloud Public (formerly OpenTelekomCloud).
---

# opentelekomcloud_aom_alarm_action_rule_v2

Manages an AOM alarm action rule resource within T-Cloud Public (formerly OpenTelekomCloud).

Up-to-date reference of API arguments for AOM alarm you can get at
[documentation portal](https://docs.otc.t-systems.com/application-operations-management-2/api-ref/apis/alarm/index.html)

## Example Usage

```hcl
variable "alarm_topic_urn" {}

resource "opentelekomcloud_aom_alarm_action_rule_v2" "test" {
  name                  = "alarm_action_rule"
  type                  = "1"
  description           = "alarm action rule"
  notification_template = "aom.built-in.template.en"

  smn_topic {
    topic_urn   = var.alarm_topic_urn
    name        = "alarm_topic"
    push_policy = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The name of the alarm action rule.
  A maximum of 100 characters are allowed. The name must not start or end with
  a special character and only letters, digits, and underscores (_) are allowed.
  Changing this parameter will create a new resource.

* `type` - (Required, String) The rule type.
  Valid values are `1` (Prometheus monitoring) and `3` (log monitoring).

* `notification_template` - (Required, String) The message template.
  The template name must reference an existing AOM notification template.

  -> **Note:** The `notification_template` must match the rule type
  (e.g. the built-in template `aom.built-in.template.en` is only valid for
  `type = "1"`), otherwise the API returns an
  `AOM.08018056: notification rule and template is not matched` error.

* `smn_topic` - (Required, List) The list of SMN topics to notify.
  A maximum of 5 topics are allowed.
  The [smn_topic](#smn_topic) structure is documented below.

* `description` - (Optional, String) The description of the rule.
  A maximum of 1,024 characters are allowed. The description must not start or
  end with an underscore (_), and only digits, letters, underscores (_),
  asterisks (*), and spaces are allowed.

* `user_name` - (Optional, String, ForceNew) The member account name.
  If omitted, the username the provider is authenticated as is used.
  Changing this parameter will create a new resource.

* `time_zone` - (Optional, String) The time zone of the rule.

<a name="smn_topic"></a>The `smn_topic` block supports:

* `topic_urn` - (Required, String) The unique resource identifier of the topic.

* `name` - (Required, String) The name of the topic.

* `push_policy` - (Required, Int) The SMN message push policy.
  Valid values are `0` and `1`.

* `display_name` - (Optional, String) The topic display name, which will be the
  name of an email sender. A maximum of 192 bytes are allowed.
  If omitted, the value is read back from the API.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the name of the alarm action rule.

* `create_time` - The creation time (UTC timestamp, in milliseconds).

* `update_time` - The modification time (UTC timestamp, in milliseconds).

* `smn_topic.status` - The status of the topic subscriber:
  `0` (the topic has been deleted or the subscription list of this topic is
  empty), `1` (the subscription object is in the subscribed state), or
  `2` (the subscription object is in the unsubscribed or canceled state).

## Import

AOM alarm action rules can be imported using the rule name, e.g.

```
terraform import opentelekomcloud_aom_alarm_action_rule_v2.test alarm_action_rule
```
