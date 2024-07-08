---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_subscription_v2"
sidebar_current: "docs-opentelekomcloud-resource-smn-subscription-v2"
description: |-
Manages an SMN Subscription resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SMN subscription you can get at
[documentation portal](https://docs.otc.t-systems.com/simple-message-notification/api-ref/apis/subscription_operations)

# opentelekomcloud_smn_subscription_v2

Manages a V2 subscription resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_smn_subscription_v2" "subscription_1" {
  topic_urn = opentelekomcloud_smn_topic_v2.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}

resource "opentelekomcloud_smn_subscription_v2" "subscription_2" {
  topic_urn = opentelekomcloud_smn_topic_v2.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "O&M"
}
```

## Argument Reference

The following arguments are supported:

* `topic_urn` - (Required) Specifies the resource identifier of a topic, which is unique.

* `endpoint` - (Required) Specifies the message endpoint.
  * For an HTTP subscription, the endpoint starts with http\://.
  * For an HTTPS subscription, the endpoint starts with https\://.
  * For an email subscription, the endpoint is a mail address.
  * For an SMS message subscription, the endpoint is a phone number.

* `protocol` - (Required) Specifies protocol of the message endpoint. Currently, `email`,
  `sms`, `http`, and `https` are supported.

* `remark` - (Optional) Specifies the remark information. The remarks must be a UTF-8-coded
  character string containing 128 bytes.

* `project_name` - (Optional) The project name for the subscription.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `subscription_urn` - The resource identifier of a subscription.

* `owner` - The project ID of the topic creator.

* `status` - The subscription status.
  * `0` indicates that the subscription is not confirmed.
  * `1` indicates that the subscription is confirmed.
  * `3` indicates that the subscription is canceled.
