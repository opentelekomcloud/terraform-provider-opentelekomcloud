---
subcategory: "Object Storage Service (OBS)"
---

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
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Sid": "IPAllow",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::${opentelekomcloud_obs_bucket.bucket.id}/*",
      "Condition": {
         "IpAddress": {"aws:SourceIp": "8.8.8.8/32"}
      }
    }
  ]}
POLICY
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which to apply the policy.

* `policy` - (Required) The text of the policy.
