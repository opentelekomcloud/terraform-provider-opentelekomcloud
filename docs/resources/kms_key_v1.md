---
subcategory: "Key Management Service (KMS)"
---

# opentelekomcloud_kms_key_v1

Manages a V1 KMS key resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_kms_key_v1" "key_1" {
  key_alias       = "key_1"
  pending_days    = "7"
  key_description = "first test key"
  realm           = "eu-de-01"
  is_enabled      = true

  tags = {
    muh = "kuh"
  }
}
```

## Argument Reference

The following arguments are supported:

* `allow_cancel_deletion` - (Optional) Specifies whether the key is enabled from Pending Deletion state. The value `true` indicates
  that the key state Pending Deletion will be cancelled.

* `key_alias` - (Required) The alias in which to create the key. It is required when
  we create a new key. Changing this updates the alias of key.

* `key_description` - (Optional) The description of the key as viewed in OpenTelekomCloud console.
  Changing this updates the description of key.

* `realm` - (Optional) Region where a key resides. Changing this creates a new key.

* `pending_days` - (Optional) Duration in days after which the key is deleted
  after destruction of the resource, must be between 7 and 1096 days. Defaults to 7.
  It only is used when delete a key.

* `is_enabled` - (Optional) Specifies whether the key is enabled. Defaults to true.
  Changing this updates the state of existing key.

* `rotation_interval` - (Optional) Rotation interval. The value is an integer ranging from 30 to 365.
   Set the interval based on how often a CMK is used.
   If it is frequently used, set a short interval; otherwise, set a long one.

* `rotation_enabled` - (Optional) Specifies whether the key is enabled for rotation.

* `tags` - (Optional) Tags key/value pairs to associate with the AutoScaling Group.


## Attributes Reference

The following attributes are exported:

* `id` - The globally unique identifier for the key.

* `key_alias` - See Argument Reference above.

* `key_description` - See Argument Reference above.

* `realm` - See Argument Reference above.

* `default_key_flag` - Identification of a Master Key. The value `1` indicates a Default
  Master Key, and the value `0` indicates a key.

* `origin` - Origin of a key. The default value is kms.

* `scheduled_deletion_date` - Scheduled deletion time (time stamp) of a key.

* `domain_id` - ID of a user domain for the key.

* `expiration_time` - Expiration time.

* `creation_date` - Creation time (time stamp) of a key.

* `is_enabled` - See Argument Reference above.

* `tags` - See Argument Reference above.

* `rotation_number` - Number of key rotations.

## Import

KMS Keys can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_kms_key_v1.key_1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```
