---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_ip_address_group_v3"
sidebar_current: "docs-opentelekomcloud-resource-vpc-ip-address-group-v3"
description: |-
  Manages a VPC IP Address Group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC IP address group v3 you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/vpc_apis_v3/ip_address_group/index.html)

# opentelekomcloud_vpc_ip_address_group_v3

Manages a V3 VPC IP address group resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_vpc_ip_address_group_v3" "group_1" {
  name        = "group_1"
  description = "My VPC IP address group"
  ip_version  = 4
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) A unique name for the IP address group.
  The value can contain 1 to 64 characters, including letters, digits, underscores (`_`), hyphens (`-`), and periods (`.`).

* `description` - (Optional, String) Description about the IP address group.
  The value can contain up to 255 characters and cannot contain angle brackets (`<` or `>`).

* `ip_version` - (Required, Int, ForceNew) IP address version of the IP address group.
  Valid values are `4` (IPv4) and `6` (IPv6).

* `enterprise_project_id` - (Optional, String, ForceNew) ID of the enterprise project to which the IP address group belongs.
  Default: `0`.

* `max_capacity` - (Optional, Int, ForceNew) Maximum number of IP address entries in the IP address group.
  Range: 0 to 20. Default: `0`.

* `ip_set` - (Optional, List) IP address entries in the IP address group.
  Both IPv4 and IPv6 addresses are supported. An IP address entry can be:
  - A single IP address, for example, `192.168.21.25`
  - An IP address range, for example, `192.168.21.25-192.168.21.30`
  - A CIDR block, for example, `192.168.21.0/24`

* `ip_extra_set` - (Optional, List) IP address entries with remarks.
  The `ip_set` and `ip_extra_set` parameters cannot be both specified in a request.
  Structure is documented below.

The `ip_extra_set` block supports:

* `ip` - (Required, String) IP address entry. Can be a single IP address, an IP address range, or a CIDR block.

* `remarks` - (Optional, String) Remarks of the IP address entry.
  The value can contain 0 to 255 characters and cannot contain angle brackets (`<` or `>`).

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - IP address group ID.

* `project_id` - Indicates the project ID that the IP address group belongs to.

* `created_at` - Indicates the time when the IP address group was created.
  It is a UTC time in the format of `yyyy-MM-ddTHH:mm:ss`.

* `updated_at` - Indicates the time when the IP address group was last updated.
  It is a UTC time in the format of `yyyy-MM-ddTHH:mm:ss`.

* `status` - Status of the IP address group.
  Valid values are `NORMAL` (normal), `UPDATING` (being updated), and `UPDATE_FAILED` (update failed).

* `status_message` - Details about the IP address group status.

## Import

VPC IP Address Group V3 can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_ip_address_group_v3.group_1 <id>
```
