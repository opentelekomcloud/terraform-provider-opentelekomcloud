---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_ips_protection_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-ips-protection-v1"
description: |-
  Configure IPS protection associated with CFW firewall within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW IPS protection you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/ips_management/index.html)

# opentelekomcloud_cfw_ips_protection_v1

  Configure IPS protection associated with CFW firewall within OpenTelekomCloud.

## Example Usage:
```hcl
variable object_id {}

resource "opentelekomcloud_cfw_ips_protection_v1" "protect_1" {
  object_id      = var.object_id
  ips_type       = 2
  feature_status = 1
  mode           = 0
}
```

## Argument Reference

The following arguments are supported:

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID, which is used to distinguish between Internet border protection and VPC border protection after a cloud firewall is created. If the value of type is `0`, the protected object ID belongs to the Internet border. If the value of type is `1`, the protected object ID belongs to the VPC border.

* `ips_type` - (Optional, Integer, ForceNew) Specifies the IPS patch type. Its value can only be `2` (virtual patch). Default: `2`.

* `feature_status` - (Required, Integer, ForceNew) Specifies the desired IPS virtual patching status: `0` (disabled), `1` (enabled).

* `mode` - (Required, Integer, ForceNew) Specifies the IPS protection mode: `0` (observation mode), `1` (strict mode), `2` (medium mode), or `3` (loose mode).

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `basic_defense_status` - Indicates the basic defense status: `0` (disabled), `1` (enabled).
* `ips_switch_id` - Indicates the IPS switch ID.
* `ips_protection_mode_id` - Indicates the IPS protection mode ID.



## Timeouts

This resource provides the following timeout configuration options:

* `create` - Default is 30 minutes.
