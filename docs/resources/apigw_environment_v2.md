---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_environment_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-environment-v2"
description: |-
Manages a APIGW Environment resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway environment service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/index.html)

# opentelekomcloud_apigw_environment_v2

API Gateway (APIGW) is a high-performance, high-availability, and high-security API hosting service that helps you build,
manage, and deploy APIs at any scale.
With just a few clicks, you can integrate internal systems, and selectively expose capabilities with minimal costs and risks.

## Example Usage

```hcl
variable "instance_id" {}
variable "environment_name" {}
variable "description" {}

resource "opentelekomcloud_apigw_environment_v2" "test" {
  instance_id = var.instance_id
  name        = var.environment_name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dedicated instance is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the environment
  belongs.
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the environment name.
  The valid length is limited from `3` to `64`, only letters, digits and underscores (_) are allowed.
  The name must start with a letter.

* `description` - (Optional, String) Specifies the environment description.
  The value can contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the dedicated environment.

* `created_at` - The time when the environment was created.

## Import

Environments can be imported using their `name` and the ID of the related dedicated instance, separated by a slash, e.g.

```
$ terraform import opentelekomcloud_apigw_environment_v2.test instance_id/name
```
