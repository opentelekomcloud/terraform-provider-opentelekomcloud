### create random number

resource "random_id" "as" {
  byte_length = 4
}
### create keypair
resource "opentelekomcloud_compute_keypair_v2" "as-keypair" {
  name       = "${var.key_name_as}-${random_id.as.id}"
  public_key = file("public-2048.txt")
}

### Basic AS Configuration
resource "opentelekomcloud_as_configuration_v1" "as_config" {
  scaling_configuration_name = "as-config-basic"
  instance_config {
    flavor = var.flavor_id_as
    image  = var.image_id_as
    disk {
      size        = var.volume_size_as
      volume_type = var.volume_type_as ####  SATA (common I/O disk type) Or SSD (ultra-high I/O disk type).
      disk_type   = "SYS"    ####  DATA Or SYS
    }

    key_name = opentelekomcloud_compute_keypair_v2.as-keypair.name
    public_ip {
      eip {
        ip_type = "5_bgp"
        bandwidth {
          size          = 5
          share_type    = "PER"
          charging_mode = "traffic"
        }
      }
    }
  }
}

### used existing ecs
resource "opentelekomcloud_as_configuration_v1" "my_as_config" {
  scaling_configuration_name = "my_as_config_existingecs"
  instance_config {
    instance_id = "f36ea54d-b848-4db3-bf72-a182a70da1f4"
    key_name    = opentelekomcloud_compute_keypair_v2.as-keypair.name
  }
}
### as with user data and Metadata
resource "opentelekomcloud_as_configuration_v1" "as_config_1" {
  scaling_configuration_name = "my_as_config_userdata"
  instance_config {
    flavor = var.flavor_id_as
    image  = var.image_id_as
    disk {
      size        = var.volume_size_as
      volume_type = var.volume_type_as ####  SATA (common I/O disk type) Or SSD (ultra-high I/O disk type).
      disk_type   = "SYS"     ####  DATA Or SYS
    }
    key_name  = opentelekomcloud_compute_keypair_v2.as-keypair.name
    user_data = file("/opt/terraform/terraformTest/terraform-DT/modules/as/userdata.txt")
    metadata = {
      some_key = "some_value"
    }
  }
}

### create vpc

resource "opentelekomcloud_vpc_v1" "vpc_as" {
  name = "${var.vpc_name_as}-${random_id.as.id}"
  cidr = var.vpc_cidr_as
}
#create subnet1
resource "opentelekomcloud_vpc_subnet_v1" "subnet_as" {
  name          = var.subnet_name1_as
  cidr          = var.subnet_cidr1_as
  gateway_ip    = var.subnet_gateway_ip1_as
  vpc_id        = opentelekomcloud_vpc_v1.vpc_as.id
  primary_dns   = var.primary_dns_as
  secondary_dns = var.secondary_dns_as
}

### create security group
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_as" {
  name        = "${var.secgroup_name_as}-${random_id.as.id}"
  description = "Created By Terraform."
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_as" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "All"
  port_range_min    = 0
  port_range_max    = 0
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_as.id
}

resource "opentelekomcloud_elb_loadbalancer" "elb_as" {
  name           = "elb_as"
  type           = "External"
  vpc_id         = opentelekomcloud_vpc_v1.vpc_as.id
  admin_state_up = "true"
  #vip_subnet_id  = opentelekomcloud_vpc_subnet_v1.subnet_as.subnet_id

  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_as.id
  bandwidth         = "5"
}

resource "opentelekomcloud_elb_listener" "listener_as" {
  name             = "elb-listener-as"
  description      = "great listener"
  protocol         = "TCP"
  backend_protocol = "TCP"
  protocol_port    = 12345
  backend_port     = 8080
  lb_algorithm     = "roundrobin"
  loadbalancer_id  = opentelekomcloud_elb_loadbalancer.elb_as.id
}

###Basic Autoscaling Group
resource "opentelekomcloud_as_group_v1" "my_as_group" {
  scaling_group_name       = "as_group_basic"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  desire_instance_number   = 1
  min_instance_number      = 0
  max_instance_number      = 3
  cool_down_time           = 902
  networks        {
    id = opentelekomcloud_vpc_subnet_v1.subnet_as.id
  }
  security_groups  {
    id = opentelekomcloud_networking_secgroup_v2.secgroup_as.id
  }
  vpc_id           = opentelekomcloud_vpc_v1.vpc_as.id
  delete_publicip  = true
  delete_instances = "yes"
  available_zones  = [var.availability_zone_as]
  region           = var.region_as
  depends_on       = [opentelekomcloud_as_configuration_v1.as_config_1]
}


resource "opentelekomcloud_as_group_v1" "my_as_group2" {
  scaling_group_name       = "as_group_required"
  networks {
    id = opentelekomcloud_vpc_subnet_v1.subnet_as.id
  }
  security_groups {
    id = opentelekomcloud_networking_secgroup_v2.secgroup_as.id
  }
  vpc_id                   = opentelekomcloud_vpc_v1.vpc_as.id
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.as_config.id
  delete_instances         = "yes"
  delete_publicip          = true
}

### ADD AS Recurrence Policy

resource "opentelekomcloud_as_policy_v1" "hth_aspolicy" {
  scaling_policy_name = "hth_aspolicy_01"
  scaling_group_id    = opentelekomcloud_as_group_v1.my_as_group.id
  cool_down_time      = 900
  scaling_policy_type = "RECURRENCE"
  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time     = "07:00"
    recurrence_type = "Daily"
    start_time      = "2019-11-30T12:00Z"
    end_time        = "2019-12-30T12:00Z"
  }
}

### AS Scheduled Policy

resource "opentelekomcloud_as_policy_v1" "hth_aspolicy_1" {
  scaling_policy_name = "hth_aspolicy_02_${random_id.as.id}"
  scaling_group_id    = opentelekomcloud_as_group_v1.my_as_group.id
  cool_down_time      = 900
  scaling_policy_type = "SCHEDULED"
  scaling_policy_action  {
    operation       = "REMOVE"
    instance_number = 1
  }
  scheduled_policy  {
    launch_time = "2019-12-22T12:00Z"
  }
}

###  AS Alarm Policy

resource "opentelekomcloud_as_policy_v1" "hth_aspolicy_2" {
  scaling_policy_name = "hth_aspolicy_03_${random_id.as.id}"
  scaling_group_id    = opentelekomcloud_as_group_v1.my_as_group.id
  cool_down_time      = 900
  scaling_policy_type = "ALARM"
  alarm_id            = var.alarm_id
  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
}
