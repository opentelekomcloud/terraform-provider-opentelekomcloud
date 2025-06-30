---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_topic_subscription_v2"
sidebar_current: "docs-opentelekomcloud-datasource-smn-topic-subscription-v2"
description: |-
  Get details about an SMN  Topic Subscription resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SMN topic subscription you can get at
[documentation portal](https://docs.otc.t-systems.com/simple-message-notification/api-ref/apis/subscription_operations/querying_subscriptions_of_a_specified_topic.html#)

# opentelekomcloud_smn_topic_subscription_v2

Get details about an SMN topic subscription V2 resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "topic_urn" {}

data "opentelekomcloud_smn_topic_subscription_v2" "data-sub" {
  topoc_urn = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `topic_urn` - (Required, String) The resource ID of a topic

* `protocol` - (Optional, String) Subscription protocol.

* `endpoint` - (Optional, String) Message receiving endpoint.

## Attributes Reference

The following attributes are exported:

* `status` - Subscription status.

* `subscription_urn` - Resource identifier of a subscription, which is unique.

* `owner` - Project ID of the topic creator.

* `remark` - Remarks.
