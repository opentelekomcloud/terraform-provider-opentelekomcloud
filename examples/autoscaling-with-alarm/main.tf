resource "opentelekomcloud_networking_router_v2" "router" {
  name             = "terraform"
  admin_state_up   = "true"
}

resource "opentelekomcloud_networking_network_v2" "network" {
  name           = "terraform"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet" {
  name            = "terraform"
  network_id      = "${opentelekomcloud_networking_network_v2.network.id}"
  cidr            = "172.16.10.0/24"
  ip_version      = 4
  dns_nameservers = ["100.125.1.250", "114.114.115.115"]
}

resource "opentelekomcloud_networking_router_interface_v2" "int_01" {
  router_id = "${opentelekomcloud_networking_router_v2.router.id}"
  subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet.id}"
}

resource "opentelekomcloud_compute_secgroup_v2" "secgroup" {
  name        = "terraform"
  description = "This is a terraform test security group"
  region      = "${var.region}"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 80
    to_port     = 80
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_as_group_v1" "as_group"{
  scaling_group_name = "terraform"
  scaling_configuration_id = "${opentelekomcloud_as_configuration_v1.as_configuration.id}"
  desire_instance_number = 2
  min_instance_number = 0
  max_instance_number = 3
  networks = [
    {id = "${opentelekomcloud_networking_network_v2.network.id}"},
  ]
  security_groups = [
    {id = "${opentelekomcloud_compute_secgroup_v2.secgroup.id}"},
  ]
  vpc_id = "${opentelekomcloud_networking_router_v2.router.id}"
  delete_publicip = true
  delete_instances = "yes" 
  depends_on = ["opentelekomcloud_networking_router_interface_v2.int_01"]
}

resource "opentelekomcloud_as_policy_v1" "as_policy"{
  scaling_policy_name = "terraform"
  scaling_group_id = "${opentelekomcloud_as_group_v1.as_group.id}"
  scaling_policy_type= "ALARM"
  scaling_policy_action = {
    operation = "ADD"
    instance_number = 1
  }
  alarm_id="${opentelekomcloud_ces_alarmrule.alarm_rule.id}"
}

resource "opentelekomcloud_as_configuration_v1" "as_configuration"{
  scaling_configuration_name = "terraform"
  instance_config = {
    flavor = "${var.flavor}"
    image = "${var.image_id}"
    disk = [
      {size = 40
      volume_type = "SATA"
      disk_type = "SYS"}
    ]
    key_name = "${var.keyname}"
    user_data = "${file("userdata.txt")}"
  }
}

resource "opentelekomcloud_ces_alarmrule" "alarm_rule" {
  "alarm_action_enabled" = "false" 
  "alarm_name" = "terraform"
  "metric" {
    "namespace" = "SYS.AS"
    "metric_name" = "cpu_util"
    "dimensions" {
        "name" = "AutoScalingGroup"
        "value" = "${opentelekomcloud_as_group_v1.as_group.id}"
    }
  }
  "condition"  {
    "period" = 300
    "filter" = "average"
    "comparison_operator" = ">"
    "value" = 80
    "unit" = "%"
    "count" = 2
  }
}
