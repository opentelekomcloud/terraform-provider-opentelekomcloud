terraform {
  required_version = ">= 0.13"

  required_providers {
    opentelekomcloud = {
      source  = "opentelekomcloud/opentelekomcloud"
      version = ">= 1.35.9"
    }
  }
}
