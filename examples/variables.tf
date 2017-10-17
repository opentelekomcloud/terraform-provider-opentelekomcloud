### OpenTelekomCloud Credentials
variable "username" {
  default = "tf_user"
}

variable "password" {
  default = "Huawei@123.UhOh"
}

variable "domain_name" {
  default = "OTC-EU-DE-00000000001000022296"
}

variable "tenant_name" {
  default = "eu-de"
}

variable "endpoint" {
  default = "https://iam.eu-de.otc.t-systems.com:443/v3"
}

### OTC Specific Settings
variable "external_network" {
  default = "admin_external_net"
}

### Project Settings
variable "project" {
  default = "terraform"
}

variable "subnet_cidr" {
  default = "192.168.10.0/24"
}

variable "ssh_pub_key" {
  default = "~/.ssh/id_rsa.pub"
}

### DNS Settings
variable "dnszone" {
  default = ""
}

variable "dnsname" {
  default = "webserver"
}

### VM (Instance) Settings
variable "instance_count" {
  default = "1"
}

variable "disk_size_gb" {
  default = "0"
}

variable "flavor_name" {
  default = "s1.medium"
}

variable "image_name" {
  default = "Standard_CentOS_7_latest"
}
