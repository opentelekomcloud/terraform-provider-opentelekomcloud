resource "opentelekomcloud_networking_secgroup_v2" "secgroup_mysqlrds" {
  name        = var.secgroup_name
  description = "My neutron security group"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mysqlrds_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mysqlrds_dbport" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = var.db_port
  port_range_max    = var.db_port
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mysqlrds_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id
}

resource "opentelekomcloud_vpc_v1" "this" {
  name = "test-vpc-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "this" {
  name       = "${opentelekomcloud_vpc_v1.this.name}-private"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.this.cidr, 8, 0)
  vpc_id     = opentelekomcloud_vpc_v1.this.id
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.this.cidr, 8, 0), 1)
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance"
  availability_zone = [var.availability_zone]
  db {
    password = var.db_passwd
    type     = var.db_type
    version  = var.db_version
    port     = var.db_port
  }
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id
  subnet_id         = opentelekomcloud_vpc_subnet_v1.this.id
  vpc_id            = opentelekomcloud_vpc_v1.this.id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = var.db_flavor
}
