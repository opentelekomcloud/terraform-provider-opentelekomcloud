---
subcategory: "Object Storage Service (OBS)"
---

Up-to-date reference of API arguments for OBS bucket policy you can get at
`https://docs.otc.t-systems.com/object-storage-service/api-ref/apis/advanced_bucket_settings`.

# opentelekomcloud_obs_bucket_policy

Attaches a policy to an OBS bucket resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "my-tf-test-bucket"
}

resource "opentelekomcloud_obs_bucket_policy" "policy" {
  bucket = opentelekomcloud_obs_bucket.bucket.id
  policy = <<POLICY
{
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "ID": ["*"]
    },
    "Action": [
      "ListBucket",
      "ListBucketVersions"
    ],
    "Resource": [
      "${opentelekomcloud_obs_bucket.bucket.bucket}/*"
    ]
  }]
}
POLICY
}
```

~>
  Please note that used policy syntax is OBS-specific. For s3-compatible policies check
  `opentelekomcloud_s3_bucket_policy` resource.

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which to apply the policy.

* `policy` - (Required) The text of the policy.
