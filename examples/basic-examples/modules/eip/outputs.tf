
output "eip_address"{
  value = "${opentelekomcloud_vpc_eip_v1.eip_1.publicip.0.ip_address}"
}

output "eip_id"{
  value = "${opentelekomcloud_vpc_eip_v1.eip_2.id}"
}

output "lb_eip"{
  value = "${opentelekomcloud_vpc_eip_v1.eip_lb.publicip.0.ip_address}"
}
