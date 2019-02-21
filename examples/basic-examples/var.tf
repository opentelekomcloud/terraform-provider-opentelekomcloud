####### Provider variable #######

variable user_name {
   default = "your user_name"
}
variable tenant_name {
   default = "eu-de"
}
variable password {
   default = "your passwd"

}
variable auth_url {
   default = "https://iam.eu-de.otc.t-systems.com/v3"
}
variable domain_id {
   default = "domain_id"
}

####### Public or global variable #######

variable region {
   default = "eu-de"
}
variable availability_zone {
   default = "eu-de-01"
}

####### vpc variable #######

variable vpc_name {
   default = "terraform-test-vpc"
}
variable vpc_cidr {
   default = "192.168.0.0/16"
}
variable subnet_name1 {
   default = "terrtest-subnet-1"
}
variable subnet_cidr1 {
   default = "192.168.0.0/24"
}
variable subnet_gateway_ip1 {
   default = "192.168.0.1"
}
variable primary_dns {
   default = "100.125.1.250"
}
variable secondary_dns {
   default = "100.125.21.250"
}
variable subnet_name2 {
   default = "terrtest-subnet-2"
}
variable subnet_cidr2 {
   default = "192.168.10.0/24"
}
variable subnet_gateway_ip2 {
   default = "192.168.10.1"
}

####### securitygroup variable #######

variable secgroup_name {
   default = "secgroup_tftest"
}

####### keypair variable #######

variable key_name {
   default = "KeyPair-h"
}


####### ecs variable #######

variable ecs_name {
   default = "ecs-zht-tftest"
}
variable image_id {
   default = "53b2fbb5-ef2c-412a-bb0a-571436fa78ad"
}
variable flavor_id {
   default = "m1.large"
}
variable admin_pass{
   default = "your pass"
}
variable auto_recovery{
   default = "true"
}
####### eip variable #######

variable bw_name {
   default = "bandwidth_1546875"
}
variable bw_size {
   default = 5
}

variable bucket_name{
   default = "test_bucket"
}

####### elb variable #######

variable elb_name {
   default = "tfelb_45654654"
}

variable elb_desc {
   default = "create by terraform"
}

####### rds variable #######

variable datastore_name {
   default = "MySQL"
}
variable datastore_version {
   default = "5.6.30"
}
variable flavor_name {
   default = "rds.mysql.s1.medium"
}
variable rds_name {
   default = "tf_test_rds_mysql_0221"
}
variable rds_volume_type {
   default = "COMMON"
}
variable rds_volume_size {
   default = 50
}
variable availabilityzone {
   default = "eu-de-01"
}
variable passwd {
   default = "Huawei@12"
}

####### sfs variable #######

variable sfs_size {
   default = 1
}

variable sfs_share_name {
   default = "sfs_share"
}
##kms variable
variable kms_name {
   default = "kms_terraform"
}
#
######## evs variable #######

variable volume_name {
   default = "volume_lmm"
}

variable desc { 
   default = "create a volume"
  
}  
variable volume_size{
   default  = "1"
  
}  
variable volume_type{
   default = "SATA"
 
}

######## dns variable #######

variable zone_name {
   default = "example.com"
}
variable email {
   default = "example@example.com"
}
  
variable zone_desc 
{
   default = "create a private dns"
}
variable zone_type {
   default = "private"
} 
variable recordset_name{
   default = "test.example.com"
}
variable recordset_desc{
   default = "create a record set"
} 
  
variable recordset_type {
   default = "A"
} 

######## nat variable #######

variable "nat_name" { 
   default = "nat_test"
}

variable "nat_desc" {
   default = "nat_createByTerrafrom"
}

######## ces alarm-rules#######

variable "alarm_name" { 
   default = "alarm_test"
}

variable "alarm_desc" {
   default = "alarm_createByTerrafrom"
}

######## define smn_topic variable #######

variable "topic_name" {
   default = "terrtest-smn-topic-name"
}
variable "display_name" {
   default = "terrtest-smn-topic-display-name"
}


######## define smn_subscription variable #######

variable "subscription_endpoint" {
   default = "hanmeina2@huawei.com"   
}
variable "subscription_protocol" {
   default = "email"
}
variable "subscription_remark" {
   default = "O&M"
}

######## vpc_peering variable #######
#
variable "vpc_name_peering" {
   default = "vpc_peer_01,vpc_peer_02"
}
variable "vpc_cidr_peering" {
   default  = "172.16.0.0/16,10.0.0.0/24"
}
#
variable "peering_name" {
   default = "peering_huaweicloud"
}
#

###########tags variable################
variable "tags" {
   default = ["key1.value1","key2.value2"]
}


######## metadata variable #######
variable "metadata" {
   type = "map"
   default = {meta1="",meta2="data2"}
}

######## vbs variable #######
variable "backup_name" {
   default = "backup_terraform"
}

variable "backup_desc" {
   default = "created by terraform"
}

variable "to_project_ids" {
   default = "9d3f1f127e944a00811cddedb108dda1"
}

######## deh variable #######
variable  "deh_name" {
   default = "deh_test_terraform"
}
variable  "host_type" {
   default = "h1"
}

######## dms variable #######
variable dms_name {
   default = "dms_test"
}
variable dms_desc {
   default = "create by terraform"
}
variable queue_mode {
   default = "FIFO"
}
variable redrive_policy {
   default = "enable"
}
variable max_consume_count {
   default = "80"
}
variable group_name {
   default = "group_test"
}
######## firewall variable #######
variable "rule_protocol" {
   default = "tcp"
}
variable "rule_action" {
   default = "allow"
}
variable  "rule_name"{
   default = "rule_test"
}
variable  "rule_desc"{
   default = "create by terraform"
}
variable  "policy_name"{
   default = "policy_test"

}
variable  "policy_desc"{
   default = "create by terraform"
}
variable  "firewall_group_name"{
   default = "firewall_group_test"
}
variable  "firewall_group_desc"{
   default = "create by terraform"
}

######## dcs variable #######
variable dcs_name {
   default = "dcs_test"
}

variable capacity {
   default = 2
}
variable dcs_desc {
   default = "create by terraform"
}
######## cts variable #######
variable "cts_bucket_name" { 
   default = "obs-8f58"   
}

######## as variable #######
variable flavor_id_as {default = "m1.large"}
variable image_id_as {default = "53b2fbb5-ef2c-412a-bb0a-571436fa78ad"}
variable key_name_as {default = "KeyPair-as"}
variable vpc_name_as {default = "vpc-as"}
variable vpc_cidr_as {default = "192.168.0.0/16"}
variable subnet_name1_as { default = "terrtest-subnet-1"}
variable subnet_cidr1_as {default = "192.168.0.0/24"}
variable subnet_gateway_ip1_as { default = "192.168.0.1"}
variable primary_dns_as { default = "100.125.1.250"}
variable secondary_dns_as {  default = "100.125.21.250"}
variable availability_zone_as {default = "eu-de-01"}
variable region_as {default = "eu-de"}
variable volume_type_as { default ="SATA"}
variable volume_size_as {default = "40" }
variable secgroup_name_as {default = "secgroup-As"}

######## mrs variable #######
variable mrs_name {
   default = "mrs_test"
}
variable  job_name{
   default = "job_test"
}


######## cce variable #######
variable  cce_name {
   default = "cce-test"
}
variable  cce_type {
   default = "VirtualMachine"
}
variable  cce_node_name {
   default = "cce_node_test"
}
 
 
 
######## rts variable #######
variable  rts_name {
   default = "rts_test"
}
variable  config_name {
   default = "config_test"
}
######## iam variable #######
variable project_name {
   default = "eu-de_project_test"
}
variable project_desc {
   default = "terraform test project"
}
variable parent_id {
   default = "b730519ca7064da2a3233e86bee139e4"
}
variable iam_user_name {
   default = "user_test"
}
variable user_desc {
   default = "terraform test user"
}
variable user_passd {
   default = "Huawei@123"
}
variable user_status {
   default = "true"
}
variable user_group_name {
   default = "group_test"
}
variable group_desc {
   default = "group_desc"
}
variable role_name {
   default = "secu_admin"
}
