---
subcategory: "Elastic Cloud Server (ECS)"
---

# opentelekomcloud_compute_keypair_v2

Manages a V2 keypair resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_compute_keypair_v2" "test-keypair" {
  name       = "my-keypair"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the keypair. Changing this creates a new keypair.

* `public_key` - (Optional) A pre-generated OpenSSH-formatted public key.
  Changing this creates a new keypair.

->
If both `name` and `public_key` duplicate the existing keypair value, the new keypair won't be
managed by the Terraform. Keypair resource will be marked as `shared.`

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `public_key` - See Argument Reference above.

* `private_key` - The information about the private for an SSH key.

* `value_specs` - See Argument Reference above.

* `shared` - Indicates that keypair is shared (global) and not managed by Terraform.

## Import

Keypair can be imported using the `name`, e.g.

```sh
terraform import opentelekomcloud_compute_keypair_v2.my-keypair test-keypair
```

Imported key pairs are considered to be not shared.
