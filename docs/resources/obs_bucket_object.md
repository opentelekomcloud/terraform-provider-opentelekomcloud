---
subcategory: "Object Storage Service (OBS)"
---

Up-to-date reference of API arguments for OBS bucket object you can get at
`https://docs.otc.t-systems.com/object-storage-service/api-ref/apis/operations_on_objects`.

# opentelekomcloud_obs_bucket_object

Provides an OBS bucket object resource within OpenTelekomCloud.

## Example Usage

### Uploading to a bucket

```hcl
resource "opentelekomcloud_obs_bucket_object" "object" {
  bucket       = "your_bucket_name"
  key          = "new_key_from_content"
  content      = "some object content"
  content_type = "application/xml"
}
```

### Uploading a file to a bucket

```hcl
resource "opentelekomcloud_obs_bucket" "examplebucket" {
  bucket = "examplebuckettftest"
  acl    = "private"
}

resource "opentelekomcloud_obs_bucket_object" "object" {
  bucket = opentelekomcloud_obs_bucket.examplebucket.bucket
  key    = "new_key_from_file"
  source = "index.html"
}
```

### Server Side Encryption with OBS Default Master Key

```hcl
resource "opentelekomcloud_obs_bucket_object" "examplebucket_object" {
  bucket     = "your_bucket_name"
  key        = "someobject"
  source     = "index.html"
  encryption = true
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to put the file in.

* `key` - (Required) The name of the object once it is in the bucket.

* `source` - (Optional) The path to the source file being uploaded to the bucket.

* `content` - (Optional) The literal content being uploaded to the bucket.

* `acl` - (Optional) The ACL policy to apply. Defaults to `private`.

* `storage_class` - (Optioanl) Specifies the storage class of the object. Defaults to `STANDARD`.

* `content_type` - (Optional) A standard MIME type describing the format of the object data, e.g. application/octet-stream.
  All Valid MIME Types are valid for this input.

* `encryption` - (Optional) Whether enable server-side encryption of the object in SSE-KMS mode.

* `sse_kms_key_id` - (Optional) The ID of the kms key. If omitted, the default master key will be used.

* `etag` - (Optional) Specifies the unique identifier of the object content. It can be used to trigger updates.
  The only meaningful value is `md5(file("path_to_file"))`.

Either `source` or `content` must be provided to specify the bucket content.
These two arguments are mutually-exclusive.

## Attributes Reference

The following attributes are exported

* `id` - the `key` of the resource supplied above.

* `etag` - the ETag generated for the object (an MD5 sum of the object content).
  When the object is encrypted on the server side, the ETag value is not the MD5 value of the object,
  but the unique identifier calculated through the server-side encryption.

* `size` - the size of the object in bytes.

* `version_id` - A unique version ID value for the object, if bucket versioning is enabled.
