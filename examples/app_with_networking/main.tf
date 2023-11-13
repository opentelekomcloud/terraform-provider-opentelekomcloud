resource "opentelekomcloud_compute_keypair_v2" "terraform" {
  name       = "terraform"
  public_key = file("${var.ssh_key_file}.pub")
}

resource "opentelekomcloud_networking_network_v2" "terraform" {
  name           = "terraform"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "terraform" {
  name            = "terraform"
  network_id      = opentelekomcloud_networking_network_v2.terraform.id
  cidr            = "10.0.0.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "opentelekomcloud_networking_router_v2" "terraform" {
  name             = "terraform"
  admin_state_up   = "true"
}

resource "opentelekomcloud_networking_router_interface_v2" "terraform" {
  router_id = opentelekomcloud_networking_router_v2.terraform.id
  subnet_id = opentelekomcloud_networking_subnet_v2.terraform.id
}

resource "opentelekomcloud_networking_secgroup_v2" "terraform" {
  name        = "terraform"
  description = "Security group for the Terraform example instances"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.terraform.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_2" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.terraform.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_3" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.terraform.id
}

resource "opentelekomcloud_networking_floatingip_v2" "terraform" {
  depends_on = [opentelekomcloud_networking_router_interface_v2.terraform]
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "myip" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.terraform.address
  port_id     = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  admin_state_up = true
  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.terraform.id
    ip_address = "10.0.0.5"
  }
  network_id = opentelekomcloud_networking_network_v2.terraform.id
}

resource "opentelekomcloud_compute_instance_v2" "terraform" {
  name            = "terraform"
  image_name      = var.image
  flavor_id       = "s3.large.1"
  key_pair        = opentelekomcloud_compute_keypair_v2.terraform.name
  security_groups = [opentelekomcloud_networking_secgroup_v2.terraform.name]
  depends_on      = [ opentelekomcloud_networking_floatingip_associate_v2.myip ]

  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }

  connection {
    user        = "ubuntu"
    host        = opentelekomcloud_networking_floatingip_v2.terraform.address
    private_key = file(var.ssh_key_file)
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt -y update",
      "sudo apt -y install nginx",
      "sudo service nginx start",
    ]
  }
}
