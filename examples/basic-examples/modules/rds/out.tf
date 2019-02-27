output "flavor_id"{
  value = "${data.opentelekomcloud_rds_flavors_v1.flavor.id}"
}

output "hostname"{
  value = "${opentelekomcloud_rds_instance_v1.instance.hostname}"
}
