variable "region" {
  default = "eu-de"
}

variable "db_flavor" {
  default = "rds.mysql.c2.medium"
}

variable "db_name" {
  default = "test-rds"
}

variable "db_type" {
  default = "MySQL"
}

variable "db_version" {
  default = "5.6"
}

variable "db_passwd" {
  default = "TestPasswd!#@112"
}

variable "db_port" {
  default = "8635"
}

variable "availability_zone" {
  default = "eu-de-01"
}

variable "secgroup_name" {
  default = "rds-secgroup"
}
