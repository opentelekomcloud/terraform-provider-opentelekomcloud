---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sdrs_replication_pair_v1"
sidebar_current: "docs-opentelekomcloud-resource-sdrs-replication-pair-v1"
description: |-
Manages an SDRS Replication Pair resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SDRS replication pair you can get at
[documentation portal](https://docs.otc.t-systems.com/storage-disaster-recovery-service/api-ref/sdrs_apis/replication_pair/)

# opentelekomcloud_sdrs_replication_pair_v1

Manages a SDRS replication pair resource within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_sdrs_domain_v1" "dom_1" {}

resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "eu-de-02"
  volume_type       = "SATA"
  size              = 12
}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name        = "group_1"
  description = "test description"

  source_availability_zone = "eu-de-01"
  target_availability_zone = "eu-de-02"

  domain_id     = data.opentelekomcloud_sdrs_domain_v1.dom_1.id
  source_vpc_id = var.vpc_id
  dr_type       = "migration"
}

resource "opentelekomcloud_sdrs_replication_pair_v1" "pair_1" {
  name                 = "replication_1"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  volume_id            = opentelekomcloud_evs_volume_v3.volume_1.id
  description          = "description"
  delete_target_volume = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of a replication pair. The name can contain a maximum of 64 characters.
  The value can contain only letters (a to z and A to Z), digits (0 to 9), dots (.), underscores (_), and hyphens (-).

* `group_id` - (Required, String, ForceNew) Specifies the ID of a protection group.

* `volume_id` - (Required, String, ForceNew) Specifies the ID of the production site disk.
  When the provider is successfully invoked, the disaster recovery site disk will be automatically created.

* `delete_target_volume` - (Optional, Bool) Specifies whether to delete the disaster recovery site disk.
  The default value is **false**.

* `description` - (Optional, String, ForceNew) Specifies the description of a replication pair. The value can contain
  a maximum of 64 characters and angle brackets (<) and (>) are not allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `region` -  The resource region.

* `fault_level` - The fault level of a replication pair.
    + 0: No fault occurs.
    + 2: The disk of the current production site does not have read/write permissions. In this case, you are advised to
      perform a failover.
    + 5: The replication link is disconnected. In this case, a failover is not allowed. Contact the customer service to
      obtain service support.

* `replication_model` - The replication mode of a replication pair. The default value is **hypermetro**,
  indicating synchronous replication.

* `status` - The status of a replication pair.

* `target_volume_id` - The ID of the disk in the protection availability zone.

## Import

The SDRS replication pair can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_sdrs_replication_pair_v1.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `delete_target_volume`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```
resource "opentelekomcloud_sdrs_replication_pair_v1" "test" {
  ...

  lifecycle {
    ignore_changes = [
      delete_target_volume,
    ]
  }
}
```
