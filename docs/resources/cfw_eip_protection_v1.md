---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_eip_protection_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-eip-protection-v1"
description: |-
  Enable or Disable EIP protection using CFW firewall within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW EIP protection you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/eip_management/index.html)

# opentelekomcloud_cfw_eip_protection_v1

  Enable or Disable EIP protection using CFW firewall within OpenTelekomCloud.

## Example Usage:
```hcl
variable firewall_id {}
variable object_id {}
variable eip_id {}

resource "opentelekomcloud_cfw_eip_protection_v1" "protect_1" {
  firewall_id = var.firewall_id
  object_id   = var.object_id
  status      = 0
  eip_id      = var.eip_id
}
```

## Argument Reference

The following arguments are supported:

* `firewall_id` - (Required, String, ForceNew) Specifies the Firewall ID.

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID, which is used to distinguish between Internet border protection and VPC border protection after a cloud firewall is created. If the value of type is 0, the protected object ID belongs to the Internet border. If the value of type is 1, the protected object ID belongs to the VPC border.

* `status` - (Required, Integer, ForceNew) Specifies the desired EIP protection status: `0` (protected), `1` (unprotected). 

* `eip_id` - (Required, String, ForceNew) Specifies the EIP ID.

* `public_ip` - (Optional, String, ForceNew) Specifies the EIP IPV4 address.

* `public_ipv6` - (Optional, String, ForceNew) Specifies the EIP IPV6 address.

## Timeouts

This resource provides the following timeout configuration options:

* `create` - Default is 30 minutes.
