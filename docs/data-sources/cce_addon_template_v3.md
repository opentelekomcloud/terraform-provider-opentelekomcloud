---
subcategory: "Cloud Container Engine (CCE)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_addon_template_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cce-addon-template-v3"
description: |-
Get CCE Addon template info from OpenTelekomCloud
---

Up-to-date reference of API arguments for CCE Addon template you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-engine/api-ref/apis/add-on_management/reading_add-on_templates.html#cce-02-0321)

# opentelekomcloud_cce_addon_template_v3

Use this data source to get from OpenTelekomCloud a CCE Addon template info.

## Example Usage

```hcl
data "opentelekomcloud_cce_addon_template_v3" "template" {
  addon_version = "1.23.1"
  addon_name    = "coredns"
}
```

## Argument Reference

The following arguments are supported:

* `addon_version` -  (Required) The version of the CCE cluster addon. For example: `1.23.6`.

* `addon_name` - (Required) The name of the CCE addon. For example: `autoscaler`.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `id` - The ID of the addon.

* `cluster_ip` - The cluster ip.

* `image_version` - The cluster image version.

* `swr_addr` - The cluster `swr_addr`.

* `swr_user` - The cluster `swr_user`.

* `cluster_versions` - Supported cluster versions.
