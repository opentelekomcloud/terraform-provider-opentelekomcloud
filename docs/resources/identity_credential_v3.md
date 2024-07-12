---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_credential_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-credential-v3"
description: |-
  Manages a IAM Credential resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM credential you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/access_key_management)

# opentelekomcloud_identity_credential_v3

Manages permanent access key for an OpenTelekomCloud user.

## Example Usage

### Create AK/SK for yourself
```hcl
resource opentelekomcloud_identity_credential_v3 aksk {}
```

### Create user with AK/SK

```hcl
resource opentelekomcloud_identity_user_v3 user {
  name     = "user_1"
  password = "password123!"
}

resource opentelekomcloud_identity_credential_v3 aksk {
  user_id     = opentelekomcloud_identity_user_v3.user.id
  description = "Created by administrator"
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Optional) IAM user ID. If not set, will create AK/SK for yourself.

* `description` - (Optional) Description of the access key.

* `status` - (Optional) Status of the access key to be changed to. The value can be `active` or `inactive`.

* `pgp_key` - (Optional, String, ForceNew) Either a base-64 encoded PGP public key, or a keybase username in the form
  `keybase:some_person_that_exists`. Changing this creates a new resource.

## Attributes Reference

The following attributes are exported:

* `user_id` - IAM user ID. Changing this parameter will recreate the resource.

* `description` - Description of the access key.

* `status` - Status of the access key.

* `access` - Access key ID.

* `secret` - Secret key ID. If `pgp_key` is not set, secret will be in plain text.
  The encrypted secret, base64 encoded. The encrypted secret may be decrypted using the command
  line, for example: `terraform output encrypted_secret | base64 --decode | keybase pgp decrypt`.

* `create_time` - Time of the access key creation.

* `last_use_time` - Time of the access key last usage.
