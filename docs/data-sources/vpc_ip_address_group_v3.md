---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_ip_address_group_v3"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-ip-address-group-v3"
description: |-
  Get details of a VPC IP Address Group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC IP address group v3 you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/vpc_apis_v3/ip_address_group/index.html)

# opentelekomcloud_vpc_ip_address_group_v3

Get details of a V3 VPC IP address group resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "id" {}

data "opentelekomcloud_vpc_ip_address_group_v3" "group_1" {
  id = var.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) IP address group ID.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `name` - A unique name for the IP address group.

* `description` - Description about the IP address group.

* `ip_version` - IP address version of the IP address group.

* `enterprise_project_id` - ID of the enterprise project to which the IP address group belongs.

* `max_capacity` - Maximum number of IP address entries in the IP address group.

* `ip_set` - IP address entries in the IP address group.

* `ip_extra_set` - IP address entries with remarks. Structure is documented below.
    * `ip` - IP address entry. Can be a single IP address, an IP address range, or a CIDR block.
    * `remarks` - Remarks of the IP address entry.

* `project_id` - Indicates the project ID that the IP address group belongs to.

* `created_at` - Indicates the time when the IP address group was created.
  It is a UTC time in the format of `yyyy-MM-ddTHH:mm:ss`.

* `updated_at` - Indicates the time when the IP address group was last updated.
  It is a UTC time in the format of `yyyy-MM-ddTHH:mm:ss`.

* `status` - Status of the IP address group.
  Valid values are `NORMAL` (normal), `UPDATING` (being updated), and `UPDATE_FAILED` (update failed).

* `status_message` - Details about the IP address group status.
