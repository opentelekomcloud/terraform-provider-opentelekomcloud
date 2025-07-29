---
subcategory: "Tag Management Service (TMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_tms_tags_v1"
sidebar_current: "docs-opentelekomcloud-datasource-tms-tags-v1"
description: |-
  List TMS Tags resource within OpenTelekomCloud.
---

# opentelekomcloud_tms_tags_v1

Use this data source to get the list of predefined tags within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_tms_tags_v1" "test" {}
```

## Argument Reference

The following arguments are supported:

* `key` - (Optional, String) Specifies the tag key. Fuzzy search is supported. Key is case-insensitive.

* `value` - (Optional, String) Specifies the tag value. Fuzzy search is supported. Value is case-insensitive.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the value of the tag.
