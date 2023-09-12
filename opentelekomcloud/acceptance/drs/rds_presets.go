package drs

var RdsPreset = `
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}
resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {}

resource "opentelekomcloud_rds_instance_v3" "mysql_1" {
  name                = "RDS-ha"
  flavor              = "rds.mysql.s1.large.ha"
  ha_replication_mode = "semisync"
  vpc_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id   = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone   = ["eu-de-01,eu-de-03"]

  db {
    type     = "MySQL"
    version  = "5.7"
    password = "MySql_120521"
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

	public_ips = [opentelekomcloud_networking_floatingip_v2.fip_1.address]
}

resource "opentelekomcloud_rds_instance_v3" "mysql_2" {
  name                = "RDS-single"
  flavor              = "rds.mysql.m1.large"
  vpc_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id   = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone   = ["eu-de-01"]

  db {
    type     = "MySQL"
    version  = "5.7"
    password = "MySql_120521"
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

	public_ips = [opentelekomcloud_networking_floatingip_v2.fip_2.address]
}`
