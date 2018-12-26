resource "opentelekomcloud_networking_secgroup_v2" "secgroup_mysqlrds" {
  name        = "${var.secgroup_name}"
  description = "My neutron security group"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mysqlrds_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id}"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mysqlrds_dbport" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = "${var.db_port}"
  port_range_max    = "${var.db_port}"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id}"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_mysqlrds_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id}"
}


data "opentelekomcloud_rds_flavors_v1" "flavor_mysqlrds" {
  region = "${var.region}"
  datastore_name = "${var.db_type}"
  datastore_version = "${var.db_version}"
  speccode = "${var.db_flavor}"
}

resource "opentelekomcloud_rds_instance_v1" "instance_mysqlrds" {
  name = "${var.db_name}-instance"
  datastore {
    type = "${var.db_type}"
    version = "${var.db_version}"
  }
  flavorref = "${data.opentelekomcloud_rds_flavors_v1.flavor_mysqlrds.id}"
  volume {
    type = "COMMON"
    size = 100
  }
  region = "${var.region}"
  availabilityzone = "${var.availability_zone}"
  vpc = "${var.vpc_id}"
  nics {
    subnetid = "${var.existing_private_net_id}"
  }
  securitygroup {
    id = "${opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds.id}"
  }
  dbport = "${var.db_port}"
  backupstrategy = {
    starttime = "00:00:00"
    keepdays = 0
  }
  dbrtpd = "${var.db_passwd}"
  ha = {
    enable = true
    replicationmode = "async"
  }
  depends_on = ["opentelekomcloud_networking_secgroup_v2.secgroup_mysqlrds"]
}
