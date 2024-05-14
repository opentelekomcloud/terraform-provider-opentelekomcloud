---
subcategory: "APIGW"
---

Up-to-date reference of API arguments for API Gateway service you can get at
`https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/index.html`.

# opentelekomcloud_apigw_group_v2

API Gateway (APIG) is a high-performance, high-availability, and high-security API hosting service that helps you build,
manage, and deploy APIs at any scale.
With just a few clicks, you can integrate internal systems, and selectively expose capabilities with minimal costs and risks.

## Example Usage

```hcl
resource "opentelekomcloud_apigw_group_v2" "group" {
  instance_id = var.gateway_id
  name        = "group-name"
  description = "test description"

  environment {
    variable {
      name  = "test-name"
      value = "test-value"
    }
    environment_id = var.env_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APIGW group is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the group belongs.
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the group name.
  The valid length is limited from `3` to `64`, only letters, digits and hyphens (-) are
  allowed.

* `description` - (Optional, String) Specifies the group description.

* `environment` - (Optional, List) Specifies an array of one or more environments of the associated group.
  The `environment` structure is documented below.

The `environment` block supports:

* `variable` - (Required, List) Specifies an array of one or more environment variables.
  The `variable` structure is documented below.

  -> The environment variables of different groups are isolated in the same environment.

* `environment_id` - (Required, String) Specifies the environment ID of the associated group.

The `variable` block supports:

* `name` - (Required, String) Specifies the variable name.
  The valid length is limited from `3` to `32` characters.
  Only letters, digits, hyphens (-), and underscores (_) are allowed, and must start with a letter.

* `value` - (Required, String) Specifies the variable value.
  The valid length is limited from `1` to `255` characters.
  Only letters, digits and special characters (_-/.:) are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The group ID.

* `registration_time` - The registration time, in RFC-3339 format.

* `updated_at` - The time when the API group was last modified, in RFC-3339 format.

* `environment` - The array of one or more environments of the associated group.
  The `environment` structure is documented below.

The `environment` block supports:

* `variable` - The array of one or more environment variables.
  The `variable` structure is documented below.

The `variable` block supports:

* `id` - The variable ID.

## Import

API groups can be imported using their `id` and the ID of the related dedicated instance, separated by a slash, e.g.

```shell
$ terraform import opentelekomcloud_apigw_group_v2.test <instance_id>/<id>
```
