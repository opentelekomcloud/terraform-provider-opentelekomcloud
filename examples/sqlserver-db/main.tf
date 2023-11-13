resource "opentelekomcloud_networking_secgroup_v2" "secgroup_mssqlrds" {
  name        = var.secgroup_name
  description = "My neutron security group"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mssqlrds_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mssqlrds.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mssqlrds_dbport" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = var.db_port
  port_range_max    = var.db_port
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mssqlrds.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mssqlrds_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mssqlrds.id
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "${var.db_name}-instance"
  availability_zone = [var.availability_zone]

  db {
    password = var.db_passwd
    type     = var.db_type
    version  = var.db_version
    port     = var.db_port
  }

  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mssqlrds.id
  subnet_id         = var.network_id
  vpc_id            = var.vpc_id
  flavor            = var.db_flavor

  volume {
    type = "COMMON"
    size = 100
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  depends_on = [opentelekomcloud_networking_secgroup_v2.secgroup_mssqlrds]
}
