---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_topic_v2"
sidebar_current: "docs-opentelekomcloud-resource-smn-topic-v2"
description: |-
  Get details about an SMN Topic resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SMN topic you can get at
[documentation portal](https://docs.otc.t-systems.com/simple-message-notification/api-ref/apis/topic_operations)

# opentelekomcloud_smn_topic_v2

Get details about an SMN topic V2 resource within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the topic to be fetched.


## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `display_name` - Topic display name, which is presented as the
  name of the email sender in an email message.

* `topic_urn` - Resource identifier of a topic, which is unique.

* `push_policy` - Message pushing policy. 0 indicates that the message
  sending fails and the message is cached in the queue. 1 indicates that the
  failed message is discarded.

* `create_time` - Time when the topic was created.

* `update_time` - Time when the topic was updated.
