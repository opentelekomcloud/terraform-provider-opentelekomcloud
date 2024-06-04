---
subcategory: "APIGW"
---

Up-to-date reference of API arguments for API Gateway Custom Authorizer service you can get at
`https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/custom_authorizer_management/index.html`.

# opentelekomcloud_apigw_custom_authorizer_v2

Manages an APIGW custom authorizer resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "gateway_id" {}
variable "authorizer_name" {}
variable "function_urn" {}

resource "opentelekomcloud_apigw_custom_authorizer_v2" "test" {
  gateway_id   = var.gateway_id
  name         = var.authorizer_name
  function_urn = var.function_urn
  type         = "FRONTEND"
  cache_age    = 60

  identity {
    name     = "user_name"
    location = "QUERY"
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies an ID of the APIGW dedicated instance to which the
  custom authorizer belongs to.
  Changing this will create a new custom authorizer resource.

* `name` - (Required, String) Specifies the name of the custom authorizer.
  The custom authorizer name consists of `3` to `64` characters, starting with a letter.
  Only letters, digits and underscores `_` are allowed.

* `function_urn` - (Required, String) Specifies the uniform function URN of the function graph resource.

* `type` - (Optional, String, ForceNew) Specifies the custom authorize type.
  The valid values are `FRONTEND` and `BACKEND`. Defaults to `FRONTEND`.
  Changing this will create a new custom authorizer resource.

* `is_body_send` - (Optional, Bool) Specifies whether to send the body.

* `ttl` - (Optional, Int) Specifies the maximum cache age.

* `user_data` - (Optional, String) Specifies the user data, which can contain a maximum of `2,048` characters.
  The user data is used by APIGW to invoke the specified authentication function when accessing the backend service.

  -> **NOTE:** The user data will be displayed in plain text on the console.

* `identity` - (Optional, List) Specifies an array of one or more parameter identities of the custom authorizer.
  The [object](#authorizer_identity) structure is documented below.

<a name="authorizer_identity"></a>
The `identity` block supports:

* `name` - (Required, String) Specifies the name of the parameter to be verified.
  The parameter includes front-end and back-end parameters.

* `location` - (Required, String) Specifies the parameter location, which support `HEADER` and `QUERY`.

* `validation` - (Optional, String) Specifies the parameter verification expression.
  If omitted, the custom authorizer will not perform verification.
  The valid value is range form `1` to `2,048`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the custom authorizer.

* `created_at` - The creation time of the custom authorizer.

* `region` - The region in which to create the custom authorizer resource.

## Import

Custom Authorizers of the APIGW can be imported using their `name` and related dedicated gateway IDs, separated by a
slash, e.g.

```shell
$ terraform import opentelekomcloud_apigw_custom_authorizer_v2.test <gateway_id>/<name>
```
