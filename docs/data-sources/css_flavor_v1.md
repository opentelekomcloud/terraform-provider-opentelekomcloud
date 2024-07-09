---
subcategory: "Cloud Search Service (CSS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_css_flavor_v1"
sidebar_current: "docs-opentelekomcloud-datasource-css-flavor-v1"
description: |-
Get details about CSS flavor from OpenTelekomCloud
---

Up-to-date reference of API arguments for CSS flavor you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-search-service/api-ref/cluster_management_apis/obtaining_the_list_of_instance_flavors.html#listflavors)

# opentelekomcloud_css_flavor_v1

Use this data source to search matching CSS cluster flavor from OpenTelekomCloud.

## Example Usage

### Search by name

```hcl
data "opentelekomcloud_css_flavor_v1" "flavor" {
  name = "css.medium.8"
}
```

### Search by specs

```hcl
data "opentelekomcloud_css_flavor_v1" "flavor" {
  min_cpu = 4
  min_ram = 32
  disk_range {
    min_from = 320
    min_to   = 800
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the flavor.

* `version` - (Optional) Version of cluster.

* `type` - (Optional) Flavor type, one of `ess`, `ess-master`, `ess-client`, `ess-cold`. Default is `ess`.

* `min_cpu` - (Optional) Minimal count of CPU the flavor should have.

* `min_ram` - (Optional) Minimal RAM size (`GB`) the flavor should have.

* `disk_range` - (Optional) Disk range restrictions the flavor should match. Disk range describes available storage
  volume of the CSS node. Unit: `GB`.

  * `min_from` - (Optional) Minimal disk range start.

  * `min_to` - (Optional) Minimal disk range end.

## Attributes Reference

The following attributes of a single found flavor are exported:

* `id` - Flavor ID.

* `name` - Flavor name.

* `ram` - Flavor RAM (`GB`).

* `cpu` - Flavor CPU count.

* `region` - Region the flavor is available.

* `disk_range` - Disk range specifications.

  * `from` - Minimal disk volume the flavor can have

  * `to` - Maximum disk volume the flavor can have
