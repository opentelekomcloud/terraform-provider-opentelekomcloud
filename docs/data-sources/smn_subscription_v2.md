---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_subscription_v2"
sidebar_current: "docs-opentelekomcloud-datasource-smn-subscription-v2"
description: |-
  Get details about an SMN Subscription resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SMN subscription you can get at
[documentation portal](https://docs.otc.t-systems.com/simple-message-notification/api-ref/apis/subscription_operations/querying_subscriptions.html#smn-api-52001)

# opentelekomcloud_smn_subscription_v2

Get details about an SMN subscription V2 resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "subscription_urn" {}

data "opentelekomcloud_smn_subscription_v2" "data-sub" {
  subscription_urn = var.subscription_urn
}
```

## Argument Reference

The following arguments are supported:

* `subscription_urn` - (Required, String) The urn of subscription to be fetched.

## Attributes Reference

The following attributes are exported:

* `protocol` - Subscription protocol.

* `status` - Subscription status.

* `endpoint` - Message receiving endpoint.

* `topic_urn` - Resource identifier of a topic, which is unique.

* `owner` - Project ID of the topic creator.

* `remark` - Remarks.
