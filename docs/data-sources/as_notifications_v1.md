---
subcategory: "Autoscaling"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_as_notifications_v1"
sidebar_current: "docs-opentelekomcloud-datasource-as-notifications-v1"
description: |-
  Retrieves notification configurations for an Auto Scaling group.
---

# opentelekomcloud_as_notifications_v1

Retrieves all SMN notification configurations for an Auto Scaling group.
This data source does not create or modify notification configurations.

See the [Auto Scaling notification query API reference](https://docs.otc.t-systems.com/auto-scaling/api-ref/apis/notifications/querying_notification_configurations_of_an_as_group.html).

## Example Usage

```hcl
data "opentelekomcloud_as_notifications_v1" "scaling" {
  scaling_group_id = opentelekomcloud_as_group_v1.group.id
}

output "notification_topics" {
  value = data.opentelekomcloud_as_notifications_v1.scaling.notifications[*].topic_urn
}
```

The example assumes an existing `opentelekomcloud_as_group_v1.group`. When
querying configurations created in the same apply, add `depends_on` references
to the relevant `opentelekomcloud_as_notification_v1` resources so the query runs
after the configurations have been created.

## Argument Reference

* `scaling_group_id` - (Required, String) The Auto Scaling group ID.

## Attributes Reference

* `id` - The Auto Scaling group ID, stable across refreshes.
* `region` - The region in which the configurations are queried.
* `notifications` - A list of notification configurations, sorted by `topic_urn`.
  Each entry exports the following attributes:
  * `topic_urn` - The URN of the SMN topic.
  * `topic_name` - The name of the SMN topic.
  * `topic_scene` - The set of events that trigger notifications:
    `SCALING_UP`, `SCALING_UP_FAIL`, `SCALING_DOWN`, `SCALING_DOWN_FAIL`,
    and `SCALING_GROUP_ABNORMAL`.

An existing group with no notification configurations returns an empty list.
An API error, including a missing or inaccessible group, is reported as an error
rather than an empty result.
