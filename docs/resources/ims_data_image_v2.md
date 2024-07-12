---
subcategory: "Image Management Service (IMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ims_data_image_v2"
sidebar_current: "docs-opentelekomcloud-resource-ims-data-image-v2"
description: |-
  Manages a IMS Data Image resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IMS data image you can get at
[documentation portal](https://docs.otc.t-systems.com/image-management-service/api-ref/ims_apis/image)

# opentelekomcloud_ims_data_image_v2

Manages a V2 Data Image resource within OpenTelekomCloud IMS.

## Example Usage

###  Creating a data disk image using an ECS

```hcl
resource "opentelekomcloud_ims_data_image_v2" "ims_test" {
  name        = "imt_test"
  volume_id   = "54a6c3a4-8511-4d01-818f-3fe5177cbb06"
  description = "Create an image using an ECS."

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

###  Creating a data disk image in the OBS bucket

```hcl
resource "opentelekomcloud_ims_data_image_v2" "ims_test_file" {
  name        = "ims_test_file"
  image_url   = "ims-image:centos70.qcow2"
  min_disk    = 40
  os_type     = "Linux"
  description = "Create an image using a file in the OBS bucket."

  tags = {
    foo = "bar1"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the image.

* `description` - (Optional) A description of the image. Changing this creates a new image.

* `tags` - (Optional) The tags of the image.

* `volume_id` - (Optional) The ID of the ECS atatched volume that needs to be converted into an image.
  This parameter is mandatory when you create a privete image from an ECS.
  Changing this creates a new image.

* `image_url` - (Optional) The URL of the external image file in the OBS bucket.
  This parameter is mandatory when you create a private image from an external file
  uploaded to an OBS bucket. The format is *OBS bucket name:Image file name*.
  Changing this creates a new image.

* `min_disk` - (Optional) The minimum size of the system disk in the unit of GB.
  This parameter is mandatory when you create a private image from an external file
  uploaded to an OBS bucket. The value ranges from 1 GB to 1024 GB.
  Changing this creates a new image.

* `os_type` - (Optional) The OS type. It can only be Windows or Linux.
  This parameter is valid when you create a private image from an external file
  uploaded to an OBS bucket. Changing this creates a new image.

* `cmk_id` - (Optional) The master key used for encrypting an image.
  Changing this creates a new image.


## Attributes Reference

The following attributes are exported:

* `id` - A unique ID assigned by IMS.

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `tags` - See Argument Reference above.

* `volume_id` - See Argument Reference above.

* `image_url` - See Argument Reference above.

* `min_disk` - See Argument Reference above.

* `os_type` - See Argument Reference above.

* `cmk_id` - See Argument Reference above.

* `visibility` - Whether the image is visible to other tenants.

* `data_origin` - The image resource. The pattern can be 'instance,*instance_id*' or 'file,*image_url*'.

* `disk_format` - The image file format. The value can be `vhd`, `zvhd`, `raw`, `zvhd2`, or `qcow2`.

* `image_size` - The size(bytes) of the image file format.

## Import

Images can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_ims_data_image_v2.my_image 7886e623-f1b3-473e-b882-67ba1c35887f
```
