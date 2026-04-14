---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_secrets_v2"
sidebar_current: "docs-opentelekomcloud-data-source-cci-secrets-v2"
description: |-
  Get the list of CCI v2 Secrets within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI Secret you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_secrets_v2

Use this data source to get the list of CCI v2 Secrets under a namespace within OpenTelekomCloud.

## Example Usage

```hcl
variable "namespace" {}

data "opentelekomcloud_cci_secrets_v2" "test" {
  namespace = var.namespace
}
```

### Query a Single Secret by Name

```hcl
variable "namespace" {}
variable "secret_name" {}

data "opentelekomcloud_cci_secrets_v2" "test" {
  namespace = var.namespace
  name      = var.secret_name
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) Specifies the namespace to which the Secrets belong.

* `name` - (Optional, String) Specifies the name of the Secret used to query a single Secret.
  If omitted, all Secrets under the namespace are returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the Secrets are queried.

* `secrets` - The list of Secrets. The [secrets](#secrets) structure is documented below.

<a name="secrets"></a>
The `secrets` block supports:

* `name` - The name of the Secret.

* `namespace` - The namespace to which the Secret belongs.

* `api_version` - The API version of the Secret.

* `kind` - The kind of the Secret.

* `annotations` - The annotations of the Secret.

* `labels` - The labels of the Secret.

* `creation_timestamp` - The creation timestamp of the Secret.

* `resource_version` - The resource version of the Secret.

* `uid` - The UID of the Secret.

* `string_data` - The non-binary string data of the Secret. Values are plain UTF-8 strings.

* `data` - The data of the Secret. Values are base64-encoded strings.

* `type` - The type of the Secret.

* `immutable` - Whether the Secret is immutable.
