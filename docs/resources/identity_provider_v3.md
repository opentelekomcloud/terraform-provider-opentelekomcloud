---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_provider_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-provider-v3"
description: |-
Manages a IAM Provider v3 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM provider you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/federated_identity_authentication_management/identity_provider)

# opentelekomcloud_identity_provider_v3

-> You _must_ have security admin privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).


## Example Usage

```hcl
resource "opentelekomcloud_identity_provider_v3" "provider" {
  name        = "ACME"
  description = "This is simple identity provider"
  enabled     = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (ID) of the provider. Changing this creates a new provider.

* `description` - (Optional) A description of the provider.

* `enabled` - (Optional) Whether an identity provider is enabled. Default value is `false`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

`links` - Resource links of an identity provider, including `protocols` and `self`.

`remote_ids` - Federated user ID list of an identity provider.

## Import

Providers can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_identity_provider_v3.provider ACME
```

