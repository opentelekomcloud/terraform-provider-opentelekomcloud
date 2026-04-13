---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_namespaces_v2"
sidebar_current: "docs-opentelekomcloud-data-source-cci-namespaces-v2"
description: |-
  Get the list of CCI v2 namespaces within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI namespace you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_namespaces_v2

Use this data source to get the list of CCI v2 namespaces within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_cci_namespaces_v2" "all" {}
```

```hcl
variable "name" {}

data "opentelekomcloud_cci_namespaces_v2" "by_name" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the namespace used to query the namespace detail.
  If omitted, the list of all namespaces is returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the namespaces are queried.

* `namespaces` - The list of namespaces. The [namespaces](#namespaces) structure is documented below.

<a name="namespaces"></a>
The `namespaces` block supports:

* `name` - The name of the namespace.

* `api_version` - The API version of the namespace.

* `kind` - The kind of the namespace.

* `annotations` - The annotations of the namespace.

* `labels` - The labels of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `finalizers` - The finalizers of the namespace.

* `resource_version` - The resource version of the namespace.

* `uid` - The uid of the namespace.

* `status` - The status of the namespace.
