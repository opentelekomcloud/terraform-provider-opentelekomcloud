
#resource "opentelekomcloud_smn_topic_v2" "topic" {
#  name        = "topic_1"
#  display_name    = "The display name of topic by terraform"
#}

resource "opentelekomcloud_ces_alarmrule" "alarm_rule" {
  alarm_name        = var.alarm_name
  alarm_description = var.alarm_desc
  metric {
    #"namespace" = "SYS.ECS"
    namespace   = "SYS.AS"
    metric_name = "as"
    dimensions {
      name = "AutoScalingGroup"
      #"value" = opentelekomcloud_compute_instance_v2.webserver.id
      value = var.as_group_id
    }

  }
  condition {
    period              = 1
    filter              = "average"
    comparison_operator = "<"
    value               = 1
    unit                = "count"
    count               = 1
  }
  alarm_actions {
    type = "notification"
    notification_list = [
      "urn:smn:eu-de:b730519ca7064da2a3233e86bee139e4:smn-test"
    ]
  }
  ok_actions {
    type = "notification"
    notification_list = [
      "urn:smn:eu-de:b730519ca7064da2a3233e86bee139e4:smn-test"
    ]
  }
  alarm_enabled        = true
  alarm_action_enabled = true

}

resource "opentelekomcloud_ces_alarmrule" "alarm_rule2" {
  alarm_name = "${var.alarm_name}_2"
  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = var.ecs_id
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
  alarm_action_enabled = false
}
