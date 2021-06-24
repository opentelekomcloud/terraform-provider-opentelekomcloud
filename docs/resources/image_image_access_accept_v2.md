---
subcategory: "Image Management Service (IMS)"
---

# opentelekomcloud_images_image_access_v2

Manages membership status for the shared OpenTelekomCloud Glance Image within the destination project, which has a member proposal.

## Example Usage

Accept a shared image membershipship proposal within the current project.

```hcl
data "opentelekomcloud_images_image_v2" "rancheros" {
  name       = "RancherOS"
  visibility = "shared"
}

resource "opentelekomcloud_images_image_access_accept_v2" "rancheros_member" {
  image_id = data.opentelekomcloud_images_image_v2.rancheros.id
  status   = "accepted"
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The proposed image ID.

* `member_id` - (Optional) The member ID, e.g. the target project ID. Optional
  for admin accounts. Defaults to the current scope project ID.

* `status` - (Required) The membership proposal status. Can either be `accepted`, `rejected` or `pending`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Specifies the time when a shared image was created. The value is in UTC format.

* `updated_at` - Specifies the time when a shared image was updated. The value is in UTC format.

* `schema` - Specifies the sharing schema.

## Import

Image access acceptance status can be imported using the `image_id`, e.g.

```
$ terraform import opentelekomcloud_images_image_access_accept_v2 89c60255-9bd6-460c-822a-e2b959ede9d2
```
