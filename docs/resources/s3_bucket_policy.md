---
subcategory: "Object Storage Service (OBS)"
---

# opentelekomcloud_s3_bucket_policy

Attaches a policy to an S3 bucket resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
}

resource "opentelekomcloud_s3_bucket_policy" "b" {
  bucket = opentelekomcloud_s3_bucket.b.id
  policy = <<POLICY
  {
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Sid": "IPAllow",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::${opentelekomcloud_s3_bucket.b.bucket}/*",
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
