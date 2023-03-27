---
subcategory: "Image Management Service (IMS)"
---

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

* `checksum` - The checksum of the data associated with the image.

* `created_at` - The date the image was created.

* `container_format`: The format of the image's container.

* `disk_format`: The format of the image's disk.

* `file` - the trailing path after the glance endpoint that represent the
  location of the image, or the path to retrieve it.

* `metadata` - The metadata associated with the image.
  Image metadata allow for meaningfully define the image properties
  and tags. See http://docs.openstack.org/developer/glance/metadefs-concepts.html.

* `min_disk_gb`: The minimum amount of disk space required to use the image.

* `min_ram_mb`: The minimum amount of ram required to use the image.

* `properties` - Freeform information about the image.

* `protected` - Whether the image is protected.

* `schema` - The path to the JSON-schema that represent the image or image

* `size_bytes` - The size of the image (in bytes).

* `tags` - See Argument Reference above.

* `update_at` - The date the image was last updated.
