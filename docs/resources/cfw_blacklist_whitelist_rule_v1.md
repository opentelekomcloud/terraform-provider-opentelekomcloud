---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_blacklist_whitelist_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-blacklist-whitelist-rule-v1"
description: |-
  Manages a CFW blacklist/whitelist rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW blacklist/whitelist rule you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/blacklist_whitelist_management/index.html)

# opentelekomcloud_cfw_blacklist_whitelist_rule_v1

Manages a CFW blacklist/whitelist rule resource within OpenTelekomCloud.

## Example Usage:
```hcl
variable object_id {}

resource "opentelekomcloud_cfw_blacklist_whitelist_rule_v1" "rule_1" {
  object_id    = var.object_id
  list_type    = 5
  direction    = 0
  address_type = 0
  address      = "1.1.1.1"
  protocol     = 6
  port         = "1"
  description  = "Test111161"
}
```

## Argument Reference

The following arguments are supported:

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID, which is used to distinguish between Internet border protection and VPC border protection after a cloud firewall is created. If the value of type is 0, the protected object ID belongs to the Internet border. If the value of type is 1, the protected object ID belongs to the VPC border.

* `list_type` - (Required, Integer, ForceNew) Specifies the list type. `4` (blacklist), `5` (whitelist).

* `direction` - (Required, Integer) Specifies the address direction: `0` (source), `1` (destination).

* `address_type` - (Required, Integer) Specifies the Internet protocol type of an address: `0` (IPv4), `1` (IPv6).

* `address` - (Required, String) Specifies the IP address.

* `protocol` - (Required, Integer) Specifies the Protocol type: `6` (TCP), `17` (UDP), `1` (ICMP), `58` (ICMPv6), or -`1` (any). 

* `port` - (Required, String) Specifies the destination port.

* `description` - (Optional, String) Specifies the description of the blacklist or whitelist rule.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the blacklist or whitelist rule ID.

## Import

CFW Blacklist or Whitelist Rule V1 resource can be imported using the object ID, `object_id`, the type of list, `list_type` and IP address, `address`, e.g.

```shell
terraform import opentelekomcloud_cfw_blacklist_whitelist_rule_v1.rule_1 <object_id>/<list_type>/<address>
```
