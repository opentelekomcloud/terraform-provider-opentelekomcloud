#pass
resource "opentelekomcloud_s3_bucket" "bucket1" {
  region = "eu-de"
}
#pass
resource "opentelekomcloud_s3_bucket_object" "object" {

  bucket                 = opentelekomcloud_s3_bucket.bucket1.bucket
  key                    = "objectTest"
  source                 = "/opt/terraform/terraformTest/terraform-DT/modules/obs/file.txt"
  etag                   = md5(file("/opt/terraform/terraformTest/terraform-DT/modules/obs/file.txt"))
  cache_control          = "1"
  content_disposition    = "attachment;filename='fname.ext'"
  content_encoding       = "encoding"
  content_language       = "en-GB"
  website_redirect       = "/test2.txt"
  server_side_encryption = "aws:kms"
}

data "opentelekomcloud_s3_bucket_object" "b" {
  bucket = "obs-8f58"
  key    = "test-image-2.qcow2"
}

#pass
resource "opentelekomcloud_s3_bucket_object" "object2" {

  bucket = opentelekomcloud_s3_bucket.bucket1.bucket
  key    = "new_object_key"
  #source = "./policy.json"
  content                = "/opt/terraform/terraformTest/terraform-DT/modules/obs/file.txt"
  etag                   = md5(file("/opt/terraform/terraformTest/terraform-DT/modules/obs/file.txt"))
  cache_control          = "1"
  content_disposition    = "attachment;filename='fname.ext'"
  content_encoding       = "encoding"
  content_language       = "en-GB"
  website_redirect       = "/test2.txt"
  server_side_encryption = "aws:kms"
}


//pass
resource "opentelekomcloud_s3_bucket" "bucket2" {
  region = "eu-de"
  #bucket = var.bucket_name
  bucket_prefix = "test-"
  acl           = "private"
  #policy = file("/opt/terraform/terraformTest/terraform-DT/modules/obs/policy.json")
  tags {
    Name        = "Mybucket"
    Environment = "Dev"
  }
}
//pass
resource "opentelekomcloud_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
}
//pass
resource "opentelekomcloud_s3_bucket_policy" "policy" {
  bucket = opentelekomcloud_s3_bucket.b.bucket
  policy = <<POLICY
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
	      "AWS":["*"]
	  },
      "Action": [
	   "s3:DeleteBucket",
	   "s3:ListBucket",
	   "s3:ListBucketVersions"
	  ],
     "Resource": [
	    "arn:aws:s3:::my-tf-test-bucket",
	      "arn:aws:s3:::my-tf-test-bucket/*"
	  ]
    }
  ]
}
POLICY
}


##Static Website Hosting  pass
resource "opentelekomcloud_s3_bucket" "bucket3" {
  region = "eu-de"
  bucket = "s3-website-test-bucket3.hashicorp.com"
  acl    = "public-read"
  policy = file("/opt/terraform/terraformTest/terraform-DT/modules/obs/policy3.json")
  website {
    index_document = "index.html"
    error_document = "error.html"
    routing_rules  = <<EOF
	[{
    "Condition": {
        "KeyPrefixEquals": "docs/"
    },
    "Redirect": {
        "ReplaceKeyPrefixWith": "documents/"
    }
   }]
   EOF
  }
}

##Using CORS pass
resource "opentelekomcloud_s3_bucket" "bucket4" {
  region = "eu-de"
  bucket = "s3-website-test-bucket4.hashicorp.com"
  acl    = "public-read"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
##Using versioning pass
resource "opentelekomcloud_s3_bucket" "bucket5" {
  region = "eu-de"
  bucket = "my-tf-test-bucket-bucket5"
  acl    = "private"

  versioning {
    enabled = true
  }
}

##Enable Logging    pass
#resource "opentelekomcloud_s3_bucket" "log_bucket" {
#  region = "eu-de"
#  bucket = "my-tf-log-bucket-log-bucket"
#  acl    = "log-delivery-write"
#}
#resource "opentelekomcloud_s3_bucket" "bucket6" {
# region = "eu-de"
#  bucket = "my-tf-test-bucket-bucket6"
#  acl    = "private"
#  logging {
#    target_bucket = opentelekomcloud_s3_bucket.log_bucket.id
#    target_prefix = "log/"
#  }
#}

##Using object lifecycle  pass
resource "opentelekomcloud_s3_bucket" "bucket7" {
  region = "eu-de"
  bucket = "my-bucket-bucket6"
  acl    = "private"

  lifecycle_rule {
    id      = "log"
    enabled = true

    prefix = "log/"
    #tags {
    #  rule      = "log"
    #  autoclean = "true"
    #}

    expiration {
      days = 90
    }
  }

  lifecycle_rule {
    id      = "tmp"
    prefix  = "tmp/"
    enabled = true

    expiration {
      date = "2019-08-12"
    }
  }
}
//pass
resource "opentelekomcloud_s3_bucket" "versioning_bucket" {
  region = "eu-de"
  bucket = "my-versioning-bucket"
  acl    = "private"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    prefix  = "config/"
    enabled = true
    expiration {
      days = 90
    }
  }

}

