package src

const WafdVars = `

#############
# Environment
#############

variable "environment" {
  default = "throttle-wafd-test"
}

###################
#   OTC auth config
###################

variable "region" {
  default = "eu-de"
}

variable "otc_domain" {
  default = "eu-de"
}

variable "auth_url" {
  default = "https://iam.eu-de.otc.t-systems.com:443/v3"
}

variable "tenant_name" {
  default = "eu-de"
}

variable "access_key" {
  default = ""
}

variable "secret_key" {
  default = ""
}

variable "key" {
  default = ""
}

##########
# VPC vars
##########

variable "vpc_cidr" {
  description = "CIDR of the VPC"
  default     = "10.1.0.0/24"
}

#############
# Subnet vars
#############

variable "subnet_cidr" {
  description = "CIDR of the Subnet"
  default     = "10.1.0.0/24"
}

variable "subnet_gateway_ip" {
  description = "Default gateway of the Subnet"
  default     = "10.1.0.1"
}

variable "subnet_primary_dns" {
  description = "Primary DNS server of the Subnet"
  default     = "100.125.4.25"
}

variable "subnet_secondary_dns" {
  description = "Secondary DNS server of the Subnet"
  default     = "100.125.129.199"
}

###########
# WAFD vars
###########

variable "wafd_az" {
  description = "Availability Zone 1 (Biere)"
  default     = "eu-de-01"
}

variable "wafd_flavor" {
  description = "Name of the wafd flavor"
  default     = "s2.large.2"
}

variable "wafd_arch" {
  description = "Name of the wafd arch"
  default     = "x86"
}
`
