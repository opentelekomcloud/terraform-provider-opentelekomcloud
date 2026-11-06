---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_rocketmq_users_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dms-rocketmq-users-v2"
description: |-
  Get the list of DMS RocketMQ users within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS RocketMQ users you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service-rocketmq/api-ref/apis_v2_recommended/user_management/querying_the_user_list.html)

# opentelekomcloud_dms_rocketmq_users_v2

Use this data source to get the list of users of a DMS RocketMQ instance within OpenTelekomCloud.

-> The instance must have ACL enabled (`enable_acl = true`).

## Example Usage

```hcl
variable "instance_id" {}
variable "access_key" {}

data "opentelekomcloud_dms_rocketmq_users_v2" "test" {
  instance_id = var.instance_id
  access_key  = var.access_key
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `access_key` - (Optional, String) Specifies the username.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the users are queried.

* `users` - The list of users.
  The [users](#dms_rocketmq_users) structure is documented below.

<a name="dms_rocketmq_users"></a>
The `users` block supports:

* `access_key` - The username.

* `white_remote_address` - The IP address whitelist.

* `admin` - Whether the user is an administrator.

* `default_topic_perm` - The default topic permissions.

* `default_group_perm` - The default consumer group permissions.

* `topic_perms` - The special topic permissions.
  The [topic_perms](#dms_rocketmq_users_perms) structure is documented below.

* `group_perms` - The special consumer group permissions.
  The [group_perms](#dms_rocketmq_users_perms) structure is documented below.

<a name="dms_rocketmq_users_perms"></a>
The `topic_perms` and `group_perms` blocks support:

* `name` - The topic or consumer group name.

* `perm` - The permissions.
