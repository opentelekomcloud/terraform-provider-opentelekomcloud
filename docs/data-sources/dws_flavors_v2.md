---
subcategory: "Data Warehouse Service (DWS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dws_flavors_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dws-flavors-v2"
description: |-
Get details about DWS flavors from OpenTelekomCloud
---

# opentelekomcloud_dws_flavors_v2

Use this data source to get details about flavors from OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_dws_flavors_v2" "flavor" {
  vcpus = 32
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the dws cluster client. If omitted, the
  provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the availability zone name.

* `vcpus` - (Optional, String) Specifies the vcpus of the dws node flavor.

* `memory` - (Optional, String) Specifies the ram of the dws node flavor in GB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID in UUID format.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `flavor_id` - The name of the dws node flavor. It is referenced by `node_type` in `opentelekomcloud_dws_flavors_v2`.
* `vcpus` - Indicates the vcpus of the dws node flavor.
* `memory` - Indicates the ram of the dws node flavor in GB.
* `volumetype` - Indicates Disk type.
    + **LOCAL_DISK**: common I/O disk
    + **SSD**: ultra-high I/O disk
* `size` - Indicates the Disk size in GB.
* `availability_zone` - Indicates the availability zone where the node resides.
