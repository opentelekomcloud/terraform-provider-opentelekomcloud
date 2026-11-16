---
subcategory: "Autoscaling"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_as_notification_v1"
sidebar_current: "docs-opentelekomcloud-resource-as-notification-v1"
description: |-
  Manages an Auto Scaling notification configuration within OpenTelekomCloud.
---

# opentelekomcloud_as_notification_v1

Manages the association between an Auto Scaling group and an SMN topic. Auto
Scaling publishes notifications to the topic when the configured events occur.

See the [Auto Scaling notification API reference](https://docs.otc.t-systems.com/auto-scaling/api-ref/apis/notifications/index.html).

Create the SMN topic and its subscriptions before configuring notifications.
Ensure that the topic's access policy permits Auto Scaling to publish messages.
An Auto Scaling group supports at most five notification topics.

Each resource manages one group/topic pair, leaving other notification topics
unchanged. Do not manage the same pair with multiple resources. The API does not
document replacing an existing set of events in place, so changing `topic_scene`
deletes and recreates the association. Notifications to that topic may be
interrupted during replacement.

## Example Usage

```hcl
resource "opentelekomcloud_smn_topic_v2" "scaling" {
  name = "scaling-events"
}

resource "opentelekomcloud_smn_subscription_v2" "scaling" {
  topic_urn = opentelekomcloud_smn_topic_v2.scaling.id
  protocol  = "email"
  endpoint  = "scaling-alerts@example.com"
}

resource "opentelekomcloud_as_notification_v1" "scaling" {
  depends_on = [opentelekomcloud_smn_subscription_v2.scaling]

  scaling_group_id = opentelekomcloud_as_group_v1.group.id
  topic_urn        = opentelekomcloud_smn_topic_v2.scaling.id
  topic_scene      = ["SCALING_UP_FAIL", "SCALING_DOWN_FAIL", "SCALING_GROUP_ABNORMAL"]
}
```

The example assumes an existing `opentelekomcloud_as_group_v1.group`. Confirm the
SMN subscription before expecting notifications to be delivered.

## Argument Reference

* `scaling_group_id` - (Required, String, ForceNew) The Auto Scaling group ID.
  Changing this creates a new notification configuration.

* `topic_urn` - (Required, String, ForceNew) The URN of the SMN topic.
  Changing this creates a new notification configuration.

* `topic_scene` - (Required, Set of Strings, ForceNew) One or more events that
  trigger notifications. Order and duplicate entries do not affect the set.
  Supported values are `SCALING_UP`, `SCALING_UP_FAIL`, `SCALING_DOWN`,
  `SCALING_DOWN_FAIL`, and `SCALING_GROUP_ABNORMAL`.
  Changing this replaces the notification configuration.

## Attributes Reference

* `id` - The composite identifier `<scaling_group_id>/<topic_urn>`.
* `region` - The region in which the notification configuration is managed.
* `topic_name` - The name of the SMN topic.

## Timeouts

* `create` - (Default `5 minutes`) Waiting for the configuration to become visible.
* `delete` - (Default `5 minutes`) Waiting for the configuration to disappear.

## Import

Import an existing notification configuration using the group ID and topic URN
separated by a slash:

```shell
terraform import opentelekomcloud_as_notification_v1.scaling 'group-id/urn:smn:eu-de:project-id:scaling-events'
```

Import reads the existing configuration; it does not create an SMN topic or
subscription.
