terraform {
  required_version = ">= 1.6.3"

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
