
resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
  bucket_name               = var.bucket_name
  file_prefix_name          = "yO8Q"
  is_support_smn            = true
  topic_id                  = "urn:smn:eu-de:b730519ca7064da2a3233e86bee139e4:smn-test"
  is_send_all_key_operation = false
  operations                = ["login", "create"]
  need_notify_user_list     = ["user1"]

}


#data "opentelekomcloud_cts_tracker_v1" "tracker" {
# bucket_name = opentelekomcloud_cts_tracker_v1.tracker_v1.bucket_name
# }
