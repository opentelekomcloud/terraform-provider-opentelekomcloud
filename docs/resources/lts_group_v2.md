---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_group_v2"
sidebar_current: "docs-opentelekomcloud-resource-lts-group-v2"
description: |-
  Manages a LTS Log Group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log group you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/log_group_management/index.html)

# opentelekomcloud_lts_group_v2

Manages a log group resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "log_group_1"
  ttl_in_days = 30
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String, ForceNew) Specifies the log group name. Changing this parameter will create a new resource.

* `ttl_in_days` - (Required, Int) Specifies the log expiration time(days).
  The value is range from `1` to `365`.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the log group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log group ID.

* `created_at` - The creation time of the log group.

* `enterprise_project_id` - Shows the enterprise project ID to which the log group belongs.

* `region` - Shows the region in the log group resource created.

## Import

The log group can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_lts_group_v2.test <id>
```
