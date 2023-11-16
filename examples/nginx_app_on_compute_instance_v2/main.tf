module "network" {
  source = "../modules/network"
}

resource "opentelekomcloud_compute_keypair_v2" "terraform" {
  name       = "terraform"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLotBCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAnOfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZqd9LvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TaIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIF61p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = module.network.default_security_group_id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_2" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = module.network.default_security_group_id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_3" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = module.network.default_security_group_id
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "fip_associated_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_1.address
  port_id     = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  admin_state_up = true
  fixed_ip {
    subnet_id  = module.network.shared_subnet.network_id
    ip_address = "10.0.0.5"
  }
  network_id = module.network.shared_subnet.vpc_id
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "terraform"
  image_name      = var.image
  flavor_id       = "s3.large.1"
  key_pair        = opentelekomcloud_compute_keypair_v2.terraform.name
  security_groups = [module.network.default_security_group_id]
  depends_on      = [ opentelekomcloud_networking_floatingip_associate_v2.fip_associated_1 ]

  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }

  connection {
    user        = "linux"
    host        = opentelekomcloud_networking_floatingip_v2.fip_1.address
    private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIBUwIBADANBgkqhkiG9w0BAQEFAASCAT0wggE5AgEAAkEAu+qgVpV6mqbaGW1Q
n6eDPzhwentQPPiXwG1665M9+gjW4pUQ0RudBc0fkUU/O+Q0UMT8ZV/I2hSenCVy
JoyPEwIDAQABAkAbyksEAv8qt9oxQHVX5xIF23bm5i2rlqf6kTZIeHIF89/NNJ2E
sejiqFIWqPc5a00Scn+ymdCvjC25JVyup9cBAiEA4a+7WhPmgS54yNHjwkG2pflz
cfH1V7qPqlBKIGLwZbMCIQDVKCsZ6eoNdQoLVmK0zii8XDCgL8HWMrm/bytbYM9B
IQIgVdcAXKebEeF6IW/rwDQ8Y2644UsVdTPJdw8o0p6vLw8CIDqm191EiPt09fOS
rIxVoc3ajCK3oV2ADa5IN6ToKX8hAiBPuNCCIYcZz0tAzWX7I1OYMI3UhJjtrESg
mYFrsJ4gHw==
-----END RSA PRIVATE KEY-----
EOT
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt -y update",
      "sudo apt -y install nginx",
      "sudo service nginx start",
    ]
  }
}
