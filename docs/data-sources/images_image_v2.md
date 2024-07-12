---
subcategory: "Image Management Service (IMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_images_image_v2"
sidebar_current: "docs-opentelekomcloud-datasource-images-image-v2"
description: |-
  Get the ID of an available image from OpenTelekomCloud
---

Up-to-date reference of API arguments for Image you can get at
[documentation portal](https://docs.otc.t-systems.com/image-management-service/api-ref/native_openstack_apis/image_native_openstack_apis/querying_images_native_openstack_api.html#en-us-topic-0060804959)

# opentelekomcloud_images_image_v2

Use this data source to get the ID of an available OpenTelekomCloud image.

## Example Usage

### Get Ubuntu_20.04 latest

```hcl
data "opentelekomcloud_images_image_v2" "ubuntu" {
  name = "Standard_Ubuntu_20.04_latest"
}
```

### Get most recent Debian

```hcl
data "opentelekomcloud_images_image_v2" "latest-debian" {
  name_regex  = "^Standard_Debian.?"
  most_recent = true
}
```

## Argument Reference

* `most_recent` - (Optional) If more than one result is returned, use the most recent image.

* `name` - (Optional) The name of the image.

* `name_regex` - (Optional) A regex string to apply to the images list.
  This allows more advanced filtering not supported from the OpenTelekomCloud API.
  This filtering is done locally on what OpenTelekomCloud returns.

* `owner` - (Optional) The owner (UUID) of the image.

* `size_min` - (Optional) The minimum size (in bytes) of the image to return.

* `size_max` - (Optional) The maximum size (in bytes) of the image to return.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.

* `sort_key` - (Optional) Sort images based on a certain key. Defaults to `name`.

* `tag` - (Optional) Search for images with a specific tag.

* `visibility` - (Optional) The visibility of the image. Must be one of
   `public`, `private`, `community`, or `shared`. Defaults to `private`.

-> If more or less than a single match is returned by the search, Terraform will fail.
Ensure that your search is specific enough to return a single IMS ID only, or use `most_recent`
to choose the most recent one.

## Attributes Reference

`id` is set to the ID of the found image. In addition, the following attributes are exported:

* `backup_id` - Specifies the backup ID.

* `checksum` - The checksum of the data associated with the image.

* `created_at` - The date the image was created.

* `container_format` - The format of the image's container.

* `data_origin` - Specifies the image source.

* `description` - Specifies the image description.

* `disk_format` - The format of the image's disk.

* `file` - the trailing path after the glance endpoint that represent the
  location of the image, or the path to retrieve it.

* `image_source_type` - Specifies the image backend storage type. Only `UDS` is currently supported.

* `image_type` - Specifies the image type.

* `is_registered` - Specifies whether the image is available.

* `login_user` - Specifies default image login user.

* `metadata` - The metadata associated with the image.
  Image metadata allow for meaningfully define the image properties
  and tags. See http://docs.openstack.org/developer/glance/metadefs-concepts.html.

* `min_disk` - The minimum amount of disk space required to use the image.

* `min_ram` - The minimum amount of ram required to use the image.

* `original_image_name` - Specifies the parent image ID.

* `os_bit` - Specifies the OS architecture, 32 bit or 64 bit.

* `os_type` - Specifies the OS type. The value can be Linux, Windows, or Other.

* `os_version` - Specifies the OS version.

* `platform` - Specifies the image platform type. The value can be Windows, Ubuntu, Red Hat, SUSE, CentOS,
  Debian, OpenSUSE, Oracle Linux, Fedora, Other, CoreOS, or EulerOS.

* `properties` - Freeform information about the image.

* `protected` - Whether the image is protected.

* `schema` - The path to the JSON-schema that represent the image or image.

* `size_bytes` - The size of the image (in bytes).

* `status` - The image status.

* `support_disk_intensive` - Specifies whether the image supports disk-intensive ECSs.

* `support_high_performance` - Specifies whether the image supports high-performance ECSs.

* `support_kvm` - Specifies whether the image supports KVM.

* `support_kvm_gpu_type` - Specifies whether the image supports GPU-accelerated ECSs on the KVM platform.

* `support_kvm_infiniband` - Specifies whether the image supports ECSs with the InfiniBand NIC on the KVM platform.

* `support_large_memory` - Specifies whether the image supports large-memory ECSs.

* `support_xen` - Specifies whether the image supports Xen.

* `support_xen_gpu_type` - Specifies whether the image supports GPU-accelerated ECSs on the Xen platform.

* `support_xen_hana` - Specifies whether the image supports HANA ECSs on the Xen platform.

* `system_cmk_id` - Specifies the ID of the key used to encrypt the image.

* `tags` - See Argument Reference above.

* `virtual_env_type` - Specifies the environment where the image is used.
  The value can be `FusionCompute`, `Ironic`, `DataImage`, or `IsoImage`.

* `updated_at` - The date the image was modified.
