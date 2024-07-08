---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sdrs_protected_instance_v1"
sidebar_current: "docs-opentelekomcloud-resource-sdrs-protected-instance-v1"
description: |-
Manages an SDRS Protected Instance resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SDRS protected instance you can get at
[documentation portal](https://docs.otc.t-systems.com/storage-disaster-recovery-service/api-ref/sdrs_apis/protected_instance)

# opentelekomcloud_sdrs_protected_instance_v1

Manages a SDRS protected instance resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name = "group_1"

  source_availability_zone = "eu-de-01"
  target_availability_zone = "eu-de-02"

  domain_id     = var.domain_id
  source_vpc_id = var.vpc_id
  dr_type       = "migration"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = var.image_id
  flavor   = var.flavor
  vpc_id   = var.vpc_id

  nics {
    network_id = var.network_id
  }

  availability_zone = "eu-de-01"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_create"
  description          = "some interesting description"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a protected instance.

* `description` - (Optional) The description of a protected instance. Changing this creates a new instance. Changing this will create a new resource.

* `group_id` - (Required) Specifies the ID of the protection group where a protected instance is added. Changing this will create a new resource.

* `server_id` - (Required) Specifies the ID of the protected ECS instance. Changing this will create a new resource.

-> When the API is successfully invoked, the disaster recovery instance will be automatically created.

* `subnet_id` - (Optional) Specifies the network ID of the subnet. Changing this will create a new resource.

* `ip_address` - (Optional) Specifies the IP address of the primary NIC on the DR site server.
  This parameter is valid only when `subnet_id` is specified. Changing this will create a new resource.

* `delete_target_server` - (Optional) Specifies whether to delete the DR site server. The default value is `false`.

* `delete_target_eip` - (Optional) Specifies whether to delete the EIP of the DR site server. The default value is `false`.

* `tags` - (Optional) Tags key/value pairs to associate with the instance.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` -  ID of the protected instance.

* `priority_station` - Specifies the current production site AZ of the protection group containing the protected instance.
  * `source`: indicates that the current production site AZ is the `source_availability_zone` value.
  * `target`: indicates that the current production site AZ is the `target_availability_zone` value.

* `target_id` - Specifies the DR site server ID.

* `created_at` - Specifies the time when a protected instance was created.

* `updated_at` - Specifies the time when a protected instance was updated.

## Import

Protected instances can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_sdrs_protected_instance_v1.instance_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
