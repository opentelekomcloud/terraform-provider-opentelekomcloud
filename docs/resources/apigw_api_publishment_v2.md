---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_api_publishment_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-api-publishment-v2"
description: |-
Manages a APIGW API publishment resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway Api publishment you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/api_management/publishing_an_api_or_taking_an_api_offline.html#createordeletepublishrecordforapiv2-1)

# opentelekomcloud_apigw_api_publishment_v2

Using this resource to publish an API to the environment or manage a historical publishing version within OpenTelekomCloud.

~> If you republish on the same environment or switch versions through other ways (such as console) after the API is
published through terraform, the current resource attributes will be affected, resulting in data inconsistency.

## Example Usage

### Publish a new version of the API

```hcl
variable "gateway_id" {}
variable "environment_id" {}
variable "api_id" {}

resource "opentelekomcloud_apigw_api_publishment_v2" "default" {
  gateway_id     = var.gateway_id
  environment_id = var.environment_id
  api_id         = var.api_id
}
```

### Switch to a specified version of the API which is published

```hcl
variable "gateway_id" {}
variable "environment_id" {}
variable "api_id" {}
variable "version_id" {}

resource "opentelekomcloud_apigw_api_publishment_v2" "default" {
  gateway_id     = var.gateway_id
  environment_id = var.environment_id
  api_id         = var.api_id
  version_id     = var.version_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIGW dedicated instance to which the API belongs
  to. Changing this will create a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the environment to which the current version of the API
  will be published or has been published.
  Changing this will create a new resource.

* `api_id` - (Required, String, ForceNew) Specifies the ID of the API to be published or already published.
  Changing this will create a new resource.

* `description` - (Optional, String) Specifies the description of the current publish.

* `version_id` - (Optional, String) Specifies the version ID of the current publish.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is constructed from the instance ID, environment ID, and API ID, separated by slashes.

* `environment_name` - The name of the environment to which the current version of the API is published.

* `published_at` - Time when the current version was published.

* `publish_id` - The publishing ID of the API in current environment.

* `history` - All publish history of the API.
  The [object](#publishment_history) structure is documented below.

* `region` - The region in which to APIs was published.

<a name="publishment_history"></a>
The `history` block supports:

* `version_id` - The version ID of the API publish.

* `description` - The version description of the API publish.

## Import

The publishment can be imported using related `instance_id`, `environment_id` and `api_id`, separated by slashes, e.g.

```shell
$ terraform import opentelekomcloud_apigw_api_publishment_v2.pub <instance_id>/<environment_id>/<api_id>
```
