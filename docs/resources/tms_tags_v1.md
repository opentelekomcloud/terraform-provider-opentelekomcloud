---
subcategory: "Tag Management Service (TMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_tms_tags_v1"
sidebar_current: "docs-opentelekomcloud-resource-tms-tags-v1"
description: |-
  Manages an TMS Tags resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for TMS Tags you can get at
[documentation portal](https://docs.otc.t-systems.com/tag-management-service/api-ref/)

# opentelekomcloud_tms_tags_v1

Manages TMS tags resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_tms_tags_v1" "test" {
  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `tags` - (Required, List, ForceNew) Specifies an array of one or more predefined tags.

The `tags` block supports:

* `key` - (Required, String, ForceNew) Specifies the tag key. The value can contain up to 36 characters.
  Only letters, digits, hyphens (-), underscores (_), and Unicode characters from \u4e00 to \u9fff are allowed.

* `value` - (Required, String, ForceNew) Specifies the tag value. The value can contain up to 43 characters.
  Only letters, digits, periods (.), hyphens (-), and underscores (_), and Unicode characters from \u4e00 to \u9fff
  are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minute.
* `delete` - Default is 3 minute.
