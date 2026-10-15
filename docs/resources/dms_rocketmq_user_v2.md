---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_user_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-rocketmq-user-v2"
description: |-
  Manages a DMS RocketMQ user resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ user you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/user_management/index.html)

# opentelekomcloud_dms_rocketmq_user_v2

Manages a DMS RocketMQ user in the OpenTelekomCloud DMS Service.

-> User management requires access control to be enabled on the instance (`enable_acl = true`
of `opentelekomcloud_dms_rocketmq_instance_v2`).

## Example Usage

```hcl
variable "instance_id" {}
variable "secret_key" {
  sensitive = true
}

resource "opentelekomcloud_dms_rocketmq_user_v2" "test" {
  instance_id          = var.instance_id
  access_key           = "user_test"
  secret_key           = var.secret_key
  white_remote_address = "192.168.1.*"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `access_key` - (Required, String, ForceNew) Specifies the username. Starts with a letter, consists of 7 to 64
  characters, and contains only letters, digits, hyphens (-), and underscores (_).
  Changing this creates a new resource.

* `secret_key` - (Required, String) Specifies the secret key of the user. Enter 8 to 32 characters containing
  at least three of the following character types: uppercase letters, lowercase letters, digits and special characters.
  Cannot be the username or the username spelled backwards.

* `white_remote_address` - (Optional, String) Specifies the IP address whitelist, for example `192.168.1.*`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<access_key>`.
* `admin` - Indicates whether the user is an administrator.
* `default_topic_perm` - Indicates the default topic permissions. Values: **PUB**, **SUB**, **PUB|SUB**, **DENY**.
* `default_group_perm` - Indicates the default consumer group permissions. Values: **SUB**, **DENY**.

## Import

DMS RocketMQ user can be imported using the instance ID and username separated by a slash, e.g.

```
 $ terraform import opentelekomcloud_dms_rocketmq_user_v2.test 8d3c7938-dc47-4937-a30f-c80de381c5e3/user_test
```
