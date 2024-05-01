---
subcategory: "APIGW"
---

# opentelekomcloud_apigw_throttling_policy_associate_v2

This API is used to bind a request throttling policy to an API that has been published in an environment within OpenTelekomCloud.
You can bind different request throttling policies to an API in different environments,
but can bind only one request throttling policy to the API in each environment.

## Example Usage

```hcl
variable "gateway_id" {}
variable "policy_id" {}
variable "publish_ids" {
  type = list(string)
}

resource "opentelekomcloud_apigw_throttling_policy_associate_v2" "tpa" {
  gateway_id  = var.gateway_id
  policy_id   = var.policy_id
  publish_ids = var.publish_ids
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated gateway to which the APIs and the
  throttling policy belongs.
  Changing this will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the ID of the throttling policy.
  Changing this will create a new resource.

* `publish_ids` - (Required, List) Specifies the publishing IDs corresponding to the APIs bound by the throttling policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The format is `<instance_id>/<policy_id>`.

* `region` - Specifies the region where the dedicated instance and the throttling policy are located.

## Import

Resources can be imported using their `policy_id` and the APIGW dedicated gateway ID to which the policy
belongs, separated by a slash, e.g.

```shell
$ terraform import opentelekomcloud_apigw_throttling_policy_associate_v2.tpa <instance_id>/<policy_id>
```
