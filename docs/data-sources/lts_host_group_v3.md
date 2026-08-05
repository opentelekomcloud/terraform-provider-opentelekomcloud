---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_host_group_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lts-host-group-v3"
description: |-
  Use this data source to query LTS host groups within T-Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for LTS host group you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/host_group_management/index.html)

# opentelekomcloud_lts_host_group_v3

Use this data source to query LTS host groups within T-Cloud Public (formerly OpenTelekomCloud).

## Example Usage

```hcl
variable "group_id" {}

data "opentelekomcloud_lts_host_group_v3" "test" {
  id = var.group_id
}
```

```hcl
variable "group_name" {}

data "opentelekomcloud_lts_host_group_v3" "test" {
  name = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) Specifies the ID of the host group.

* `name` - (Optional, String) Specifies the name of the host group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `host_groups` - List of host groups that match the filter parameters.

  The [host_groups](#host_groups_struct) structure is documented below.

<a name="host_groups_struct"></a>
The `host_groups` block supports:

* `id` - The host group ID.

* `name` - The host group name.

* `type` - The type of the host. Valid values are `linux` and `windows`.

* `host_ids` - List of host IDs in the group.

* `agent_access_type` - The type of host group access. Valid values are `IP` and `LABEL`.

* `labels` - Custom label list of the host group.

* `tags` - Key/value pairs attached to the host group.

* `created_at` - The creation time of the host group, in RFC3339 format.

* `updated_at` - The latest update time of the host group, in RFC3339 format.
