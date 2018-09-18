variable "region" {
  default = "eu-de"
}

variable "db_flavor" {
  default = "rds.mysql.s1.medium"
}

variable "db_name" {
  default = "<YOUR_DBNAME>"
}

variable "db_type" {
  default = "MySQL"
}

variable "db_version" {
  default = "5.6.33"
}

variable "vpc_id" {
  default = "<YOUR_VPC_ID>"
}

variable "existing_private_net_id" {
  default = "<YOUR_NETWORK_ID>"
}

variable "db_passwd" {
  default = "<YOUR_DB_PASSWORD>"
}

variable "db_port" {
  default = "<YOUR_DB_PORT>"
}

variable "availability_zone" {
  default = "<YOUR_AZ>"
}

variable "secgroup_name" {
  default = "<YOUR_SG_NAME>"
}
