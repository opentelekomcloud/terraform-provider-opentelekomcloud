terraform {
  required_version = ">= 0.13"

  required_providers {
    opentelekomcloud = {
      source  = "opentelekomcloud/opentelekomcloud"
      version = ">= 1.35.9"
    }
  }
  backend "s3" {
    endpoint = "https://obs.eu-de.otc.t-systems.com"
    skip_region_validation      = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true
  }
}

provider "opentelekomcloud" {
  cloud = "functest_cloud"
}

variable "module_files" {
  default = ""
}

variable "key" {
  default = ""
}

variable "index_files" {
  default = ""
}

variable "bucket" {
  default = ""
}

locals {
  main = compact(split(",", var.index_files))
  files = compact(split(",", var.module_files))
}

resource "opentelekomcloud_obs_bucket_object" "object" {
  count = length(local.files)

  bucket = var.bucket
  key    = "${var.key}/${split("terraform-visual-report/", local.files[count.index])[1]}"
  source = local.files[count.index]
  etag   = filemd5(local.files[count.index])
}

resource "opentelekomcloud_obs_bucket_object" "index" {
  count = length(local.main)

  bucket = var.bucket
  key    = split("main_page/", local.main[count.index])[1]
  source = local.main[count.index]
  etag   = filemd5(local.main[count.index])
}
