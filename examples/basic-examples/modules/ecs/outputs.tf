
output "ecs_id" {
  value = "${opentelekomcloud_compute_instance_v2.ecs_1.id}"
}
output "ecs2_id" {
  value = "${opentelekomcloud_compute_instance_v2.ecs_2.id}"
}
output "ecs_port" {
  value = "${opentelekomcloud_compute_instance_v2.ecs_1.network.0.port}"
}

output "ecs_ip" {
  value = "${opentelekomcloud_compute_instance_v2.ecs_1.access_ip_v4}"
}
output "ecs2_ip" {
  value = "${opentelekomcloud_compute_instance_v2.ecs_2.access_ip_v4}"
}

#output "hostname"{
#  value = "${opentelekomcloud_compute_instance_v2.ecs_1.hostname}"
#}

output  "access_ip_v4"{
 value = "${opentelekomcloud_compute_instance_v2.ecs_1.access_ip_v4}"
}
