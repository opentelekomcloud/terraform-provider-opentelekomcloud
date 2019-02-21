
output "vpc_id" {
  value = "${opentelekomcloud_vpc_v1.vpc_1.id}"
}

output "subnet1_id" {
  value = "${opentelekomcloud_vpc_subnet_v1.subnet_1.id}"
}

output "subnetid1" {
  value = "${opentelekomcloud_vpc_subnet_v1.subnet_1.subnet_id}"
}

output "subnet2_id" {
  value = "${opentelekomcloud_vpc_subnet_v1.subnet_2.id}"
}

