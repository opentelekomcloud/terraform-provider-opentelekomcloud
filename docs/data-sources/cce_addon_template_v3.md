---
subcategory: "Cloud Container Engine (CCE)"
---

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
