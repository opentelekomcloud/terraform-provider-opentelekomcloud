---
subcategory: "Elastic Volume Service (EVS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_evs_volume_v3"
sidebar_current: "docs-opentelekomcloud-resource-evs-volume-v3"
description: |-
  Manages an EVS resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for EVS you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-volume-service/api-ref/unrecommended_apis/openstack_cinder_api_v3)

# opentelekomcloud_evs_volume_v3

Manages a V3 volume resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "eu-de-01"
  volume_type       = "SATA"
  size              = 20

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Use KMS encryption

```hcl
resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "eu-de-01"
  volume_type       = "SATA"
  kms_id            = var.kms_id
  size              = 20

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) The availability zone for the volume.
  Changing this creates a new volume.

* `volume_type` - (Required) The type of volume to create.
  Currently, the value can be `SSD`, `SAS`, `SATA`, `co-p1`, or `uh-l1`.
  Changing this creates a new volume.

* `name` - (Optional) A unique name for the volume. Changing this updates the volume's name.

* `size` - (Optional) The size of the volume to create (in gigabytes). This parameter is mandatory when
  you create an empty EVS disk or use an image or a snapshot to create an EVS disk.
  _Decreasing_ this value creates a new volume.

* `description` - (Optional) A description of the volume. Changing this updates the volume's description.

* `image_id` - (Optional) The image ID from which to create the volume.
  Changing this creates a new volume.

* `backup_id` - (Optional) The backup ID from which to create the volume.
  Changing this creates a new volume.

* `snapshot_id` - (Optional) The snapshot ID from which to create the volume.
  Changing this creates a new volume.

* `tags` - (Optional) Tags key/value pairs to associate with the volume.
  Changing this updates the existing volume tags.

* `multiattach` - (Optional) Specifies whether the disk is shareable. The default value is `false`.
  Changing this creates a new volume.

* `kms_id` - (Optional) The Encryption KMS ID to create the volume.
  Changing this creates a new volume.

* `device_type` - (Optional) The device type of volume to create. Valid options are VBD and SCSI.
  Defaults to `VBD`. Changing this creates a new volume.

* `cascade` - (Optional) Specifies to delete all snapshots associated with the EVS disk. Default is `false`.

## Attributes Reference

The following attributes are exported:

* `availability_zone` - See Argument Reference above.

* `volume_type` - See Argument Reference above.

* `name` - See Argument Reference above.

* `size` - See Argument Reference above.

* `description` - See Argument Reference above.

* `image_id` - See Argument Reference above.

* `backup_id` - See Argument Reference above.

* `snapshot_id` - See Argument Reference above.

* `tags` - See Argument Reference above.

* `multiattach` - See Argument Reference above.

* `kms_id` - See Argument Reference above.

* `device_type` - See Argument Reference above.

* `attachment` - If a volume is attached to an instance, this attribute will
  display the Attachment ID, Instance ID, and the Device as the Instance sees it.

* `wwn` - Specifies the unique identifier used for mounting the EVS disk.

## Import

Volumes can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_evs_volume_v3.volume_1 14a80bc7-c12c-4fe0-a38a-cb77eeac9bd6
```
