---
subcategory: "Object Storage Service (OBS)"
---

Up-to-date reference of API arguments for OBS bucket inventory you can get at
`https://docs.otc.t-systems.com/object-storage-service/api-ref/apis/advanced_bucket_settings/configuring_bucket_inventories.html`.

# opentelekomcloud_obs_bucket_inventory

Configures OBS bucket inventory resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "my-tf-test-bucket"
}

resource "opentelekomcloud_obs_bucket_inventory" "inventory" {
  bucket           = opentelekomcloud_obs_bucket.bucket.bucket
  configuration_id = "test-id"
  is_enabled       = true
  frequency        = "Weekly"
  destination {
    bucket = opentelekomcloud_obs_bucket.bucket.bucket
    format = "CSV"
    prefix = "test-"
  }
  filter_prefix            = "test-filter-prefix"
  included_object_versions = "Current"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew, String) Name of the bucket for saving inventories..

* `configuration_id` - (Required, ForceNew, String) ID of the inventory configuration. Valid characters: letters, digits, hyphens (-),
  periods (.) and underscores (_).

* `is_enabled` - (Required, Bool) Indicates whether the rule is enabled. If this parameter is set to `true`, the inventory is generated.

* `frequency` - (Required, String) Intervals when inventories are generated.
  An inventory is generated within one hour after it is configured for the first time. Then it is generated at the specified intervals.
  Possible values:
    * `Daily`
    * `Weekly`

* `destination` - (Required, List) Destination bucket settings of an inventory.
  The structure is documented below.

* `included_object_versions` - (Required, String) Indicates whether versions of objects are included in an inventory.
  Possible values:
    * `All`
    * `Current`

* `filter_prefix` - (Optional, String) Filtering by name prefix. Only objects with the specified name prefix are included in the inventory.

The `destination` block supports:

* `format` - (Required, String) Inventory format. Only the `CSV` format is supported.

* `bucket` - (Required, String) Name of the bucket for saving inventories.

* `prefix` - (Optional, String) The name prefix of inventory files. If no prefix is configured, the names of inventory files will start with the `BucketInventory` by default.

## Attributes Reference

The following attributes are exported

* `region` - Specifies the bucket region.

## Import

Inventories can be imported using related `bucket` and their `configuration_id` separated by the slashes, e.g.

```bash
$ terraform import opentelekomcloud_obs_bucket_inventory.inv <bucket>/<configuration_id>
```
