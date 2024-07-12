---
subcategory: "Image Management Service (IMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_images_image_access_accept_v2"
sidebar_current: "docs-opentelekomcloud-resource-images-image-access-accept-v2"
description: |-
  Manages an Image Sharing Accept resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Image sharing accept you can get at
[documentation portal](https://docs.otc.t-systems.com/image-management-service/api-ref/native_openstack_apis/image_sharing_native_openstack_apis)

# opentelekomcloud_images_image_access_accept_v2

Manages membership status for the shared OpenTelekomCloud Glance Image within the destination project, which has a member proposal.

## Example Usage

Accept a shared image membershipship proposal within the current project.

```hcl
data "opentelekomcloud_images_image_v2" "rancheros" {
  name       = "RancherOS"
  visibility = "shared"
}

resource "opentelekomcloud_images_image_access_accept_v2" "rancheros_member" {
  image_id  = data.opentelekomcloud_images_image_v2.rancheros.id
  member_id = "bed6b6cbb86a4e2d8dc2735c2f1000e4"
  status    = "accepted"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The proposed image ID.

* `member_id` - (Required) The member ID, e.g. the target project ID.

* `status` - (Required) The membership proposal status. Can either be `accepted`, `rejected` or `pending`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Specifies the time when a shared image was created. The value is in UTC format.

* `updated_at` - Specifies the time when a shared image was updated. The value is in UTC format.

* `schema` - Specifies the sharing schema.

## Import

Image access can be imported using the `image_id` and the `member_id`, separated by a slash, e.g.

```
$ terraform import opentelekomcloud_images_image_access_accept_v2 89c60255-9bd6-460c-822a-e2b959ede9d2/bed6b6cbb86a4e2d8dc2735c2f1000e4
```
