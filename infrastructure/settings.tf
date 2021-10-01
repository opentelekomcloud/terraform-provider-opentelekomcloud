terraform {
  required_version = ">= 0.13"

  required_providers {
    opentelekomcloud = {
      source  = "opentelekomcloud/opentelekomcloud"
      version = ">= 1.23.9"
    }
  }

  backend "s3" {
    key                         = "tf_state"
    endpoint                    = "obs.eu-de.otc.t-systems.com"
    bucket                      = "obs-tf"
    region                      = "eu-de"
    skip_region_validation      = true
    skip_credentials_validation = true
  }
}

provider "opentelekomcloud" {}
