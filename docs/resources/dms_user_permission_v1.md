---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_user_permission_v1"
sidebar_current: "docs-opentelekomcloud-resource-dms-user-permission-v1"
description: |-
Manages a DMS User Permissions resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS user permissions you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service/api-ref/apis_v2_recommended/user_management/index.html)

# opentelekomcloud_dms_user_permission_v1

Manages a DMS topic permissions for users for the OpenTelekomCloud DMS Service Instances (Kafka Premium/Platinum).

~>
  Topic permission management is supported only when SASL is enabled for the Kafka instance.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_dms_user_v2" "user_1" {
  instance_id = instance_id
  username    = "Test-user"
  password    = "Dmstest@123"
}

resource "opentelekomcloud_dms_user_v2" "user_2" {
  instance_id = instance_id
  username    = "Test-user2"
  password    = "Dmstest@123"
}

resource "opentelekomcloud_dms_user_permission_v1" "perm_1" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  topic_name  = "test-topic"
  policies {
    username      = opentelekomcloud_dms_user_v2.user_1.id
    access_policy = "all"
  }

  policies {
    username      = opentelekomcloud_dms_user_v2.user_2.id
    access_policy = "sub"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Indicates the ID of primary DMS instance.

* `topic_name` - (Required) Indicates the name of a topic.

* `policies` - (Required) Indicates policy configuration for the topic.
  Supported fields:
  * `username` - (Required) DMS instance user name.
  * `access_policy` - (Required) Permission type. Possible values:
    * `all`: publish and subscribe permissions.
    * `pub`: publish permissions.
    * `sub`: subscribe permissions.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `owner` - Indicates whether the user is the one selected during topic creation.

* `topic_type` - Indicates topic type. `0`: common topic; `1`: system (internal) topic.
