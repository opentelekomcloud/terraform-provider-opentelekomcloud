module "network" {
  source = "../modules/network"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = module.network.default_security_group_id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_2" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = module.network.default_security_group_id
}

resource "opentelekomcloud_as_group_v1" "as_group" {
  scaling_group_name       = "terraform"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_configuration.id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 3
  networks {
    id = module.network.shared_subnet.network_id
  }

  security_groups {
    id = module.network.default_security_group_id
  }
  vpc_id           = module.network.shared_subnet.vpc_id
  delete_publicip  = true
  delete_instances = "yes"
}

resource "opentelekomcloud_as_policy_v1" "as_policy" {
  scaling_policy_name = "terraform"
  scaling_group_id    = opentelekomcloud_as_group_v1.as_group.id
  scaling_policy_type = "ALARM"
  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  alarm_id = opentelekomcloud_ces_alarmrule.alarm_rule.id
}

data "opentelekomcloud_images_image_v2" "latest_image" {
  name        = var.image_name
  most_recent = true
}

resource "opentelekomcloud_as_configuration_v1" "as_configuration" {
  scaling_configuration_name = "terraform"
  instance_config  {
    flavor = var.flavor
    image  = data.opentelekomcloud_images_image_v2.latest_image.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type = "SYS"
    }

    key_name  = var.keyname
    user_data = file("userdata.txt")
  }
}

resource "opentelekomcloud_ces_alarmrule" "alarm_rule" {
  alarm_action_enabled = "false"
  alarm_name           = "terraform"
  metric {
    namespace   = "SYS.AS"
    metric_name = "cpu_util"
    dimensions {
      name  = "AutoScalingGroup"
      value = opentelekomcloud_as_group_v1.as_group.id
    }
  }
  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%"
    count               = 2
  }
}
