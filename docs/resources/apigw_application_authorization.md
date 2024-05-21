---
subcategory: "APIGW"
---

# opentelekomcloud_apigw_application_authorization_v2

Using this resource to authorize APIs for application, allowing it to access the published APIs within OpenTelekomCloud.

-> For an application, an environment can only create one `opentelekomcloud_apigw_application_authorization_v2` resource (all
   published APIs must belong to an environment).

## Example Usage

```hcl
variable "gateway_id" {}
variable "application_id" {}
variable "published_env_id" {}
variable "published_api_ids" {
  type = list(string)
}

resource "opentelekomcloud_apigw_application_authorization_v2" "auth" {
  gateway_id     = var.gateway_id
  application_id = var.application_id
  env_id         = var.published_env_id
  api_ids        = var.published_api_ids
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated gateway instance to which the application
  and APIs belong.
  Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the ID of the application authorized to access the APIs.
  Changing this will create a new resource.

* `env_id` - (Required, String, ForceNew) Specifies the environment ID where the APIs were published.
  Changing this will create a new resource.

* `api_ids` - (Required, List) Specifies the authorized API IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also `<env_id>/<application_id>`.

* `region` - Region where the application and APIs are located.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `update` - Default is 3 minutes.
* `delete` - Default is 3 minutes.

## Import

Auth of application can be imported using related `gateway_id` and their `id` (also consists of
`env_id` and `application_id`), separated by the slashes, e.g.

```bash
$ terraform import opentelekomcloud_apigw_application_authorization_v2.auth <gateway_id>/<env_id>/<application_id>
```
