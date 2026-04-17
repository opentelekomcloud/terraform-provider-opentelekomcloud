---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_service_group_member_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-service-group-member-v1"
description: |-
  Manages a CFW Service group member resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW service group member you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/service_group_management/index.html)

# opentelekomcloud_cfw_service_group_member_v1

Manages a CFW Service Group Member resource within OpenTelekomCloud.

## Example Usage:
```hcl
variable set_id {}

resource "opentelekomcloud_cfw_service_group_member_v1" "member_1" {
  set_id      = var.set_id
  protocol    = 6
  source_port = "1"
  dest_port   = "1"
  description = "Test611"
}
```

## Argument Reference

The following arguments are supported:

* `set_id` - (Required, String, ForceNew) Specifies the service group ID.

* `protocol` - (Required, Integer, ForceNew) Specifies the Protocol type: `6` (TCP), `17` (UDP), `1` (ICMP), `58` (ICMPv6), or -`1` (any). 

* `source_port` - (Required, String, ForceNew) Specifies the source port.

* `dest_port` - (Required, String, ForceNew) Specifies the destination port.

* `description` - (Optional, String, ForceNew) Specifies the description of the service group member.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the service group member ID.

## Import

CFW Service Group Member V1 resource can be imported using the service group ID, `set_id` and member ID, `id`, e.g.

```shell
terraform import opentelekomcloud_cfw_service_group_member_v1.member_1 <set_id>/<id>
```
