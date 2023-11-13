variable "db_flavor" {
  default = "rds.mssql.s1.2xlarge"
}

variable "db_name" {
  default = "<YOUR_DBNAME>"
}

variable "db_type" {
  default = "SQLServer"
}

variable "db_version" {
  default = "2019_SE"
}

variable "vpc_id" {
  default = "<YOUR_VPC_ID>"
}

variable "network_id" {
  default = "<YOUR_NETWORK_ID>"
}

variable "db_passwd" {
  default = "<YOUR_DB_PASSWORD>"
}

variable "db_port" {
  default = "3365"
}

variable "availability_zone" {
  default = "<YOUR_AZ>"
}

variable "secgroup_name" {
  default = "<YOUR_SG_NAME>"
}
