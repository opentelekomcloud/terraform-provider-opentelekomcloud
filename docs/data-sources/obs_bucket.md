---
subcategory: "Object Storage Service (OBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_obs_bucket"
sidebar_current: "docs-opentelekomcloud-datasource-obs-bucket"
description: |-
  Get details about OBS bucket from OpenTelekomCloud
---

Up-to-date reference of API arguments for OBS bucket you can get at
[documentation portal](https://docs.otc.t-systems.com/object-storage-service/api-ref/apis/operations_on_buckets/listing_buckets.html#obs-04-0020)

# opentelekomcloud_obs_bucket

Use this data source to get details about bucket within OpenTelekomCloud.


## Example Usage

```hcl
data "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "my-test-bucket"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to read.

## Attributes Reference

The following attributes are exported:

* `bucket_domain_name` - The bucket domain name. Will be of format `bucketname.obs.region.otc.t-systems.com`.

* `region` - The region this bucket resides in.

* `storage_class` - Specifies the storage class of the bucket. OBS provides three storage classes:
  `STANDARD`, `WARM` (Infrequent Access) and `COLD` (Archive).

* `tags` - A mapping of tags to assign to the bucket. Each tag is represented by one key-value pair.

* `logging` - A settings of bucket logging (documented below).

* `lifecycle_rule` - A configuration of object lifecycle management (documented below).

* `website` - A website object (documented below).

* `cors_rule` - A rule of Cross-Origin Resource Sharing (documented below).

* `lifecycle_rule` - A configuration of object lifecycle management (documented below).

* `server_side_encryption` - A configuration of server side encryption (documented below).

* `event_notifications` - A configuration of bucket event notifications (documented below).

The `logging` object supports the following:

* `target_bucket` - The name of the bucket that will receive the log objects.

* `target_prefix` - To specify a key prefix for log objects.

The `website` object supports the following:

* `index_document` - Specifies the default homepage of the static website, only HTML web pages are supported.

* `error_document` - Specifies the error page returned when an error occurs during static website access.

* `redirect_all_requests_to` - A hostname to redirect all website requests for this bucket to.

* `routing_rules` - A JSON or XML format containing routing rules describing redirect
  behavior and when redirects are applied.

The `lifecycle_rule` object supports the following:

* `name` - Unique identifier for lifecycle rules. The Rule Name contains a maximum of 255 characters.

* `enabled` - Specifies lifecycle rule status.

* `prefix` - Object key prefix identifying one or more objects to which the rule applies.

* `expiration` - Specifies a period when objects that have been last updated are automatically
  deleted. (documented below).

* `transition` - Specifies a period when objects that have been last updated are automatically
  transitioned to `WARM` or `COLD` storage class (documented below).

* `noncurrent_version_expiration` - Specifies a period when noncurrent object versions are
  automatically deleted. (documented below).

* `noncurrent_version_transition` - Specifies a period when noncurrent object versions are
  automatically transitioned to `WARM` or `COLD` storage class (documented below).

The `expiration` object supports the following

* `days` - Specifies the number of days when objects that have been last updated are automatically deleted.
  The expiration time must be greater than the transition times.

The `transition` object supports the following

* `days` - Specifies the number of days when objects that have been last updated are automatically
  transitioned to the specified storage class.

* `storage_class` - The class of storage used to store the object.

The `noncurrent_version_expiration` object supports the following

* `days` - Specifies the number of days when noncurrent object versions are automatically deleted.

The `noncurrent_version_transition` object supports the following

* `days` - Specifies the number of days when noncurrent object versions are automatically
  transitioned to the specified storage class.

* `storage_class` - The class of storage used to store the object.
