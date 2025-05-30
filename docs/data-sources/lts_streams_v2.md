---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_streams_v2"
sidebar_current: "docs-opentelekomcloud-datasource-lts-streams-v2"
description: |-
  Use this data source to get the list of LTS log streams within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log streams service you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/log_stream_management/querying_log_streams.html#listlogstreams)

# opentelekomcloud_lts_streams_v2

Use this data source to get the list of LTS log streams.

## Example Usage

```hcl
data "opentelekomcloud_lts_streams_v2" "test" {}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the log stream.

* `log_group_name` - (Optional, String) Specifies the name of the log group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - Shows the region in the log stream resources created.

* `streams` - All log streams that match the filter parameters.

  The [streams](#streams_struct) structure is documented below.

<a name="streams_struct"></a>
The `streams` block supports:

* `id` - The ID of the log stream.

* `name` - The name of the log stream.

* `ttl_in_days` - The log expiration time (days).

* `tags` - The key/value pairs to associate with the log stream.

* `created_at` - The creation time of the log stream, in RFC3339 format.
