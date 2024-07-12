---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_role_custom_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-role-custom-v3"
description: |-
  Get a IAM user-defined role information from OpenTelekomCloud
---

Up-to-date reference of API arguments for IAM role you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/permission_management/querying_role_details.html)

# opentelekomcloud_identity_role_custom_v3

Use this data source to get the info of custom OpenTelekomCloud role.

-> For pre-defined user roles usage please refer to `opentelekomcloud_identity_role_v3`

## Example Usage

### Querying custom role by `display_name`

```hcl
data "opentelekomcloud_identity_role_custom_v3" "auth_admin" {
  display_name = "my-custom-policy"
}
```

### Querying custom role by resource `id`

```hcl
data "opentelekomcloud_identity_role_custom_v3" "auth_admin" {
  id = "13f0e753101649699664672d7b7af752"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The name of the role.

* `type` - (Optional) Display layer of a role.
    * `domain` - A role is displayed at the domain layer.
    * `project` - A role is displayed at the project layer.

* `id` - (Optional) The `id` of custom role.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `description` - Description of a role.

* `catalog` - Directory where a role locates

* `domain_id` - ID of the domain to which a role belongs

* `name` - Name of a role

* `statement` - Statement: The Statement field contains the Effect and Action
  elements. Effect indicates whether the policy allows or denies
  access. Action indicates authorization items. The number of
  statements cannot exceed 8. Structure is documented below.

The `statement` block supports:

* `action` - Permission set, which specifies the operation permissions on
  resources. The number of permission sets cannot exceed 100.
  Format:  The value format is Service name:Resource type:Action,
  for example, vpc:ports:create. Service name: indicates the
  product name, such as ecs, evs, or vpc. Only lowercase letters
  are allowed. Resource type and Action: The values are
  case-insensitive, and the wildcard (*) are allowed. A wildcard
  (*) can represent all or part of information about resource
  types and actions for the specific service.

* `effect` - The value can be Allow and Deny. If both Allow and Deny are
  found in statements, the policy evaluation starts with Deny.

* `resource` -  The resources which will be granted/denied accesses.
  Format: `Service:*:*:resource:resource_path`.
  Examples: `KMS:*:*:KeyId:your_key`, `OBS:*:*:bucket:your_bucket`, `OBS:*:*:object:your_object`.

* `condition` - The conditions for the permission to take effect.
