---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_address_group_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-address-group-v1"
description: |-
  Manages a CFW Address group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW address group you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/address_group_management/index.html)

# opentelekomcloud_cfw_address_group_v1

Manages a CFW Address Group resource within OpenTelekomCloud.

## Example Usage:
```hcl
variable object_id {}

resource "opentelekomcloud_cfw_address_group_v1" "group_1" {
  object_id    = var.object_id
  name         = "test-acc-tf-address-group"
  address_type = 0
}
```

## Argument Reference

The following arguments are supported:

* `object_id` - (Required, String, ForceNew) Protected object ID, which is used to distinguish between Internet border protection and VPC border protection after a cloud firewall is created. If the value of type is 0, the protected object ID belongs to the Internet border. If the value of type is 1, the protected object ID belongs to the VPC border.

* `name` - (Required, String) Specifies the CFW Address group name. The CFW address group name of the same type is unique in the same firewall.

* `description` - (Optional, String) Specifies the description of the address group.

* `address_type` - (Optional, Integer, ForceNew) Specifies the Internet protocol type of an address: `0` (IPv4), `1` (IPv6).

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the address group ID.
* `address_set_type` - Indicates the Address group type: `0` (user-defined address group), `1` (WAF back-to-source IP address group), `2` (DDoS back-to-source IP address group), or `3` (NAT64 address group).

## Import

CFW Address Group V1 resource can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_cfw_address_group_v1.group_1 b4cd6aeb0b7445d3bf271457c6941544in09
```

## Notes

But due to some attributes missing from the API response, it's required to ignore changes as below:

```hcl
resource "opentelekomcloud_cfw_address_group_v1" "group_1" {
  # ...

  lifecycle {
    ignore_changes = [
      object_id,
    ]
  }
}
```
