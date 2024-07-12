---
subcategory: "Cloud Container Engine (CCE)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_addon_templates_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cce-addon-templates-v3"
description: |-
  Get CCE Addon template versions and additional info based on cluster version from OpenTelekomCloud
---

Up-to-date reference of API arguments for CCE Addon template you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-engine/api-ref/apis/add-on_management/reading_add-on_templates.html#cce-02-0321)

# opentelekomcloud_cce_addon_templates_v3

Use this data source to get from OpenTelekomCloud a CCE Addon template versions and additional info based on cluster version.

## Example Usage

```hcl
data "opentelekomcloud_cce_addon_templates_v3" "templates" {
  cluster_version = "1.21.3"
  addon_name      = "volcano"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_version` -  (Required) The version of the CCE cluster. For example: `1.23.6`.

* `addon_name` - (Required) The name of the CCE addon. For example: `autoscaler`.

* `cluster_type` - (Optional) The type of the CCE cluster. Default value: `VirtualMachine`.
  The valid values are as follows:
    + **VirtualMachine**: The instance is running properly.
    + **ARM64**: The instance has been properly stopped.
    + **BareMetal**: An error has occurred on the instance.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `addons` - List of CCE addons details. The object structure of each CCE addon is documented below.

The `addons` block supports:

* `cluster_ip` - The cluster ip.

* `image_version` - The cluster image version.

* `platform` - The cluster image version.

* `euleros_version` - The euler os version.

* `obs_url` - The obs endpoint url.

* `swr_addr` - The cluster `swr_addr`.

* `swr_user` - The cluster `swr_user`.

* `addon_version` - Supported addon version.
