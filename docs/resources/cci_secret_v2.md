---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_secret_v2"
sidebar_current: "docs-opentelekomcloud-resource-cci-secret-v2"
description: |-
  Manages a CCI v2 Secret resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI secret you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_secret_v2

Manages a CCI v2 Secret resource within OpenTelekomCloud.

## Example Usage

### Opaque Secret

```hcl
resource "opentelekomcloud_cci_secret_v2" "opaque" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "my-opaque-secret"
  type      = "Opaque"

  data = {
    "username" = base64encode("admin")
    "password" = base64encode("p@ssw0rd")
  }
}
```

### Docker Registry Secret

```hcl
resource "opentelekomcloud_cci_secret_v2" "registry" {
  namespace = opentelekomcloud_cci_namespace_v2.test.name
  name      = "my-registry-secret"
  type      = "kubernetes.io/dockerconfigjson"

  data = {
    ".dockerconfigjson" = base64encode(jsonencode({
      auths = {
        "swr.eu-de.otc.t-systems.com" = {
          username = "user"
          password = "token"
          auth     = base64encode("user:token")
        }
      }
    }))
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String, ForceNew) Specifies the namespace of the CCI Secret.

* `name` - (Required, String, ForceNew) Specifies the name of the CCI Secret.

* `string_data` - (Optional, Map) Specifies string data of the CCI Secret.
  The values will be encoded to base64 by the API before storage.

* `data` - (Optional, Map) Specifies the data of the CCI Secret.
  The values must be base64-encoded strings.

* `type` - (Optional, String) Specifies the type of the CCI Secret.
  For example, `Opaque`, `kubernetes.io/dockerconfigjson`, etc.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<namespace>/<name>`.

* `region` - The region of the CCI Secret.

* `api_version` - The API version of the CCI Secret.

* `kind` - The kind of the CCI Secret.

* `annotations` - The annotations of the CCI Secret.

* `labels` - The labels of the CCI Secret.

* `creation_timestamp` - The creation timestamp of the CCI Secret.

* `resource_version` - The resource version of the CCI Secret.

* `uid` - The uid of the CCI Secret.

* `immutable` - Whether the CCI Secret is immutable.

## Import

The CCI v2 Secret can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_cci_secret_v2.test <namespace>/<name>
```
