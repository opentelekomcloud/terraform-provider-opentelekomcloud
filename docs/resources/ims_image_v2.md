---
subcategory: "Image Management Service (IMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ims_image_v2"
sidebar_current: "docs-opentelekomcloud-resource-ims-image-v2"
description: |-
  Manages a IMS Image resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IMS image you can get at
[documentation portal](https://docs.otc.t-systems.com/image-management-service/api-ref/ims_apis/image)

# opentelekomcloud_ims_image_v2

Manages a V2 Image resource within OpenTelekomCloud IMS.

## Example Usage

###  Creating an image using an ECS

```hcl
resource "opentelekomcloud_ims_image_v2" "ims_test" {
  name        = "imt_test"
  instance_id = "54a6c3a4-8511-4d01-818f-3fe5177cbb06"
  description = "Create an image using an ECS."

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

###  Creating an image in the OBS bucket

```hcl
resource "opentelekomcloud_ims_image_v2" "ims_test_file" {
  name        = "ims_test_file"
  image_url   = "ims-image:centos70.qcow2"
  min_disk    = 40
  description = "Create an image using a file in the OBS bucket."

  tags = {
    foo = "bar1"
    key = "value"
  }
}
```

###  Creating an image using an Volume

```hcl
resource "opentelekomcloud_ims_image_v2" "image_volume" {
  name        = "image_volume"
  volume_id   = "54a6c3a4-8511-4d01-818f-3fe5177cbb07"
  os_version  = "Debian GNU/Linux 10.0.0 64bit"
  description = "created by Terraform"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the image.

* `description` - (Optional, String, ForceNew) A description of the image. Changing this creates a new image.

* `min_ram` - (Optional, Integer, ForceNew) The minimum memory of the image in the unit of MB.
  The default value is 0, indicating that the memory is not restricted.
  Changing this creates a new image.

* `max_ram` - (Optional, Integer, ForceNew) The maximum memory of the image in the unit of MB.
  Changing this creates a new image.

* `tags` - (Optional, List) The tags of the image.

* `instance_id` - (Optional, String, ForceNew) The ID of the ECS that needs to be converted into an image.
  This parameter is mandatory when you create a private image from an ECS.
  Changing this creates a new image.

* `volume_id` - (Optional, String, ForceNew) Specifies the data disk ID.
  This parameter is mandatory when you create a private image from a volume.
  Changing this creates a new image.

* `image_url` - (Optional, String, ForceNew) The URL of the external image file in the OBS bucket.
  This parameter is mandatory when you create a private image from an external file
  uploaded to an OBS bucket. The format is *OBS bucket name:Image file name*.
  Changing this creates a new image.

* `min_disk` - (Optional, Integer, ForceNew) The minimum size of the system disk in the unit of GB.
  This parameter is mandatory when you create a private image from an external file
  uploaded to an OBS bucket. The value ranges from 1 GB to 1024 GB.
  Changing this creates a new image.

* `os_version` - (Optional, String, ForceNew) The OS version.
  This parameter is valid when you create a private image from an external file.
  This parameter is mandatory when you create a private image from a volume.
  uploaded to an OBS bucket. Changing this creates a new image.

* `is_config` - (Optional, Boolean, ForceNew) If automatic configuration is required, set the value to true.
  Otherwise, set the value to false. Changing this creates a new image.

* `cmk_id` - (Optional, String, ForceNew) The master key used for encrypting an image.
  Changing this creates a new image.

* `type` - (Optional, String, ForceNew) The image type. Must be one of `ECS`, `FusionCompute`, `BMS`,
  `Ironic` or `IsoImage`. Changing this creates a new image.

* `hw_firmware_type` - (Optional, String) Specifies the boot mode. The value can be `bios` or `uefi`.

## Attributes Reference

In additin to the arguments defined above, the following attributes are exported:

* `id` - A unique ID assigned by IMS.

* `visibility` - Whether the image is visible to other tenants.

* `data_origin` - The image resource. The pattern can be 'instance,*instance_id*' or 'file,*image_url*'.

* `disk_format` - The image file format. The value can be `vhd`, `zvhd`, `raw`, `zvhd2`, or `qcow2`.

* `image_size` - The size(bytes) of the image file format.

* `file` - The URL for uploading and downloading the image file.

## Import

Images can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_ims_image_v2.my_image 7886e623-f1b3-473e-b882-67ba1c35887f
```
