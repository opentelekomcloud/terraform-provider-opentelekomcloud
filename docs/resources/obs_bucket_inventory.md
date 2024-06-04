---
subcategory: "Object Storage Service (OBS)"
---

Up-to-date reference of API arguments for OBS bucket inventory you can get at
`https://docs.otc.t-systems.com/object-storage-service/api-ref/apis/advanced_bucket_settings`.

# opentelekomcloud_obs_bucket_inventory

Provides an OBS bucket inventory resource within OpenTelekomCloud.
Now respects HTTP_PROXY, HTTPS_PROXY environment variables.

## Example Usage

### Private Bucket with Tags

```hcl
variable "bucket_name" {}

resource "opentelekomcloud_obs_bucket_inventory" "inventory" {
  bucket           = var.bucket_name
  configuration_id = "test-configuration"
  is_enabled       = true
  frequency        = "Weekly"
  destination {
    bucket = var.bucket_name
    format = "CSV"
    prefix = "test-"
  }
  filter_prefix            = "test-filter-prefix"
  included_object_versions = "Current"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew, String) The name of the bucket to which to apply the inventory.

* `configuration_id` - (Required, ForceNew, String) ID of the inventory configuration. Valid characters:
  letters, digits, hyphens (-), periods (.) and underscores (_)

* `is_enabled` - (Required, Bool) Indicates whether the rule is enabled. If this parameter is set to `true`,
  the inventory is generated. If not, the inventory will not be generated.

* `frequency` - (Required, ForceNew, String) Intervals when inventories are generated.
  This resource provides the following frequency options:
    - `Daily`
    - `Weekly`

* `destination` - (Required, List) Specifies the destination bucket inventory configuration.
  The [object](#bucket_inventory_destination) structure is documented below.

* `included_object_versions` - (Required, String) Indicates whether versions of objects are included in an inventory.
  This resource provides the following versions options:
    - `All`
    - `Current`

* `filter_prefix` - (Optional, String) Option to filter by name prefix.
  Only objects with the specified name prefix are included in the inventory.

<a name="bucket_inventory_destination"></a>
The `destination` block supports:

* `format` - (Required, String) Inventory format. Currently only the `CSV` format is supported.

* `bucket` - (Required, String) Name of the bucket for saving inventories.

* `prefix` - (Optional, String) The name prefix of inventory files. If no prefix is configured,
  the names of inventory files will start with the `BucketInventory` by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the inventory configuration.

* `region` - The region where bucket is located.

## Import

OBS bucket can be imported using their `bucket_name` and the `configuration_id`, separated by a slash, e.g.

```shell
terraform import opentelekomcloud_obs_bucket_inventory.inventory test-bucket/test-config
```
