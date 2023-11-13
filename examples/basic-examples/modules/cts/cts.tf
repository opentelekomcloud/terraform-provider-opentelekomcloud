resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = opentelekomcloud_obs_bucket.bucket.bucket
  file_prefix_name = "yO8Q1"
  is_lts_enabled   = true
}

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = var.bucket_name
  acl           = "public-read"
  force_destroy = true
}
