---
subcategory: "APIGW"
---

Up-to-date reference of API arguments for API Gateway Acl associate service you can get at
`https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/binding_unbinding_access_control_policies/index.html`.

# opentelekomcloud_apigw_acl_policy_associate_v2

Use this resource to bind the APIs to the ACL policy within OpenTelekomCloud.

-> An ACL policy can only create one `opentelekomcloud_apigw_acl_policy_associate_v2` resource.

## Example Usage

```hcl
variable "gateway_id" {}
variable "policy_id" {}
variable "api_publish_ids" {
  type = list(string)
}

resource "opentelekomcloud_apigw_acl_policy_associate_v2" "test" {
  gateway_id  = var.gateway_id
  policy_id   = var.policy_id
  publish_ids = var.api_publish_ids
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated gateway instance to which the APIs and the
  ACL policy belong. Changing this will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the ACL Policy ID for APIs binding.
  Changing this will create a new resource.

* `publish_ids` - (Required, List) Specifies the publishing IDs corresponding to the APIs bound by the ACL policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The format is `<gateway_id>/<policy_id>`.

* `region` - Specifies the region where the dedicated instance and the throttling policy are located.

## Import

Associate resources can be imported using their `policy_id` and the APIG gateway instance ID to which the policy
belongs, separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_apigw_acl_policy_associate_v2.test <gateway_id>/<policy_id>
```
