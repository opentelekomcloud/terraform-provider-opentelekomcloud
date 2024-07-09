---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_environment_variable_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-environment-variable-v2"
description: |-
Manages a APIGW Environment variable resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway environment variable service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/environment_variable_management/index.html)

# opentelekomcloud_apigw_environment_variable_v2

Manages an APIGW environment variable resource within OpenTelekomCloud.

-> A maximum of `50` variables can be created on the same environment.

## Example Usage

```hcl
variable "gateway_id" {}
variable "environment_id" {}
variable "group_id" {}
variable "variable_name" {}
variable "variable_value" {}

resource "opentelekomcloud_apigw_environment_variable_v2" "var" {
  gateway_id     = var.gateway_id
  environment_id = var.environment_id
  group_id       = var.group_id
  name           = var.variable_name
  value          = var.variable_value
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated gateway instance to which the environment
  variable belongs. Changing this creates a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the group to which the environment variable belongs.
  Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the environment to which the environment variable belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the environment variable.
  The valid length is limited from `3` to `32` characters.
  Only letters, digits, hyphens `-`, and underscores `_` are allowed, and must start with a letter.
  In the definition of an API, the `name` (`case-sensitive`) indicates a variable, for example, `#Name#`.
  It is replaced by the actual value when the API is published in an environment. The variable name must be unique.
  Changing this creates a new resource.

* `value` - (Required, String, ForceNew) Specifies the value of the environment variable.
  The valid length is limited from `1` to `255` characters. Only letters, digits and special characters `_-/.:` are allowed.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `region` - The region where the dedicated instance is located.

## Import

The resource can be imported using `gateway_id`, `group_id` and `name`, separated by slashes (/), e.g.

```bash
$ terraform import opentelekomcloud_apigw_environment_variable_v2.test <gateway_id>/<group_id>/<name>
```
