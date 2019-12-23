---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_logtank_topic_v2"
sidebar_current: "docs-opentelekomcloud-resource-logtank-topic-v2"
description: |-
  log topic management
---

# opentelekomcloud\_logtank\_topic\_v2

Manage a log topic resource within OpenTelekomCloud.

## Example Usage

### create a log topic

```hcl
resource "opentelekomcloud_logtank_group_v2" "test_group" {
  topic_name  = "test_group"
}

resource "opentelekomcloud_logtank_topic_v2" "test_topic" {
  group_id   = "${opentelekomcloud_logtank_group_v2.test_group.id}"
  topic_name = "test1"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required)
  Specifies the ID of a created log group.
  Changing this parameter will create a new resource.

* `topic_name` - (Required)
  Specifies the log topic name.
  Changing this parameter will create a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The log topic ID.

* `group_id` - See Argument Reference above.

* `topic_name` - See Argument Reference above.

* `index_enabled` - Specifies the search switch. When index is enabled, the topic allows you to search for logs by keyword.

## Import

Log topic can be imported using the logtank group ID and topic ID separated by a slash, e.g.

```
$ terraform import opentelekomcloud_logtank_topic_v2.topic_1 393f2bfd-2244-11ea-adb7-286ed488c87f/72855918-20b1-11ea-80e0-286ed488c880
```
