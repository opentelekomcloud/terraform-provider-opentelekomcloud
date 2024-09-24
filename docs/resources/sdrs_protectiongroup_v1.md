---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sdrs_protectiongroup_v1"
sidebar_current: "docs-opentelekomcloud-resource-sdrs-protectiongroup-v1"
description: |-
  Manages an SDRS Protection Group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SDRS protection group you can get at
[documentation portal](https://docs.otc.t-systems.com/storage-disaster-recovery-service/api-ref/sdrs_apis/protection_group)

# opentelekomcloud_sdrs_protectiongroup_v1

Manages a SDRS protection group resource within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_sdrs_domain_v1" "dom_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name        = "group_1"
  description = "test description"

  source_availability_zone = "eu-de-01"
  target_availability_zone = "eu-de-02"

  domain_id     = data.opentelekomcloud_sdrs_domain_v1.dom_1.id
  source_vpc_id = var.vpc_id
  dr_type       = "migration"
  enable        = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of a protection group.

* `description` - (Optional, String, ForceNew) The description of a protection group. Changing this creates a new group.

* `source_availability_zone` - (Required, String, ForceNew) Specifies the source AZ of a protection group. Changing this creates a new group.

* `target_availability_zone` - (Required, String, ForceNew) Specifies the target AZ of a protection group. Changing this creates a new group.

* `domain_id` - (Required, String, ForceNew) Specifies the ID of an ``active-active domain``. Changing this creates a new group.
  An ``active-active domain`` id can be extracted from ``data/opentelekomcloud_sdrs_domain_v1`` and shouldn't be confused
  with tenant ``domain``.

* `source_vpc_id` - (Required, String, ForceNew) Specifies the ID of the source VPC. Changing this creates a new group.

* `dr_type` - (Optional, String, ForceNew) Specifies the deployment model. The default value is migration indicating migration within a VPC.
  Changing this creates a new group.
* `enable` - (Optional, Boolean) Enables or disables the Protection group.


## Attributes Reference

The following attributes are exported:

* `id` - (String) ID of the protection group.
* `created_at` - (String) Time of creation of the protection group.
* `updated_at` - (String) Time of last update of the protection group.

## Import

Protection groups can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_sdrs_protectiongroup_v1.group_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
