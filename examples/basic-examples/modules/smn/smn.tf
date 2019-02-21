resource "opentelekomcloud_smn_topic_v2" "topic_1" {
  name        = "${var.topic_name}"
  display_name    = "${var.display_name}"
  # topic_urn    =  "urn:smn:eu-de:b730519ca7064da2a3233e86bee139e4:${var.topic_name}"
  push_policy  =  0
 # create_time  = 
 # update_time  = 
}
resource "opentelekomcloud_smn_topic_v2" "topic_2" {
  name        = "${var.topic_name}_2"
}
resource "opentelekomcloud_smn_subscription_v2" "subscription_1" {
  topic_urn       = "${opentelekomcloud_smn_topic_v2.topic_1.id}"
  endpoint        = "${var.subscription_endpoint}"
  protocol        = "${var.subscription_protocol}"
  remark          = "${var.subscription_remark}"
}
