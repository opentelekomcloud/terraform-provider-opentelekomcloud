---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_topic_attribute_v2"
sidebar_current: "docs-opentelekomcloud-resource-smn-topic-attribute-v2"
description: |-
Manages an SMN Topic Attribute resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SMN topic attribute you can get at
[documentation portal](https://docs.otc.t-systems.com/simple-message-notification/api-ref/apis/topic_operations)

# opentelekomcloud_smn_topic_attribute_v2

Manages a V2 Topic Attribute resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "opentelekomcloud_smn_topic_attribute_v2" "attribute_1" {
  topic_urn       = opentelekomcloud_smn_topic_v2.topic_1.topic_urn
  attribute_name  = "access_policy"
  topic_attribute = <<EOF
{
  "Version": "2016-09-07",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__service_pub_0",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "OBS"
        ]
      },
      "Action": [
        "SMN:Publish",
        "SMN:QueryTopicDetail"
      ],
      "Resource": "${opentelekomcloud_smn_topic_v2.topic_1.topic_urn}"
    }
  ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `topic_urn` - (Required) Resource identifier of a topic, which is unique.

* `attribute_name` - (Required) Attribute name. Valid value is `access_policy`.

* `topic_attribute` - (Required) Topic attribute value. The value cannot exceed 30 KB.

## Attributes Reference

The following attributes are exported:

* `topic_urn` - See Argument Reference above.

* `attribute_name` - See Argument Reference above.

* `topic_attribute` - See Argument Reference above.

## Import

SMNv2 Topic Attribute can be imported using the `<topic_urn>/<attribute_name>`, e.g.

```shell
terraform import opentelekomcloud_smn_topic_attribute_v2.attribute_1 urn:smn:eu-de:5045c215010c440d91b2f7dce1f3753b:example/access_policy
```
