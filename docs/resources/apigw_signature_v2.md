---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_signature_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-signature-v2"
description: |-
  Manages a APIGW Signature resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway signature service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/signature_key_management/index.html)

# opentelekomcloud_apigw_signature_v2

Manages a signature resource within OpenTelekomCloud.

## Example Usage

### Create a signature of the HMAC type

```hcl
variable "gateway_id" {}
variable "signature_name" {}
variable "signature_key" {}
variable "signature_secret" {}

resource "opentelekomcloud_apigw_signature_v2" "hmac" {
  instance_id = var.gateway_id
  name        = var.signature_name
  type        = "hmac"
  key         = var.signature_key
  secret      = var.signature_secret
}
```

### Create a signature and automatically generate key and secret

```hcl
variable "gateway_id" {}
variable "signature_name" {}

resource "opentelekomcloud_apigw_signature_v2" "hmac" {
  instance_id = var.gateway_id
  name        = var.signature_name
  type        = "hmac"
}
```

### Create a signature of the AES type

```hcl
variable "gateway_id" {}
variable "signature_name" {}
variable "signature_key" {}
variable "signature_secret" {}

resource "opentelekomcloud_apigw_signature_v2" "aes" {
  instance_id = var.gateway_id
  name        = var.signature_name
  type        = "aes"
  algorithm   = "aes-128-cfb"
  key         = var.signature_key
  secret      = var.signature_secret
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated gateway to which the signature
  belongs. Changing this will create a new resource.

* `name` - (Required, String) Specifies the signature name.
  The valid length is limited from `3` to `64`, only English letters, Chinese characters, digits and underscores `_` are
  allowed. The name must start with an English letter or Chinese character.

* `type` - (Required, String) Specifies the type of signature.
  The valid values are as follows:
  + `basic`: Basic auth type.
  + `hmac`: HMAC type.
  + `aes`: AES type
  Changing this will create a new resource.

* `key` - (Optional, String) Specifies the signature key.
  + For `basic` type: The value contains `4` to `32` characters, including letters, digits, underscores `_` and
    hyphens `-`. It must start with a letter.
  + For `hmac` type: The value contains `8` to `32` characters, including letters, digits, underscores `_` and
    hyphens `-`. It must start with a letter or digit.
  + For `aes` type: The value contains `16` characters if the `aes-128-cfb` algorithm is used, or `32` characters if the
    `aes-256-cfb` algorithm is used. Only letters, digits, and special characters `_-!@#$%+/=` are allowed.
    It must start with a letter, digit, plus sign `+`, or slash `/`.
  If not specified, the key will automatically generated. The auto-generation is only supported on first creation.
  Changing this will create a new resource.

* `secret` - (Optional, String) Specifies the signature secret.
  If not specified, the secret will automatically generated. The auto-generation is only supported on first creation.
  Changing this will create a new resource.
  + For `basic` type: The value contains `8` to `64` characters. Letters, digits, and special characters `_-!@#$%` are
   allowed. It must start with a letter or digit. If not specified, a value is automatically generated.
  + For `hmac` type: The value contains `16` to `64` characters. Letters, digits, and special characters `_-!@#$%` are
   allowed. It must start with a letter or digit. If not specified, a value is automatically generated.
  + For `aes` type: The value contains `16` characters, including letters, digits, and special
   characters `_-!@#$%+/=`. It must start with a letter, digit, plus sign `+`, or slash `/`. If not specified, a
   value is automatically generated.

* `algorithm` - (Optional, String) Specifies the signature algorithm.
  This parameter is required and only available when signature `type` is `aes`.
  The valid values are as follows:
  + `aes-128-cfb`
  + `aes-256-cfb`
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the signature.

* `region` - Region where the signature is located.

* `created_at` - The creation time of the signature.

* `updated_at` - The latest update time of the signature.

## Import

Signatures can be imported using their `id` and related dedicated gateway ID, separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_apigw_signature_v2.sig <gateway_id>/<id>
```
