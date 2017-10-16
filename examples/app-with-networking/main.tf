resource "opentelekomcloud_compute_keypair_v2" "terraform" {
  name       = "terraform"
  public_key = "${file("${var.ssh_key_file}.pub")}"
}

resource "opentelekomcloud_networking_network_v2" "terraform" {
  name           = "terraform"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "terraform" {
  name            = "terraform"
  network_id      = "${opentelekomcloud_networking_network_v2.terraform.id}"
  cidr            = "10.0.0.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "opentelekomcloud_networking_router_v2" "terraform" {
  name             = "terraform"
  admin_state_up   = "true"
  external_gateway = "${var.external_gateway}"
}

resource "opentelekomcloud_networking_router_interface_v2" "terraform" {
  router_id = "${opentelekomcloud_networking_router_v2.terraform.id}"
  subnet_id = "${opentelekomcloud_networking_subnet_v2.terraform.id}"
}

resource "opentelekomcloud_compute_secgroup_v2" "terraform" {
  name        = "terraform"
  description = "Security group for the Terraform example instances"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 80
    to_port     = 80
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_floatingip_v2" "terraform" {
  pool       = "${var.pool}"
  depends_on = ["opentelekomcloud_networking_router_interface_v2.terraform"]
}

resource "opentelekomcloud_compute_instance_v2" "terraform" {
  name            = "terraform"
  image_name      = "${var.image}"
  flavor_name     = "${var.flavor}"
  key_pair        = "${opentelekomcloud_compute_keypair_v2.terraform.name}"
  security_groups = ["${opentelekomcloud_compute_secgroup_v2.terraform.name}"]
  floating_ip     = "${opentelekomcloud_compute_floatingip_v2.terraform.address}"

  network {
    uuid = "${opentelekomcloud_networking_network_v2.terraform.id}"
  }

  provisioner "remote-exec" {
    connection {
      user     = "${var.ssh_user_name}"
      private_key = "${file(var.ssh_key_file)}"
    }

    inline = [
      "sudo apt-get -y update",
      "sudo apt-get -y install nginx",
      "sudo service nginx start",
    ]
  }
}
