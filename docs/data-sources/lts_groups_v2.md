---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_groups_v2"
sidebar_current: "docs-opentelekomcloud-datasource-lts-groups-v2"
description: |-
  Use this data source to get the list of LTS log groups within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS Groups you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/log_group_management/querying_all_log_groups_of_an_account.html#listloggroups)

# opentelekomcloud_lts_groups_v2

Use this data source to get the list of LTS log groups.

## Example Usage

```hcl
data "opentelekomcloud_lts_groups_v2" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - Region in which to query the resources are placed.

* `groups` - All log groups that match the filter parameters.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - The log group ID.

* `name` - The log group name.

* `ttl_in_days` - The log expiration time(days).

* `tags` - The key/value pairs to associate with the log group.

* `enterprise_project_id` - The enterprise project ID to which the log group belongs.

* `created_at` - The creation time of the log group, in RFC3339 format.
