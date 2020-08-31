resource "opentelekomcloud_fw_rule_v2" "rule_1" {
  protocol = "${var.rule_protocol}"
  action   = "${var.rule_action}"
}

resource "opentelekomcloud_fw_rule_v2" "rule_2" {
  name             = "${var.rule_name}"
  description      = "${var.rule_desc}"
  action           = "deny"
  protocol         = "udp"
  ip_version       = 4
  destination_port = "23"
  enabled          = "true"
}

resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name        = "${var.policy_name}"
  description = "${var.policy_desc}"
  rules = ["${opentelekomcloud_fw_rule_v2.rule_1.id}",
    "${opentelekomcloud_fw_rule_v2.rule_2.id}",
  ]
  audited = "false"
  shared  = "false"
}
#resource "opentelekomcloud_networking_network_v2" "network_1" {
#  name           = "network_1"
#  admin_state_up = "true"
#}

#resource "opentelekomcloud_networking_port_v2" "port_1" {
#  name           = "port_1"
#  network_id     = "${opentelekomcloud_networking_network_v2.network_1.id}"
#  admin_state_up = "true"
#}

resource "opentelekomcloud_fw_firewall_group_v2" "firewall_group_1" {
  name              = "${var.firewall_group_name}"
  ingress_policy_id = "${opentelekomcloud_fw_policy_v2.policy_1.id}"
  egress_policy_id  = "${opentelekomcloud_fw_policy_v2.policy_1.id}"
  description       = "${var.firewall_group_desc}"
  #ports = ["${opentelekomcloud_networking_port_v2.port_1.id}"]
}
