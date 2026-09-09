---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_auth_scope_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-auth-scope-v3"
description: |-
  Get authentication information about the current auth scope in use within OpenTelekomCloud
---

# opentelekomcloud_identity_auth_scope_v3

Use this data source to get authentication information about the current
auth scope in use within OpenTelekomCloud. This can be used as self-discovery or introspection of
the username or project name currently in use.

## Example Usage

```hcl
data "opentelekomcloud_identity_auth_scope_v3" "scope" {
  name = "my_scope"
}
```

## Argument Reference

* `name` - (Required) The name of the scope. This is an arbitrary name which is
  only used as a unique identifier so an actual token isn't used as the ID.

## Authentication

The data source supports token, username/password, and permanent AK/SK authentication.

With token or username/password authentication, the attributes are populated from the
token representing the current authentication scope. Project attributes are empty for
a domain-scoped token, and domain attributes are empty for a project-scoped token.

With AK/SK authentication, no token is issued. The current user is resolved from the
configured `access_key`, while the project and user domain attributes are obtained from
the provider authentication context. The `roles`, `domain_id`, and `domain_name`
attributes are empty because AK/SK authentication is project-scoped.

-> Temporary AK/SK credentials configured with `security_token` are not supported.

## Attributes Reference

`id` is set to the name given to the scope. In addition, the following attributes are exported:

* `user_name` - The username of the scope.

* `user_id` - The user ID the of the scope.

* `user_domain_name` - The domain name of the user.

* `user_domain_id` - The domain ID of the user.

* `domain_name` - The domain name of the scope. Only set for a domain-scoped token, empty otherwise.

* `domain_id` - The domain ID of the scope. Only set for a domain-scoped token, empty otherwise.

* `project_name` - The project name of the scope.

* `project_id` - The project ID of the scope.

* `project_domain_name` - The domain name of the project.

* `project_domain_id` - The domain ID of the project.

* `roles` - A list of roles in the current scope. See reference below.

The `roles` block contains:

* `role_id` - The ID of the role.

* `role_name` - The name of the role.
