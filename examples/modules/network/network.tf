data "opentelekomcloud_networking_secgroup_v2" "default_security_group" {
  name = "default"
}

resource "opentelekomcloud_vpc_v1" "shared_vpc" {
  name   = "vpc_default"
  cidr   = "192.168.0.0/16"
  shared = true
}

resource "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name        = "subnet_default"
  cidr        = "192.168.0.0/16"
  gateway_ip  = "192.168.0.1"
  vpc_id      = opentelekomcloud_vpc_v1.shared_vpc.id
  dhcp_enable = true
  dns_list    = ["100.125.4.25", "1.1.1.1"]
}

output "default_security_group_id" {
  value = data.opentelekomcloud_networking_secgroup_v2.default_security_group.id
}

output "shared_vpc" {
  value = opentelekomcloud_vpc_v1.shared_vpc
}

output "shared_subnet" {
  value = opentelekomcloud_vpc_subnet_v1.shared_subnet
}
