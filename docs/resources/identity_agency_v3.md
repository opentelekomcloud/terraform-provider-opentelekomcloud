---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_agency_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-agency-v3"
description: |-
  Manages a IAM Agency resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM agency you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/agency_management)

# opentelekomcloud_identity_agency_v3

Manages an agency resource within OpenTelekomcloud.

## Example Usage

```hcl
resource "opentelekomcloud_identity_agency_v3" "agency" {
  name                  = "test_agency"
  description           = "test agency"
  delegated_domain_name = "***"
  project_role {
    project = "eu-de"
    roles = [
      "KMS Administrator",
      "CCE ReadOnlyAccess",
    ]
  }
  project_role {
    all_projects = true
    roles = [
      "CES Administrator",
      "ER ReadOnlyAccess",
    ]
  }
  domain_roles = ["Anti-DDoS Administrator", ]
}
```

-> **Note**: It can not set `tenant_name` in `provider "opentelekomcloud"` when using this resource.

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The name of agency. The name is a string of 1 to 64
  characters.

* `description` - (Optional, String, ForceNew) Provides supplementary information about the
  agency. The value is a string of 0 to 255 characters.

* `delegated_domain_name` - (Required, String) The name of delegated domain.

* `project_role` - (Optional, List) An array of roles and projects which are used to
  grant permissions to agency on project. The structure is documented below.

* `domain_roles` - (Optional, List) An array of role names which stand for the
  permissions to be granted to agency on domain.

The `project_role` block supports:

* `project` - (Optional, String) The name of project
  Either `project` or `all_projects` must be provided to specify single `project_role` element.

* `roles` - (Required, List) An array of role names

* `all_projects` - (Optional, Bool) Whether roles are applied to all projects.
  Either `project` or `all_projects` must be provided to specify single `project_role` element.

-> **Note**: One or both of `project_role` and `domain_roles` must be input when creating an agency.

## Attributes Reference

The following attributes are exported:

* `id` - The agency ID.

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `delegated_domain_name` - See Argument Reference above.

* `project_role` - See Argument Reference above.

* `domain_roles` - See Argument Reference above.

* `duration` - Validity period of an agency. The default value is null,
  indicating that the agency is permanently valid.

* `expire_time` - The expiration time of agency

* `create_time` - The time when the agency was created.

## Import

Agencies can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_identity_agency_v3.this 1bc93b8b-37a4-4b50-92cc-daa4c89d4e4c
```
