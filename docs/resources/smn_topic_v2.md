---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_topic_v2"
sidebar_current: "docs-opentelekomcloud-resource-smn-topic-v2"
description: |-
Manages an SMN Topic resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SMN topic you can get at
[documentation portal](https://docs.otc.t-systems.com/simple-message-notification/api-ref/apis/topic_operations)

# opentelekomcloud_smn_topic_v2

Manages a V2 topic resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew, String) The name of the topic to be created.

* `display_name` - (Optional, String) Topic display name, which is presented as the
  name of the email sender in an email message.

* `project_name` - (Optional, ForceNew, String) The project name for the topic.

* `tags` - (Optional, Map) Tags key/value pairs to associate with the instance.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `display_name` - See Argument Reference above.

* `topic_urn` - Resource identifier of a topic, which is unique.

* `push_policy` - Message pushing policy. 0 indicates that the message
  sending fails and the message is cached in the queue. 1 indicates that the
  failed message is discarded.

* `create_time` - Time when the topic was created.

* `update_time` - Time when the topic was updated.
