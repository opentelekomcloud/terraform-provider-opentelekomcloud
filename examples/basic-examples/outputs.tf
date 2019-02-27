
output "vpc_id" {
  value = "${module.vpc_subnets.vpc_id}"
}

output "subnet1_id" {
  value = "${module.vpc_subnets.subnet1_id}"
}

output "subnet2_id" {
  value = "${module.vpc_subnets.subnet2_id}"
}

output "ecskey_name" {
  value = "${module.ecskey.ecskey_name}"
}

output "ecs_id" {
  value = "${module.ecs.ecs_id}"
}

output "ecs2_id" {
  value = "${module.ecs.ecs2_id}"
}

output "ecs_port" {
  value = "${module.ecs.ecs_port}"
}


output "sg_id" {
  value = "${module.sg.sg_id}"
}

output "sg_name" {
  value = "${module.sg.sg_name}"
}

output "eip_address" {
  value = "${module.eip.eip_address}"
}

output "agency_id" {
  value = "${module.iam_agency.agency_id}"
}

output "delegated_domain_name" {
  value = "${module.iam_agency.delegated_domain_name}"
}
 
output "project_role" {
  value = "${module.iam_agency.project_role}"
}

output "domain_roles" {    
  value = "${module.iam_agency.domain_roles}"
}
 
output "project_role" {
  value = "${module.iam_agency.project_role}"
}

output "domain_roles" {    
  value = "${module.iam_agency.domain_roles}"
}

output "topic_name" {
  value = "${module.smn_topic.topic_name}"
}

output "topic_urn_1" {
  value = "${module.smn_topic.topic_urn}"
}
#
output "push_policy_1" {
  value = "${module.smn_topic.push_policy}"
}
#
output "topic_urn_2" {    
  value = "${module.smn_subscription.topic_urn}"
}
#
output "push_policy_2" {
  value = "${module.smn_subscription.push_policy}"
}
# 
output "subscription_urn" {
  value = "${module.smn_subscription.subscription_urn}"
}
#
output "subscription_endpoint" {
  value = "${module.smn_subscription.subscription_endpoint}"
}
#
output "subscription_protocol" {
  value = "${module.smn_subscription.subscription_protocol}"
}
#
##### vpc_peering ####
#
output "vpc_peer_01" {
  value = "${module.vpc_peering.vpc_peer_01}"
}
output "vpc_peer_02" {
  value = "${module.vpc_peering.vpc_peer_02}"
}
output "peering_id" {
  value = "${module.vpc_peering.peering_id}"
}
#
##### as #####
#
output "as_config_id" {
  value = "${module.as.as_config_id}"
}
output "as_config_name" {
  value = "${module.as.as_config_name}"
}
output "as_group_id" {
  value = "${module.as.as_group_id}"
}
output "as_group_name" {
  value = "${module.as.as_group_name}"
}

