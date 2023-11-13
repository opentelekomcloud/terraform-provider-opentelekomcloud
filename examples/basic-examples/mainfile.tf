module "vpc_subnets" {
  source             = "./modules/vpc_subnet"
  vpc_name           = var.vpc_name
  vpc_cidr           = var.vpc_cidr
  subnet_name1       = var.subnet_name1
  subnet_cidr1       = var.subnet_cidr1
  subnet_gateway_ip1 = var.subnet_gateway_ip1
  primary_dns        = var.primary_dns
  secondary_dns      = var.secondary_dns
  subnet_name2       = var.subnet_name2
  subnet_cidr2       = var.subnet_cidr2
  subnet_gateway_ip2 = var.subnet_gateway_ip2
}

module "sg" {
  source        = "./modules/sg"
  secgroup_name = var.secgroup_name
}

module "ims" {
  source = "./modules/ims"
}

module "ecskey" {
  source     = "./modules/ecskey"
  public_key = file("public-2048.txt")
  key_name   = var.key_name
}

module "ecs" {
  source            = "./modules/ecs"
  key_name          = var.key_name
  security_groups   = module.sg.sg_name
  subnet_id         = module.vpc_subnets.subnet1_id
  subnet_id2        = module.vpc_subnets.subnet2_id
  ecs_name          = var.ecs_name
  image_id          = var.image_id
  flavor_id         = var.flavor_id
  region            = var.region
  availability_zone = var.availability_zone
  admin_pass        = var.admin_pass
  auto_recovery     = var.auto_recovery
}

module "eip" {
  source  = "./modules/eip"
  ecs_id  = module.ecs.ecs_id
  bw_name = var.bw_name
  bw_size = var.bw_size
}

module "lb" {
  source   = "./modules/lb"
  ecs_ip   = module.ecs.ecs_ip
  ecs2_ip  = module.ecs.ecs2_ip
  elb_name = var.elb_name
  lb_eip   = module.eip.lb_eip
  subnetid = module.vpc_subnets.subnetid1
  sg_id    = module.sg.sg_id
  elb_desc = var.elb_desc
}

module "rds" {
  source            = "./modules/rds"
  region            = var.region
  datastore_name    = var.datastore_name
  datastore_version = var.datastore_version
  flavor_name       = var.flavor_name
  rds_name          = var.rds_name
  rds_volume_type   = var.rds_volume_type
  rds_volume_size   = var.rds_volume_size
  availabilityzone  = var.availabilityzone
  vpc_id            = module.vpc_subnets.vpc_id
  subnetid          = module.vpc_subnets.subnet1_id
  securitygroupid   = module.sg.sg_id
  passwd            = var.passwd
}

module "sfs" {
  source            = "./modules/sfs"
  size              = var.sfs_size
  share_name        = var.sfs_share_name
  vpc_id            = module.vpc_subnets.subnet1_id
  share_description = "share_description"
  availability_zone = var.availability_zone
}

module "kms" {
  source   = "./modules/kms"
  kms_name = "kms_test"
}

module "evs" {
  source            = "./modules/evs"
  region            = var.region
  availability_zone = var.availability_zone
  volume_name       = var.volume_name
  desc              = var.desc
  volume_size       = var.volume_size
  volume_type       = var.volume_type
  ecs_id            = module.ecs.ecs_id
  ecs2_id           = module.ecs.ecs2_id
  image_id          = var.image_id
}

module "keypair" {
  source = "./modules/keypair"
}

module "dns" {
  source         = "./modules/dns"
  region         = var.region
  zone_name      = var.zone_name
  email          = var.email
  zone_desc      = var.zone_desc
  zone_type      = var.zone_type
  vpc_id         = module.vpc_subnets.vpc_id
  recordset_name = var.recordset_name
  recordset_desc = var.recordset_desc
  recordset_type = var.recordset_type
  records        = module.ecs.access_ip_v4
}

module "nat" {
  source     = "./modules/nat"
  nat_name   = var.nat_name
  nat_desc   = var.nat_desc
  vpc_id     = module.vpc_subnets.vpc_id
  subnet1_id = module.vpc_subnets.subnet1_id
  eip_id     = module.eip.eip_id
}

module "obs" {
  source      = "./modules/obs"
  bucket_name = var.bucket_name
}

module "iam" {
  source       = "./modules/iam"
  project_name = var.project_name
  project_desc = var.project_desc
  parent_id    = var.parent_id
  user_name    = var.iam_user_name
  user_desc    = var.user_desc
  user_passd   = var.user_passd
  group_name   = var.user_group_name
  group_desc   = var.group_desc
  user_status  = var.user_status
  role_name    = var.role_name
  domain_id    = var.domain_id
  region       = var.region

}


module "smn_topic" {
  source                = "./modules/smn"
  topic_name            = var.topic_name
  display_name          = var.display_name
  subscription_endpoint = var.subscription_endpoint
  subscription_protocol = var.subscription_protocol
  subscription_remark   = var.subscription_remark
}


module "anti-dos" {
  source = "./modules/anti-dos"
}

module "cce" {
  source        = "./modules/cce"
  cce_name      = var.cce_name
  cce_type      = var.cce_type
  cce_node_name = var.cce_node_name
  vpc_id        = module.vpc_subnets.vpc_id
  key_name      = var.key_name
  subnet_id     = module.vpc_subnets.subnet1_id
}


module "mrs" {
  source       = "./modules/mrs"
  mrs_name     = var.mrs_name
  job_name     = var.job_name
  vpc_id       = module.vpc_subnets.vpc_id
  keypair_name = var.key_name
  subnet_id    = module.vpc_subnets.subnet1_id
}

module "cts" {
  source      = "./modules/cts"
  bucket_name = var.cts_bucket_name
}

module "dcs" {
  source    = "./modules/dcs"
  dcs_name  = var.dcs_name
  capacity  = var.capacity
  dcs_desc  = var.dcs_desc
  vpc_id    = module.vpc_subnets.vpc_id
  subnet_id = module.vpc_subnets.subnet1_id
}

module "dms" {
  source            = "./modules/dms"
  dms_name          = var.dms_name
  dms_desc          = var.dms_desc
  queue_mode        = var.queue_mode
  redrive_policy    = var.redrive_policy
  max_consume_count = var.max_consume_count
  group_name        = var.group_name
}


module "firewall" {
  source              = "./modules/fw"
  rule_protocol       = var.rule_protocol
  rule_action         = var.rule_action
  rule_name           = var.rule_name
  rule_desc           = var.rule_desc
  policy_name         = var.policy_name
  policy_desc         = var.policy_desc
  firewall_group_name = var.firewall_group_name
  firewall_group_desc = var.firewall_group_desc
}

module "deh" {
  source            = "./modules/deh"
  deh_name          = var.deh_name
  host_type         = var.host_type
  availability_zone = var.availability_zone
}
module "csbs" {
  source      = "./modules/csbs"
  backup_name = var.backup_name
  backup_desc = var.backup_desc
  resource_id = module.ecs.ecs_id
}
module "vbs" {
  source         = "./modules/vbs"
  backup_name    = var.backup_name
  backup_desc    = var.backup_desc
  volume_id      = module.evs.volume_id
  to_project_ids = var.to_project_ids
}
module "ces" {
  source      = "./modules/ces"
  alarm_name  = var.alarm_name
  alarm_desc  = var.alarm_desc
  as_group_id = module.as.as_group_id
  ecs_id      = module.ecs.ecs_id
}

module "as" {
  source                = "./modules/as"
  flavor_id_as          = var.flavor_id_as
  image_id_as           = var.image_id_as
  volume_size_as        = var.volume_size_as
  volume_type_as        = var.volume_type_as
  key_name_as           = var.key_name_as
  vpc_name_as           = var.vpc_name_as
  vpc_cidr_as           = var.vpc_cidr_as
  subnet_name1_as       = var.subnet_name1_as
  subnet_cidr1_as       = var.subnet_cidr1_as
  subnet_gateway_ip1_as = var.subnet_gateway_ip1_as
  primary_dns_as        = var.primary_dns_as
  secondary_dns_as      = var.secondary_dns_as
  secgroup_name_as      = var.secgroup_name_as
  region_as             = var.region_as
  availability_zone_as  = var.availability_zone_as
  public_key_as         = file("public-2048.txt")
  alarm_id              = module.ces.alarm_id
  listenerid            = module.lb.listenerid
}

module "rts" {
  source            = "./modules/rts"
  rts_name          = var.rts_name
  instance_type     = var.flavor_id
  image_id          = var.image_id
  availability_zone = var.availabilityzone
  subnet_id         = module.vpc_subnets.subnet1_id
  config_name       = var.config_name
  ecs_id            = module.ecs.ecs_id
  ecs2_id           = module.ecs.ecs2_id
}


