---
subcategory: "Image Management Service (IMS)"
---

# opentelekomcloud_images_member_v2

Manages a V2 Member resource within OpenTelekomCloud Glance.

## Example Usage

```hcl
resource "opentelekomcloud_images_member_v2" "member_1" {
  member   = var.project_to_share
  image_id = var.private_image
}
```

## Argument Reference

The following arguments are supported:

* `member` - (Required) Specifies the image member. The value is the project ID of a tenant.

* `image_id` - (Required) Specifies the image ID to be shared.

~> `status` and `vault_id` are valid parameters for Update operation only and
can't be used with the same user credentials. User should login with a project
credentials where the share was requested.

* `status` - (Optional) Specifies whether a shared image will be accepted or declined. Possible values:
  `accepted`, `rejected`.

* `vault_id` - (Optional) Specifies the ID of a vault. This parameter is mandatory if you want to accept
  a shared full-ECS image created from a CBR backup.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Specifies the image sharing status. After creation is `pending`.

* `created_at` - Specifies the time when a shared image was created. The value is in UTC format.

* `updated_at` - Specifies the time when a shared image was updated. The value is in UTC format.

* `schema` - Specifies the sharing schema.


