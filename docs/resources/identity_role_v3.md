---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_role_v3

Custom role management

## Example Usage

```hcl
resource "opentelekomcloud_identity_role_v3" "role" {
  description   = "role"
  display_name  = "%s"
  display_layer = "domain"
  statement {
    effect   = "Allow"
    action   = ["obs:bucket:GetBucketAcl"]
    resource = "obs:*:*:bucket:test-bucket"
  }
  statement {
    effect = "Allow"
    action = [
      "obs:bucket:HeadBucket",
      "obs:bucket:ListBucketMultipartUploads",
      "obs:bucket:ListBucket"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) Description of a role. The value cannot exceed 256 characters.

* `display_layer` - (Required) Display layer of a role.
  * `domain` - A role is displayed at the domain layer.
  * `project` - A role is displayed at the project layer.

* `display_name` - (Required) Displayed name of a role. The value cannot exceed 64 characters.

* `statement` - (Required) Statement: The Statement field contains the Effect and Action
  elements. Effect indicates whether the policy allows or denies
  access. Action indicates authorization items. The number of
  statements cannot exceed 8. Structure is documented below.

The `statement` block supports:

* `action` - (Required) Permission set, which specifies the operation permissions on
  resources. The number of permission sets cannot exceed 100.
  Format:  The value format is Service name:Resource type:Action,
  for example, vpc:ports:create. Service name: indicates the
  product name, such as ecs, evs, or vpc. Only lowercase letters
  are allowed. Resource type and Action: The values are
  case-insensitive, and the wildcard (*) are allowed. A wildcard
  (*) can represent all or part of information about resource
  types and actions for the specific service.

* `effect` - (Required) The value can be Allow and Deny. If both Allow and Deny are
  found in statements, the policy evaluation starts with Deny.

* `resource` - (Optional) The resources which will be granted/denied accesses.
  Format: `service:*:*:resource:resource_path`.
  Examples: `KMS:*:*:KeyId:your_key`, `OBS:*:*:bucket:your_bucket`, `OBS:*:*:object:your_object`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `catalog` - Directory where a role locates

* `domain_id` - ID of the domain to which a role belongs

* `name` - Name of a role

## Import

Role can be imported using the following format:

```sh
terraform import opentelekomcloud_identity_role_v3.default {{ resource id}}
```
