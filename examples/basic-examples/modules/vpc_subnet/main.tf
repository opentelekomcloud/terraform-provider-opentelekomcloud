# create random number
resource "random_id" "vpc" {
 byte_length = 4
}

#create VPC
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "${var.vpc_name}-${random_id.vpc.id}"
  cidr = "${var.vpc_cidr}"
}



resource "opentelekomcloud_vpc_v1" "vpc_peer" {
  name = "vpc_peer"
  cidr = "10.0.0.0/16"
}
#create subnet1
resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "${var.subnet_name1}"
  cidr = "${var.subnet_cidr1}"
  gateway_ip = "${var.subnet_gateway_ip1}"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
 # primary_dns = "${var.primary_dns}"
 # secondary_dns = "${var.secondary_dns}"
}

#create subnet2
resource "opentelekomcloud_vpc_subnet_v1" "subnet_2" {
  name = "${var.subnet_name2}"
  cidr = "${var.subnet_cidr2}"
  gateway_ip = "${var.subnet_gateway_ip2}"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
#  primary_dns = "${var.primary_dns}"
#  secondary_dns = "${var.secondary_dns}"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering" {
  name = "peer_test"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${opentelekomcloud_vpc_v1.vpc_peer.id}"
}

#resource "opentelekomcloud_vpc_route_v2" "vpc_route" {
#  type  = "peering"
#  nexthop  = "172.16.0.2"
#  destination = "172.16.0.0/12"
#  vpc_id = "8c2e48cc-7e52-46c1-8f7e-2ea2b893384b"
# }

#resource "opentelekomcloud_vpc_peering_connection_accepter_v2" "peer" {
#    vpc_peering_connection_id = "${opentelekomcloud_vpc_peering_connection_v2.peering.id}"
#    accept = true
#}
