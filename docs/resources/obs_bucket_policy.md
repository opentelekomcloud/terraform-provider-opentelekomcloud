---
subcategory: "Object Storage Service (OBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_obs_bucket_policy"
sidebar_current: "docs-opentelekomcloud-resource-obs-bucket-policy"
description: |-
  Manages a OBS Bucket Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for OBS bucket policy you can get at
[documentation portal](https://docs.otc.t-systems.com/object-storage-service/api-ref/apis/advanced_bucket_settings)

# opentelekomcloud_obs_bucket_policy

Attaches a policy to an OBS bucket resource within OpenTelekomCloud.
Now respects HTTP_PROXY, HTTPS_PROXY environment variables.

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
      "GetObject",
      "PutObject"
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

## Import

Import an existing OBS bucket policy using the bucket name:

```shell
terraform import opentelekomcloud_obs_bucket_policy.policy bucket-name
```

Import reads the existing policy and restores both `bucket` and `policy` in state.
The credentials must have permission to read the bucket policy. Include the
matching bucket name and policy JSON in your configuration before applying.

A bucket must have an existing policy to import this resource. A private bucket
without a custom policy returns `NoSuchBucketPolicy` and has no policy resource
to import. WORM retention is separate and belongs to `opentelekomcloud_obs_bucket`.
