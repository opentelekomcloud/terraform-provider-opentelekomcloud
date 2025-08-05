---
subcategory: "Config"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rms_resource_relationships_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rms-resource-relationships-v1"
description: |-
  Manages an RMS resource relationships data source, used to query the relationship between a resource and other resources by the resource ID, within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RMS resource relationships you can get at
[documentation portal](https://docs.otc.t-systems.com/config/api-ref/apis/resource_relationships/index.html)


# opentelekomcloud_rms_resource_relationships_v1

Manages an RMS resource relationships data source, used to query the relationship between a resource and other resources by the resource ID, within OpenTelekomCloud.

  -> **NOTE:**
  Resource recorder must be enabled to query resource relationships.

## Example Usage

```hcl
variable "resource_id" {}

data "opentelekomcloud_rms_resource_relationships_v1" "relations_1" {
  resource_id = var.resource_id
  direction   = "in"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Specifies the resource ID. Maximum length: `512`.

* `direction` - (Required, String) Specifies the direction of a resource relationship. Permitted values: `in`, `out`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `relations` - Specifies the list of the resource relationships. The structure is documented below:
    * `relation_type` - Specifies the relationship type.
    * `from_resource_type` -Specifies the type of the source resource.
    * `to_resource_type` - Specifies the type of the destination resource.
    * `from_resource_id` - Specifies the ID of the source resource.
    * `to_resource_id` - Specifies the ID of the destination resource.
