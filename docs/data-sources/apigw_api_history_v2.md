---
subcategory: "APIG"
---

# opentelekomcloud_apigw_api_history_v2

This API is used to query the historical versions of an API within OpenTelekomCloud.
APIGW retains a maximum of 10 historical versions for each API in an environment.

## Example Usage

```hcl
variable "gateway_id" {}
variable "environment_id" {}
variable "api_id" {}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub" {
  gateway_id     = var.instance_id
  environment_id = var.environment_id
  api_id         = var.api_id
}

data "opentelekomcloud_apigw_api_history_v2" "hist" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id

  depends_on = ["opentelekomcloud_apigw_api_publishment_v2.pub"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the APIG environment list.
  If omitted, the provider-level region will be used.

* `gateway_id` - (Required, String) Specifies an ID of the APIGW dedicated instance to which the API
  environment belongs.

* `api_id` - (Required, String) Specifies the ID of the API to be published or already published.

* `environment_id` - (Optional, String) Specifies the environment ID.

* `environment_name` - (Optional, String) Specifies the environment name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `history` - List of APIG history details. The structure is documented below.

The `history` block supports:

* `id` - API version ID.

* `name` - API version name.

* `description` - The description about the API publication.

* `publish_time` - Time when the APIG publication was created, in RFC-3339 format.

* `status` - Version status.
  Values:
    `1`: `effective`
    `2`: `not effective`
