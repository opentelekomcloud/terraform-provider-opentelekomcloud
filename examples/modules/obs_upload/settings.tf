terraform {
  required_version = ">= 0.13"

  required_providers {
    opentelekomcloud = {
      source  = "opentelekomcloud/opentelekomcloud"
      version = ">= 1.35.9"
    }
  }
}

provider "opentelekomcloud" {
  cloud = "functest_cloud"
}

variable "folder_path" {
  default = ""
}

variable "key" {
  default = ""
}

variable "index_path" {
  default = ""
}

variable "bucket" {
  default = ""
}

resource "opentelekomcloud_obs_bucket_object" "object" {
  for_each = fileset(var.folder_path, "**")

  bucket = var.bucket
  key    = "${var.key}/${each.value}"
  source = "${var.folder_path}/${each.value}"
  etag   = filemd5("${var.folder_path}/${each.value}")
}


resource "opentelekomcloud_obs_bucket_object" "index" {
  for_each = fileset(var.index_path, "**")

  bucket = var.bucket
  key    = each.value
  source = "${var.index_path}/${each.value}"
  etag   = filemd5("${var.index_path}/${each.value}")
}
