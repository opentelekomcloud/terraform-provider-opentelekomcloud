
resource "opentelekomcloud_dms_queue_v1" "queue_1" {
  name              = "${var.dms_name}"
  description       = "${var.dms_desc}"
  queue_mode        = "${var.queue_mode}"
  redrive_policy    = "${var.redrive_policy}"
  max_consume_count = "${var.max_consume_count}"
}
resource "opentelekomcloud_dms_queue_v1" "queue_2" {
  name = "${var.dms_name}_required"
}

resource "opentelekomcloud_dms_group_v1" "group_1" {
  name     = "${var.group_name}"
  queue_id = "${opentelekomcloud_dms_queue_v1.queue_1.id}"
}