resource "opentelekomcloud_ces_alarmrule" "alarm_rule" {
  alarm_name = "alarm_rule"
  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = opentelekomcloud_compute_instance_v2.webserver.id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }
  alarm_actions {
    type = "notification"
    notification_list = [
      opentelekomcloud_smn_topic_v2.topic.id
    ]
  }
}
