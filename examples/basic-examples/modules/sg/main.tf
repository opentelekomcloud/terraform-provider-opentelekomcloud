#create random number
resource "random_id" "sg" {
  byte_length = 4
}
#create security group for ecs

resource "opentelekomcloud_networking_secgroup_v2" "secgroup1" {
 name        = "${var.secgroup_name}-${random_id.sg.id}"
 description = "Created By Terraform."
}
resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "All"
  port_range_min    = 0
  port_range_max    = 0
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup1.id}"
}

