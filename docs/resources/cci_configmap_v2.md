---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_configmap_v2"
sidebar_current: "docs-opentelekomcloud-resource-cci-configmap-v2"
description: |-
  Manages a CCI v2 ConfigMap resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI ConfigMap you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_configmap_v2

Manages a CCI v2 ConfigMap resource within OpenTelekomCloud.

## Example Usage

### Basic ConfigMap

```hcl
resource "opentelekomcloud_cci_configmap_v2" "test" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "my-configmap"

  data = {
    "app.properties" = "debug=false\nlog_level=info"
    "timeout"        = "30s"
  }
}
```

### ConfigMap with a Certificate

```hcl
resource "opentelekomcloud_cci_configmap_v2" "tls" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "my-tls-config"

  data = {
    "ca.crt" = file("ca.crt")
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String, ForceNew) Specifies the namespace of the CCI ConfigMap.

* `name` - (Required, String, ForceNew) Specifies the name of the CCI ConfigMap.

* `data` - (Optional, Map) Specifies the configuration data of the CCI ConfigMap.
  Each value must be a UTF-8 string.

* `binary_data` - (Optional, Map) Specifies the binary data of the CCI ConfigMap.
  Each value must be a base64-encoded string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<namespace>/<name>`.

* `region` - The region of the CCI ConfigMap.

* `api_version` - The API version of the CCI ConfigMap.

* `kind` - The kind of the CCI ConfigMap.

* `annotations` - The annotations of the CCI ConfigMap.

* `labels` - The labels of the CCI ConfigMap.

* `creation_timestamp` - The creation timestamp of the CCI ConfigMap.

* `resource_version` - The resource version of the CCI ConfigMap.

* `uid` - The UID of the CCI ConfigMap.

* `immutable` - Whether the CCI ConfigMap is immutable.

## Import

The CCI v2 ConfigMap can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_cci_configmap_v2.test <namespace>/<name>
```
