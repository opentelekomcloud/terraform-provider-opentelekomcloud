resource "opentelekomcloud_vpc_v1" "this" {
  name = "test-vpc-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "this" {
  name       = "${opentelekomcloud_vpc_v1.this.name}-private"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.this.cidr, 8, 0)
  vpc_id     = opentelekomcloud_vpc_v1.this.id
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.this.cidr, 8, 0), 1)
  dns_list = [
    "1.1.1.1",
    "8.8.8.8",
  ]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  name            = "pool_1"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb_1.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "TCP"

  session_persistence {
    type                = "SOURCE_IP"
    persistence_timeout = "30"
  }
}

resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  subnet_id   = opentelekomcloud_vpc_subnet_v1.this.subnet_id
  network_ids = [opentelekomcloud_vpc_subnet_v1.this.network_id]

  availability_zones = [var.az]
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "listener_1"
  description     = "some interesting description"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb_1.id
  protocol        = "HTTP"
  protocol_port   = 8080

  advanced_forwarding = true
  sni_match_algo      = "wildcard"

  insert_headers {
    forwarded_host = true
  }

  ip_group {
    id     = opentelekomcloud_lb_ipgroup_v3.group_1.id
    enable = true
  }
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "group description"

  ip_list {
    ip          = "192.168.0.10"
    description = "one"
  }
}

resource "opentelekomcloud_as_group_v1" "group_example" {
  scaling_group_name       = "group_example"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.config_example.id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 3
  networks {
    id = opentelekomcloud_vpc_subnet_v1.this.id
  }
  security_groups {
    id = var.security_group_id
  }
  vpc_id           = opentelekomcloud_vpc_v1.this.id
  lbaas_listeners   {
    pool_id = opentelekomcloud_lb_pool_v3.pool.id
    protocol_port = opentelekomcloud_lb_listener_v3.listener_1.protocol_port
  }
  delete_publicip  = true
  delete_instances = "yes"
}

resource "opentelekomcloud_as_policy_v1" "policy_example" {
  scaling_policy_name = "policy_example"
  scaling_group_id    = opentelekomcloud_as_group_v1.group_example.id
  scaling_policy_type = "RECURRENCE"
  scaling_policy_action {
    operation = "ADD"
  }
  scheduled_policy {
    launch_time     = "07:00"
    recurrence_type = "Daily"
    end_time        = "2024-12-30T12:00Z"
  }
}

resource "opentelekomcloud_as_configuration_v1" "config_example" {
  scaling_configuration_name = "config_example"
  instance_config {
    flavor = var.flavor
    image  = var.image_id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name  = var.keyname
    user_data = file("userdata.txt")
  }
}
