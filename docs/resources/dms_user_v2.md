---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_user_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-user-v2"
description: |-
Manages a DMS User resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS user you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service/api-ref/apis_v2_recommended/user_management/index.html)

# opentelekomcloud_dms_user_v2

Manages a DMS user for the OpenTelekomCloud DMS Service Instances (Kafka Premium/Platinum).

~>
  User management is supported only when SASL is enabled for the Kafka instance.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_dms_user_v2" "user_1" {
  instance_id = instance_id
  username    = "Test-user"
  password    = "Dmstest@123@"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Indicates the ID of primary DMS instance.

* `username` - (Required) Indicates a username. A username consists of 4 to 64 characters
  and supports only letters, digits, and hyphens (-).

* `password` - (Required) Indicates the password of an instance. An instance password
  must meet the following complexity requirements: Must be 8 to 32 characters long.
  Must contain at least 2 of the following character types: lowercase letters, uppercase
  letters, digits, and special characters (`~!@#$%^&*()-_=+\|[{}]:'",<.>/?`).

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `role` - Specifies user role.

* `default_app` - Specifies whether an application is the default application.

* `creation_time` - Specifies the time when a user was created.
