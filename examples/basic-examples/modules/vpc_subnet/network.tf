resource "opentelekomcloud_networking_floatingip_v2" "floatip_1" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = "${opentelekomcloud_networking_network_v2.network_1.id}"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  name = "subnethhhaha"
  allocation_pools {
    start = "192.168.199.2"
	end = "192.168.199.251" 
  }
  gateway_ip = "192.168.199.1"
  #no_gateway = "true"
  dns_nameservers = ["8.8.8.8"]
 # host_routes= {
 #  destination_cidr = "10.0.0.0/29"  
 #  next_hop = "10.0.0.1"
 # }
}


resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "my_router"
  external_gateway = "0a2228f2-7f8a-45f1-8e09-9039e1d09975"
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = "${opentelekomcloud_networking_router_v2.router_1.id}"
  subnet_id = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
}

resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}


resource "opentelekomcloud_networking_port_v2" "port_1" {
  name               = "port_1"
  network_id         = "${opentelekomcloud_networking_network_v2.network_1.id}"
  admin_state_up     = "true"
  security_group_ids = ["${opentelekomcloud_compute_secgroup_v2.secgroup_1.id}"]
  device_owner = "neutron:VIP_PORT"  
  fixed_ip {
    "subnet_id"  = "${opentelekomcloud_networking_subnet_v2.subnet_1.id}"
    "ip_address" = "192.168.199.10"
  }
  allowed_address_pairs {
    "ip_address" = "192.168.199.10"
  }
}
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["${opentelekomcloud_compute_secgroup_v2.secgroup_1.name}"]
  network {
    port = "${opentelekomcloud_networking_port_v2.port_1.id}"
  }
}