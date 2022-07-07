---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_credential_v3

Manages permanent access key for an OpenTelekomCloud user.

## Example Usage

#### Create AK/SK for yourself
```hcl
resource opentelekomcloud_identity_credential_v3 aksk {}
```

#### Create user with AK/SK

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

## Attributes Reference

The following attributes are exported:

* `user_id` - IAM user ID. Changing this parameter will recreate the resource.

* `description` - Description of the access key.

* `status` - Status of the access key.

* `access` - Access key ID.

* `secret` - Access key secret.

* `create_time` - Time of the access key creation.

* `last_use_time` - Time of the access key last usage.
