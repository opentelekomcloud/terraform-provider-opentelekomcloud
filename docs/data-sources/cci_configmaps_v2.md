---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_configmaps_v2"
sidebar_current: "docs-opentelekomcloud-data-source-cci-configmaps-v2"
description: |-
  Get the list of CCI v2 ConfigMaps within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI ConfigMap you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_configmaps_v2

Use this data source to get the list of CCI v2 ConfigMaps under a namespace within OpenTelekomCloud.

## Example Usage

```hcl
variable "namespace" {}

data "opentelekomcloud_cci_configmaps_v2" "test" {
  namespace = var.namespace
}
```

### Query a Single ConfigMap by Name

```hcl
variable "namespace" {}
variable "configmap_name" {}

data "opentelekomcloud_cci_configmaps_v2" "test" {
  namespace = var.namespace
  name      = var.configmap_name
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) Specifies the namespace to which the ConfigMaps belong.

* `name` - (Optional, String) Specifies the name of the ConfigMap used to query a single ConfigMap.
  If omitted, all ConfigMaps under the namespace are returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the ConfigMaps are queried.

* `config_maps` - The list of ConfigMaps. The [config_maps](#config_maps) structure is documented below.

<a name="config_maps"></a>
The `config_maps` block supports:

* `name` - The name of the ConfigMap.

* `namespace` - The namespace to which the ConfigMap belongs.

* `api_version` - The API version of the ConfigMap.

* `kind` - The kind of the ConfigMap.

* `annotations` - The annotations of the ConfigMap.

* `labels` - The labels of the ConfigMap.

* `creation_timestamp` - The creation timestamp of the ConfigMap.

* `resource_version` - The resource version of the ConfigMap.

* `uid` - The UID of the ConfigMap.

* `data` - The configuration data of the ConfigMap.

* `binary_data` - The binary data of the ConfigMap. Values are base64-encoded strings.

* `immutable` - Whether the ConfigMap is immutable.
