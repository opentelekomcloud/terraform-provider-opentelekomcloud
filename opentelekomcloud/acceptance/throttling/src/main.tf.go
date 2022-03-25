package src

const Main = `
###########
# VPC part
###########
resource "opentelekomcloud_vpc_v1" "vpc" {
  name   = var.environment
  cidr   = var.vpc_cidr
  shared = true
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  name          = var.environment
  vpc_id        = opentelekomcloud_vpc_v1.vpc.id
  cidr          = var.subnet_cidr
  gateway_ip    = var.subnet_gateway_ip
  primary_dns   = var.subnet_primary_dns
  secondary_dns = var.subnet_secondary_dns
}

###########
# ELB part
###########
resource "opentelekomcloud_networking_floatingip_v2" "eip" {
  pool    = "admin_external_net"
  port_id = opentelekomcloud_lb_loadbalancer_v2.lb.vip_port_id
}

resource "opentelekomcloud_lb_loadbalancer_v2" "lb" {
  name          = "${var.environment}-lb"
  vip_subnet_id = opentelekomcloud_vpc_subnet_v1.subnet.subnet_id
}

resource "opentelekomcloud_lb_listener_v2" "listener_80" {
  protocol         = "TCP"
  name             = "${var.environment}-listener_80"
  protocol_port    = 80
  loadbalancer_id  = opentelekomcloud_lb_loadbalancer_v2.lb.id
}

resource "opentelekomcloud_lb_listener_v2" "listener_443" {
  protocol         = "TCP"
  name             = "${var.environment}-listener_443"
  protocol_port    = 443
  loadbalancer_id  = opentelekomcloud_lb_loadbalancer_v2.lb.id
}

resource "opentelekomcloud_lb_listener_v2" "listener_6443" {
  protocol         = "TCP"
  name             = "${var.environment}-listener_6443"
  protocol_port    = 6443
  loadbalancer_id  = opentelekomcloud_lb_loadbalancer_v2.lb.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_80" {
  protocol    = "TCP"
  name        = "${var.environment}-pool_80"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_80.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_443" {
  protocol    = "TCP"
  name        = "${var.environment}-pool_443"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_443.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_6443" {
  protocol    = "TCP"
  name        = "${var.environment}-pool_6443"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_6443.id
}

resource "opentelekomcloud_lb_monitor_v2" "monitor_80" {
  pool_id        = opentelekomcloud_lb_pool_v2.pool_80.id
  type           = "HTTP"
  delay          = 10
  timeout        = 5
  max_retries    = 10
  domain_name    = "${var.rancher_host}.${var.rancher_domain}"
  url_path       = "/healthz"
  http_method    = "GET"
  expected_codes = "200"
  monitor_port   = 80
}

resource "opentelekomcloud_lb_monitor_v2" "monitor_443" {
  pool_id        = opentelekomcloud_lb_pool_v2.pool_443.id
  type           = "HTTP"
  delay          = 10
  timeout        = 5
  max_retries    = 10
  domain_name    = "${var.rancher_host}.${var.rancher_domain}"
  url_path       = "/healthz"
  http_method    = "GET"
  expected_codes = "200"
  monitor_port   = 80
}

resource "opentelekomcloud_lb_monitor_v2" "monitor_6443" {
  pool_id        = opentelekomcloud_lb_pool_v2.pool_6443.id
  type           = "HTTP"
  delay          = 10
  timeout        = 5
  max_retries    = 10
  domain_name    = "${var.rancher_host}.${var.rancher_domain}"
  url_path       = "/healthz"
  http_method    = "GET"
  expected_codes = "200"
  monitor_port   = 80
}

# server 1
resource "opentelekomcloud_lb_member_v2" "member_80_1" {
  address       = opentelekomcloud_compute_instance_v2.throttle-server-1.access_ip_v4
  protocol_port = 80
  pool_id       = opentelekomcloud_lb_pool_v2.pool_80.id
  subnet_id     = opentelekomcloud_vpc_subnet_v1.subnet.subnet_id
}

resource "opentelekomcloud_lb_member_v2" "member_443_1" {
  address       = opentelekomcloud_compute_instance_v2.throttle-server-1.access_ip_v4
  protocol_port = 443
  pool_id       = opentelekomcloud_lb_pool_v2.pool_443.id
  subnet_id     = opentelekomcloud_vpc_subnet_v1.subnet.subnet_id
}

resource "opentelekomcloud_lb_member_v2" "member_6443_1" {
  address       = opentelekomcloud_compute_instance_v2.throttle-server-1.access_ip_v4
  protocol_port = 6443
  pool_id       = opentelekomcloud_lb_pool_v2.pool_6443.id
  subnet_id     = opentelekomcloud_vpc_subnet_v1.subnet.subnet_id
}

# server 2
resource "opentelekomcloud_lb_member_v2" "member_80_2" {
  address       = opentelekomcloud_compute_instance_v2.throttle-server-2.access_ip_v4
  protocol_port = 80
  pool_id       = opentelekomcloud_lb_pool_v2.pool_80.id
  subnet_id     = opentelekomcloud_vpc_subnet_v1.subnet.subnet_id
}

resource "opentelekomcloud_lb_member_v2" "member_443_2" {
  address       = opentelekomcloud_compute_instance_v2.throttle-server-2.access_ip_v4
  protocol_port = 443
  pool_id       = opentelekomcloud_lb_pool_v2.pool_443.id
  subnet_id     = opentelekomcloud_vpc_subnet_v1.subnet.subnet_id
}

resource "opentelekomcloud_lb_member_v2" "member_6443_2" {
  address       = opentelekomcloud_compute_instance_v2.throttle-server-2.access_ip_v4
  protocol_port = 6443
  pool_id       = opentelekomcloud_lb_pool_v2.pool_6443.id
  subnet_id     = opentelekomcloud_vpc_subnet_v1.subnet.subnet_id
}

###########
# RDS part
###########
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name        = "${var.environment}-rds-secgroup"
  description = "terraform security group rds"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = var.rds_port
  port_range_max    = var.rds_port
  remote_ip_prefix  = var.subnet_cidr
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
}

resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
  name        = "${var.environment}-rds-pg"
  description = "Parameter Group for ${var.environment}-rds RDS"
  values = {
    local_infile                         = "OFF"
    max_user_connections                 = "1000"
    validate_password_length             = "10"
    validate_password_number_count       = "1"
    validate_password_special_char_count = "1"
  }
  datastore {
    type    = "mysql"
    version = var.rds_version
  }
}

resource "opentelekomcloud_rds_instance_v3" "rds" {
  availability_zone = var.rds_az
  db {
    password = var.rds_root_password
    type     = var.rds_type
    version  = var.rds_version
    port     = var.rds_port
  }
  name              = "${var.environment}-rds"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  subnet_id         = opentelekomcloud_vpc_subnet_v1.subnet.id
  vpc_id            = opentelekomcloud_vpc_v1.vpc.id
  volume {
    type = var.rds_volume_type
    size = var.rds_volume_size
  }
  ha_replication_mode = var.rds_ha_mode
  param_group_id      = opentelekomcloud_rds_parametergroup_v3.pg.id
  flavor              = var.rds_flavor
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 30
  }
}

###########
# DNS part
###########
resource "opentelekomcloud_dns_zone_v2" "dns" {
  count       = var.create_dns ? 1 : 0
  name        = "${var.rancher_domain}."
  email       = var.admin_email
  description = "tf managed zone"
  ttl         = 300
  type        = "public"
}

resource "opentelekomcloud_dns_recordset_v2" "public_record" {
  count       = var.create_dns ? 1 : 0
  zone_id     = opentelekomcloud_dns_zone_v2.dns[0].id
  name        = "${var.rancher_host}.${var.rancher_domain}."
  description = "tf managed zone"
  type        = "A"
  ttl         = 300
  records     = [ opentelekomcloud_networking_floatingip_v2.eip.address ]
}

###########
# ECS part
###########
locals {
  throttle_server = "throttle_server"

  throttle_node = "throttle_node"

  wireguard = "wireguard"
}

data "opentelekomcloud_images_image_v2" "image-1" {
  name        = var.image_name_server-1
  most_recent = true
}

data "opentelekomcloud_images_image_v2" "image-2" {
  name        = var.image_name_server-2
  most_recent = true
}


# Secgroup part (ECS)
resource "opentelekomcloud_networking_secgroup_v2" "throttle-server-secgroup" {
  description = "throttle Server Group"
  name        = "${var.environment}-secgroup"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_all_out" {
  description       = "Rancher/throttle accept all traffic"
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_22_in_self" {
  count             = var.deploy_wireguard ? 1 : 0
  description       = "throttle Server ssh internal"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_group_id   = opentelekomcloud_networking_secgroup_v2.wireguard-secgroup[0].id
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_80_in" {
  description       = "Rancher HTTP ELB network"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "100.125.0.0/16"
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_443_in" {
  description       = "Rancher HTTPS ELB network"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "100.125.0.0/16"
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_443_in_self" {
  description       = "Rancher HTTPS internal"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_group_id   = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_6443_in" {
  description       = "Kube API ELB network"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_ip_prefix  = "100.125.0.0/16"
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_6443_in_self" {
  description       = "Kube API internal"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_group_id   = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_8472_in" {
  description       = "Flannel VXLAN"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 8472
  port_range_max    = 8472
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_9796_in" {
  description       = "Prometheus Node Exporter"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9796
  port_range_max    = 9796
  remote_group_id   = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_throttle_10250_in" {
  description       = "Kubelet"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_group_id   = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
  security_group_id = opentelekomcloud_networking_secgroup_v2.throttle-server-secgroup.id
}

resource "opentelekomcloud_networking_secgroup_v2" "wireguard-secgroup" {
  count       = var.deploy_wireguard ? 1 : 0
  description = "throttle Wireguard Server"
  name        = "${var.environment}-secgroup-wg"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_wg_all_out" {
  count             = var.deploy_wireguard ? 1 : 0
  description       = "throttle Wireguard accept all traffic"
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.wireguard-secgroup[0].id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_wg_in" {
  count             = var.deploy_wireguard ? 1 : 0
  description       = "throttle Wireguard"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = var.wg_server_port
  port_range_max    = var.wg_server_port
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.wireguard-secgroup[0].id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "sg_icmp_in" {
  count             = var.deploy_wireguard ? 1 : 0
  description       = "throttle Wireguard"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.wireguard-secgroup[0].id
}

# ssh key part
resource "opentelekomcloud_compute_keypair_v2" "throttle-server-key" {
  name       = "${var.environment}-key"
  public_key = var.public_key
}

# ECS part (instances)
resource "opentelekomcloud_compute_instance_v2" "throttle-server-1" {
  name              = "${var.environment}-server-1"
  availability_zone = var.availability_zone1
  flavor_id         = var.flavor_id
  key_pair          = opentelekomcloud_compute_keypair_v2.throttle-server-key.name
  security_groups   = ["${var.environment}-secgroup"]
  user_data         = local.throttle_server
  power_state       = var.power_state
  network {
    uuid = opentelekomcloud_vpc_subnet_v1.subnet.id
  }
  block_device {
    boot_index            = 0
    source_type           = "image"
    destination_type      = "volume"
    uuid                  = data.opentelekomcloud_images_image_v2.image-1.id
    delete_on_termination = true
    volume_size           = 30
  }
}

resource "opentelekomcloud_compute_instance_v2" "throttle-server-2" {
  name              = "${var.environment}-server-2"
  availability_zone = var.availability_zone2
  flavor_id         = var.flavor_id
  key_pair          = opentelekomcloud_compute_keypair_v2.throttle-server-key.name
  security_groups   = ["${var.environment}-secgroup"]
  user_data         = local.throttle_node
  power_state       = var.power_state
  network {
    uuid = opentelekomcloud_vpc_subnet_v1.subnet.id
  }
  block_device {
    boot_index            = 0
    source_type           = "image"
    destination_type      = "volume"
    uuid                  = data.opentelekomcloud_images_image_v2.image-2.id
    delete_on_termination = true
    volume_size           = 30
  }
}

resource "opentelekomcloud_compute_instance_v2" "wireguard" {
  count             = var.deploy_wireguard ? 1 : 0
  name              = "${var.environment}-wireguard"
  availability_zone = var.availability_zone1
  flavor_id         = var.flavor_id
  key_pair          = opentelekomcloud_compute_keypair_v2.throttle-server-key.name
  security_groups   = ["${var.environment}-secgroup-wg"]
  user_data         = local.wireguard
  power_state       = var.power_state
  network {
    uuid = opentelekomcloud_vpc_subnet_v1.subnet.id
  }
  block_device {
    boot_index            = 0
    source_type           = "image"
    destination_type      = "volume"
    uuid                  = data.opentelekomcloud_images_image_v2.image-1.id
    delete_on_termination = true
    volume_size           = 30
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "wireguard" {
  count = var.deploy_wireguard ? 1 : 0
  pool  = "admin_external_net"
}

resource "opentelekomcloud_compute_floatingip_associate_v2" "wireguard" {
  count       = var.deploy_wireguard ? 1 : 0
  floating_ip = opentelekomcloud_networking_floatingip_v2.wireguard[0].address
  instance_id = opentelekomcloud_compute_instance_v2.wireguard[0].id
}
`
