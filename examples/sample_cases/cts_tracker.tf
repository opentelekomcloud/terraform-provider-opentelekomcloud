resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket"
  acl    = "public-read"
}

resource "opentelekomcloud_cts_tracker_v3" "tracker_v3" {
  bucket_name      = opentelekomcloud_s3_bucket.bucket.bucket
  file_prefix_name = "prefix"
  is_lts_enabled   = true
  status           = "enabled"
}

resource "opentelekomcloud_cts_event_notification_v3" "notification_v3" {
  notification_name = "my_notification"
  operation_type    = "complete"
  topic_id          = opentelekomcloud_smn_topic_v2.topic.id
  status            = "enabled"
}
