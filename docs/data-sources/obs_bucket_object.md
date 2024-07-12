---
subcategory: "Object Storage Service (OBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_obs_bucket_object"
sidebar_current: "docs-opentelekomcloud-datasource-obs-bucket-object"
description: |-
  Get details about OBS bucket object from OpenTelekomCloud
---

Up-to-date reference of API arguments for OBS bucket object you can get at
[documentation portal](https://docs.otc.t-systems.com/object-storage-service/api-ref/apis/operations_on_buckets/listing_objects_in_a_bucket.html)

# opentelekomcloud_obs_bucket_object

Use this data source to get details about the metadata and
_optionally_ (see below) content of an object stored inside OBS bucket.

-> **Note:** The content of an object (`body` field) is available only for objects which have a human-readable `Content-Type` (`text/*` and `application/json`).
  This is to prevent printing unsafe characters and potentially downloading large amount of data which would be thrown away in favour of metadata.

## Example Usage

```hcl
data "opentelekomcloud_obs_bucket_object" "hello" {
  bucket = "my-test-bucket"
  key    = "hello-world.zip"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to read the object from.

* `key` - (Required) The full path to the object inside the bucket.

* `version_id` - (Optional) Specific version ID of the object returned (defaults to latest version).

## Attributes Reference

The following attributes are exported:

* `body` - Object data (see **limitations above** to understand cases in which this field is actually available)

* `cache_control` - Specifies caching behavior along the request/reply chain.

* `content_disposition` - Specifies presentational information for the object.

* `content_encoding` - Specifies what content encodings have been applied to the
  object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.

* `content_language` - The language the content is in.

* `content_length` - Size of the body in bytes.

* `content_type` - A standard MIME type describing the format of the object data.

* `etag` - [ETag](https://en.wikipedia.org/wiki/HTTP_ETag) generated for the object (an MD5 sum of the object content in case it's not encrypted).

* `expiration` - If the object expiration is configured
  (see [object lifecycle management](http://docs.opentelekomcloud.amazon.com/AmazonS3/latest/dev/object-lifecycle-mgmt.html)), the field includes this header. It includes the expiry-date and rule-id key value pairs providing object expiration information. The value of the rule-id is URL encoded.

* `expires` - The date and time at which the object is no longer cacheable.

* `last_modified` - Last modified date of the object in RFC1123 format
  (e.g. `Mon, 02 Jan 2006 15:04:05 MST`)

* `metadata` - A map of metadata stored with the object in S3

* `storage_class` - Storage class information of the object. One of `STANDARD` (OBS Standard), `STANDARD_IA` (OBS Warm), or `GLACIER` (OBS Cold).

* `version_id` - The latest version ID of the object returned.

* `website_redirect_location` - If the bucket is configured as a website,
  redirects requests for this object to another object in the same bucket or to an external URL. OBS stores the value of this header in the object metadata.
