---
subcategory: "Tag Management Service (TMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_tms_resource_tag_keys_v1"
sidebar_current: "docs-opentelekomcloud-datasource-tms-resource-tag-keys-v1"
description: |-
  Get TMS resource tag keys within OpenTelekomCloud.
---

# opentelekomcloud_tms_resource_tag_keys_v1

Use this data source to get the list of tag keys.

## Example Usage

```hcl
data "opentelekomcloud_tms_resource_tag_keys_v1" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional, String) Specifies the region ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `keys` - Indicates the tag keys.
