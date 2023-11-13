# MySQL server instance configuration

This example will show how to provisioning a mysql instance on OpenTelekomCloud.
For more detailed parameters description, please refer to the
[doc](https://www.terraform.io/docs/providers/opentelekomcloud/index.html).

The ```main.tf``` contains the major scripts, which will open the major ports
via security group rules and create mysql instance as following show:

```hcl
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


data "opentelekomcloud_rds_flavors_v1" "flavor_mssqlrds" {
  region            = var.region
  datastore_name    = var.db_type
  datastore_version = var.db_version
  speccode          = var.db_flavor
}

resource "opentelekomcloud_rds_instance_v1" "instance_mssqlrds" {
  name = "${var.db_name}-instance"
  datastore {
    type    = var.db_type
    version = var.db_version
  }
  flavorref        = data.opentelekomcloud_rds_flavors_v1.flavor_mssqlrds.id
  region           = var.region
  availabilityzone = var.availability_zone
  vpc              = var.vpc_id
  volume {
    type = "COMMON"
    size = 100
  }
  nics {
    subnetid = var.existing_private_net_id
  }
  securitygroup {
    id = opentelekomcloud_networking_secgroup_v2.secgroup_mssqlrds.id
  }
  dbport = var.db_port
  backupstrategy = {
    starttime = "00:00:00"
    keepdays  = 0
  }
  dbrtpd     = var.db_passwd
  depends_on = [opentelekomcloud_networking_secgroup_v2.secgroup_mssqlrds]
}
```

Note: Please do not forget to change the ```<YOUR_XXX>``` tag with your actual
values.
