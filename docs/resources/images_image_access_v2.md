---
subcategory: "Image Management Service (IMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_images_image_access_v2"
sidebar_current: "docs-opentelekomcloud-resource-images-image-access-v2"
description: |-
Manages an Image Sharing resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Image sharing you can get at
[documentation portal](https://docs.otc.t-systems.com/image-management-service/api-ref/native_openstack_apis/image_sharing_native_openstack_apis)

# opentelekomcloud_images_image_access_v2

Manages members for the shared OpenTelekomCloud Glance Image within the source project, which owns the Image.

## Example Usage

### Unprivileged user

Create a shared image and propose a membership to the `bed6b6cbb86a4e2d8dc2735c2f1000e4` project ID.

```hcl
resource "opentelekomcloud_images_image_v2" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}

resource "opentelekomcloud_images_image_access_v2" "rancheros_member" {
  image_id  = opentelekomcloud_images_image_v2.rancheros.id
  member_id = "bed6b6cbb86a4e2d8dc2735c2f1000e4"
}
```

### Privileged user

Create a shared image and set a membership to the `bed6b6cbb86a4e2d8dc2735c2f1000e4` project ID.

```hcl
resource "opentelekomcloud_images_image_v2" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}

resource "opentelekomcloud_images_image_access_v2" "rancheros_member" {
  image_id  = opentelekomcloud_images_image_v2.rancheros.id
  member_id = "bed6b6cbb86a4e2d8dc2735c2f1000e4"
  status    = "accepted"
}
```

## Argument Reference

The following arguments are supported:

* `member_id` - (Required) The member ID, e.g. the target project ID. Optional
  for admin accounts. Defaults to the current scope project ID.

* `image_id` - (Required) The proposed image ID.

* `status` - (Required) The member proposal status. Optional if admin wants to force the member
  proposal acceptance. Can either be `accepted`, `rejected` or `pending`. Defaults to
  `pending`. Forbidden for non-admin users.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Specifies the image sharing status. After creation is `pending`.

* `created_at` - Specifies the time when a shared image was created. The value is in UTC format.

* `updated_at` - Specifies the time when a shared image was updated. The value is in UTC format.

* `schema` - Specifies the sharing schema.

## Import

Image access can be imported using the `image_id` and the `member_id`, separated by a slash, e.g.

```
$ terraform import opentelekomcloud_images_image_access_v2 89c60255-9bd6-460c-822a-e2b959ede9d2/bed6b6cbb86a4e2d8dc2735c2f1000e4
```
