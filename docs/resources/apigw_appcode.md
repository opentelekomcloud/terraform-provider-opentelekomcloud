---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_appcode_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-appcode-v2"
description: |-
Manages a APIGW AppCode resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway AppCode service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/credential_management/index.html)

# opentelekomcloud_apigw_appcode_v2

Manages an APIGW AppCode in application resource within OpenTelekomCloud.

## Example Usage

### Auto generate AppCode

```hcl
variable "gateway_id" {}
variable "application_id" {}

resource "opentelekomcloud_apigw_appcode_v2" "code" {
  gateway_id     = var.gateway_id
  application_id = var.application_id
}
```

### Manually configure AppCode

```hcl
variable "gateway_id" {}
variable "application_id" {}
variable "app_code" {}

resource "opentelekomcloud_apigw_appcode_v2" "code" {
  gateway_id     = var.gateway_id
  application_id = var.application_id
  value          = var.app_code
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the application
  and AppCode belong. Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the ID of application to which the AppCode belongs.
  Changing this will create a new resource.

* `value` - (Optional, String, ForceNew) Specifies the AppCode value (`content`).
  The value can contain `64` to `180` characters, starting with a letter, plus sign `+`, or slash `/`, or digit.
  Only letters, digit and the following special characters are allowed: `+_!@#$%/=`.
  If omitted, a random value will be generated. Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The AppCode ID.

* `region` - The region where the application and AppCode are located.

* `created_at` - The creation time of the AppCode.

## Import

AppCode can be imported using related `gateway_id`, `application_id` and their `id`, separated by slashes, e.g.

```bash
$ terraform import opentelekomcloud_apigw_appcode_v2.code <gateway_id>/<application_id>/<id>
```
