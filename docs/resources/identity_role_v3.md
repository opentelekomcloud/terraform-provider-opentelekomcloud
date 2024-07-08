---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_role_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-role-v3"
description: |-
Manages a IAM Role resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM role you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/permission_management)

# opentelekomcloud_identity_role_v3

Custom role management

## Example Usage

```hcl
resource "opentelekomcloud_identity_role_v3" "role" {
  description   = "role"
  display_name  = "custom_role"
  display_layer = "domain"
  statement {
    effect    = "Allow"
    action    = ["obs:bucket:GetBucketAcl"]
    resource  = ["OBS:*:*:bucket:test-bucket"]
    condition = <<EOF
    {
      "StringStartWith": {
          "g:ProjectName": [
              "eu-de"
          ]
      },
      "StringNotEqualsIgnoreCase": {
          "g:ServiceName": [
              "iam"
          ]
    }
    EOF
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
  Format: `Service:*:*:resource:resource_path`.
  Examples: `KMS:*:*:KeyId:your_key`, `OBS:*:*:bucket:your_bucket`, `OBS:*:*:object:your_object`.

* `condition` - (Optional) The conditions for the permission to take effect. A maximum of 10 conditions are allowed.
  Conditions should be provided as string as in example above.

-> For the full reference checkout [Policy Syntax](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0017.html).

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
