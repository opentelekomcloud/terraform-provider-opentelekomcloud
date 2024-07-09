---
subcategory: "Elastic Volume Service (EVS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_blockstorage_volume_v2"
sidebar_current: "docs-opentelekomcloud-resource-blockstorage-volume-v2"
description: |-
Manages a BlockStorage volume resource within OpenTelekomCloud.
---

# opentelekomcloud_blockstorage_volume_v2

Manages a V2 volume resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name        = "volume_1"
  description = "first test volume"
  size        = 3

  tags = {
    foo = "bar"
    key = "value"
  }
  metadata = {
    __system__encrypted = "1"
    __system__cmkid     = "kms_id"
  }
}
```

## Argument Reference

The following arguments are supported:

* `size` - (Required) The size of the volume to create (in gigabytes). Decreasing
  this parameter creates a new volume.

* `availability_zone` - (Optional) The availability zone for the volume.
  Changing this creates a new volume.

* `consistency_group_id` - (Optional) The consistency group to place the volume in.

* `description` - (Optional) A description of the volume. Changing this updates
  the volume's description.

* `image_id` - (Optional) The image ID from which to create the volume.
  Changing this creates a new volume.

* `metadata` - (Optional) Metadata key/value pairs to associate with the volume.
  Changing this updates the existing volume metadata.
  The EVS encryption capability with KMS key can be set with the following parameters:
  * `__system__encrypted` - The default value is set to '0', which means
    the volume is not encrypted, the value '1' indicates volume is encrypted.
  * `__system__cmkid` - (Optional) The ID of the kms key.

* `tags` - (Optional) Tags key/value pairs to associate with the volume.
  Changing this updates the existing volume tags.

* `name` - (Optional) A unique name for the volume. Changing this updates the
  volume's name.

* `snapshot_id` - (Optional) The snapshot ID from which to create the volume.
  Changing this creates a new volume.

* `source_replica` - (Optional) The volume ID to replicate with.

* `source_vol_id` - (Optional) The volume ID from which to create the volume.
  Changing this creates a new volume.

* `volume_type` - (Optional) Currently, the value can be `SSD` (ultra-high I/O disk type), `SAS` (high I/O disk type), `SATA` (common I/O disk type), `co-p1` (Exclusive HPC/ SAP HANA: high I/O, performance optimized), or `uh-l1` (Exclusive HPC/ SAP HANA: ultra-high-I/O, latency optimized). Read **Note** for `uh-l1` and `co-p1`: [OTC-API](https://docs.otc.t-systems.com/en-us/api/ecs/en-us_topic_0065817708.html). Changing this creates a new volume.

* `device_type` - (Optional) The device type of volume to create. Valid options are VBD and SCSI.
  Defaults to VBD. Changing this creates a new volume.

* `cascade` - (Optional, Default:false) Specifies to delete all snapshots associated with the EVS disk.

## Attributes Reference

The following attributes are exported:

* `size` - See Argument Reference above.

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `availability_zone` - See Argument Reference above.

* `image_id` - See Argument Reference above.

* `source_vol_id` - See Argument Reference above.

* `snapshot_id` - See Argument Reference above.

* `metadata` - See Argument Reference above.

* `volume_type` - See Argument Reference above.

* `device_type` - See Argument Reference above.

* `attachment` - If a volume is attached to an instance, this attribute will
  display the Attachment ID, Instance ID, and the Device as the Instance sees it.

* `wwn` - Specifies the unique identifier used for mounting the EVS disk.

## Import

Volumes can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_blockstorage_volume_v2.volume_1 ea257959-eeb1-4c10-8d33-26f0409a755d
```
