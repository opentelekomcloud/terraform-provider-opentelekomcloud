---
subcategory: "Config"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rms_resource_tags_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rms-resource-tags-v1"
description: |-
  Manages an RMS resource tags data source, used to query resource tags visible to you, within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RMS resource tags you can get at
[documentation portal](https://docs.otc.t-systems.com/config/api-ref/apis/resource_query/querying_resource_tags.html#rms-04-0106)


# opentelekomcloud_rms_resource_tags_v1

Manages an RMS resource tags data source, used to query all resource tags under your account, within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_rms_resource_tags_v1" "test" {}
```

## Argument Reference

The following arguments are supported:

* `key` - (Optional, String) Specifies the name of the tag key. Maximum length: `128`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tags` - Specifies the list of tags. The structure is documented below:
    * `key` - Specifies the tag key.
    * `value` - Specifies the tag values.
