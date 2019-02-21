#define KMS variable
variable "kms_name" {}

variable "kms_pending_days" {
  default=7
}

variable "kms_region" {
  default="eu-de"
}
