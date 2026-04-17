---
subcategory: "Tag Management Service (TMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_tms_resource_tag_values_v1"
sidebar_current: "docs-opentelekomcloud-datasource-tms-resource-tag-values-v1"
description: |-
  Get TMS resource tag values by tag key within OpenTelekomCloud.
---

# opentelekomcloud_tms_resource_tag_values_v1

Use this data source to get the list of tag values by tag key.

## Example Usage

```hcl
data "opentelekomcloud_tms_resource_tag_values_v1" "test" {
  key = "tag_key"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, String) Specifies the tag key.

* `region_id` - (Optional, String) Specifies the region ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `values` - Indicates the tag values.
