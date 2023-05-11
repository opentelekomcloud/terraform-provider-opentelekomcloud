---
subcategory: "Key Management Service (KMS)"
---

Up-to-date reference of API arguments for KMS you can get at
`https://docs.otc.t-systems.com/key-management-service/api-ref/apis`.

# opentelekomcloud_kms_grant_v1

Manages a V1 KMS grant resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_kms_grant_v1" "grant_1" {
  key_id            = var.kms_id
  name              = "my_grant"
  grantee_principal = var.user_id
  operations        = ["describe-key", "create-datakey", "encrypt-datakey"]
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Required) Indicates the ID of the KMS. Changing this creates new grant.

* `grantee_principal` - (Required) Indicates the ID of the authorized user.
  Changing this creates new grant.

* `operations` - (Required) Permissions that can be granted.
  The valid values are: `create-datakey`, `create-datakey-without-plaintext`,
  `encrypt-datakey`, `decrypt-datakey`, `describe-key`, `create-grant`, `retire-grant`.
  Changing this creates new grant.

* `name` - (Optional) Name of a grant which can be 1 to 255 characters in length
  and matches the regular expression `^[a-zA-Z0-9:/_-]{1,255}$`.
  Changing this creates new grant.

* `retiring_principal` - (Optional) Indicates the ID of the retiring user.
  Changing this creates new grant.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `issuing_principal` - Indicates the ID of the user who created the grant.

* `creation_date` - Creation time. The value is a timestamp expressed in the number of
  seconds since 00:00:00 UTC on January 1, 1970.


## Import

KMS Grants can be imported using the `key_id/grant_id`, e.g.

```shell
terraform import opentelekomcloud_kms_grant_v1.grant_1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d/7056d636-ac60-4663-8a6c-82d3c32c1c64
```
