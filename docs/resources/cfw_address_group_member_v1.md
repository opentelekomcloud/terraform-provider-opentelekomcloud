---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_address_group_member_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-address-group-member-v1"
description: |-
  Manages a CFW Address group member resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW address group member you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/address_group_management/index.html)

# opentelekomcloud_cfw_address_group_member_v1

Manages a CFW Address Group Member resource within OpenTelekomCloud.

## Example Usage:
```hcl
variable set_id {}

resource "opentelekomcloud_cfw_address_group_member_v1" "group_1" {
  set_id  = var.set_id
  address = "1.1.1.1"
}
```

## Argument Reference

The following arguments are supported:

* `set_id` - (Required, String, ForceNew) Specifies the address group ID.

* `address_type` - (Optional, Integer, ForceNew) Specifies the Internet protocol type of an address: `0` (IPv4), `1` (IPv6).

* `address` - (Required, String, ForceNew) Specifies the IP Address.

* `description` - (Optional, String, ForceNew) Specifies the description of the address group member.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the address group member ID.
* `name` - Indicates the CFW Address group member name.

## Import

CFW Address Group V1 resource can be imported using the address group ID, `set_id` and IP address, `address`, e.g.

```shell
terraform import opentelekomcloud_cfw_address_group_member_v1.member_1 b4cd6aeb0b7445d3bf271457c6941544in09/address
```
