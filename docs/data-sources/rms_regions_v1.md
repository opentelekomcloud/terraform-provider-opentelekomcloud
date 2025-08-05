---
subcategory: "Config"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rms_regions_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rms-regions-v1"
description: |-
  Manages an RMS regions data source, used to query regions visible to you, within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RMS regions you can get at
[documentation portal](https://docs.otc.t-systems.com/config/api-ref/apis/region_management/index.html)


# opentelekomcloud_rms_regions_v1

Manages an RMS regions data source, used to query regions visible to you, within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_rms_regions_v1" "test" {}
```

## Attribute Reference

The following attributes are exported:

* `regions` - Specifies the list of region information. The structure is documented below:
    * `region_id` - Specifies the region ID.
    * `display_name` - Specifies the display name.
