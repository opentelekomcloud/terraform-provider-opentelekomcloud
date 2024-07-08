---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_application_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-application-v2"
description: |-
Manages a APIGW Application resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway App service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/credential_management/index.html)

# opentelekomcloud_apigw_application_v2

Manages an APIGW application resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "gateway_id" {}
variable "app_name" {}
variable "app_code" {}

resource "opentelekomcloud_apigw_application_v2" "test" {
  gateway_id  = var.gateway_id
  name        = var.app_name
  description = "Created by script"

  app_codes = [var.app_code]
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the application
  belongs.
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the application name.
  The valid length is limited from can contain `3` to `64`, only letters, digits and hyphens `-` are allowed.
  The name must start with a letter.

* `description` - (Optional, String) Specifies the application description.
  The description contain a maximum of 255 characters and the angle brackets (`<` and `>`) are not allowed.

  -> The description does not support updating to an empty value.

* `app_codes` - (Optional, List) Specifies an array of one or more application codes that the application has.
  Up to five application codes can be created.
  The valid length of each application code is limited from can contain `64` to `180`.
  The application code must start with a letter, digit, plus sign `+` or slash `/`.
  Only letters, digits and following special characters are allowed: `!@#$%+-_/=`.

* `secret_action` - (Optional, String) Specifies the secret action to be done for the application.
  The valid action is `RESET`.

  -> The `secret_action` is a one-time action.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The application ID.

* `region` - The region where the application is located.

* `registration_time` - The registration time.

* `updated_at` - The latest update time of the application.

* `app_key` - App key.

* `app_secret` - App secret.

## Import

Applications can be imported using their `id` and the ID of the related dedicated instance, separated by a slash, e.g.

```shell
$ terraform import opentelekomcloud_apigw_application_v2.app <gateway_id>/<id>
```
