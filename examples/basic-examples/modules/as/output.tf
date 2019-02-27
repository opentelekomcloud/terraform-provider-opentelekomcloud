output "as_config_id" {
  value= "${opentelekomcloud_as_configuration_v1.as_config.id}"
}
output "as_config_name" {
  value= "${opentelekomcloud_as_configuration_v1.as_config.scaling_configuration_name}"
}

output "as_group_id" {
 value= "${opentelekomcloud_as_group_v1.my_as_group.id}"
}
output "as_group_name" {
 value= "${opentelekomcloud_as_group_v1.my_as_group.scaling_group_name}"
} 
