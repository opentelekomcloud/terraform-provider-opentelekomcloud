---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_stream_v2"
sidebar_current: "docs-opentelekomcloud-resource-lts-stream-v2"
description: |-
  Manages a LTS Log Stream resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log group you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/log_stream_management/index.html)

# opentelekomcloud_lts_stream_v2

Manage a log stream resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "group_id" {}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = var.group_id
  stream_name = "test_stream"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) Specifies the ID of a created log group. Changing this parameter will create
  a new resource.

* `stream_name` - (Required, String, ForceNew) Specifies the log stream name. Changing this parameter will create a new
  resource.

* `ttl_in_days` - (Optional, Int) Specifies the log expiration time (days).
  The valid value is a non-zero integer from `-1` to `365`, defaults to `-1` which means inherit the log group settings.

* `tags` - (Optional, Map) Specifies the key/value pairs of the log stream.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log stream ID.

* `filter_count` - Number of log stream filters.

* `created_at` - The creation time of the log stream.

* `enterprise_project_id` - Shows the enterprise project ID to which the log stream belongs.

* `region` - Shows the region in the log group resource created.

## Import

The log stream can be imported using the group ID and stream ID separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_lts_stream_v2.test <group_id>/<id>
```
