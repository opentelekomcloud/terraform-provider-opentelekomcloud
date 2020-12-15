---
subcategory: "Bare Metal Server (BMS)"
---

# opentelekomcloud_compute_bms_keypairs_v2

Use this data source to get details about SSH key pairs of BMSs from OpenTelekomCloud.

## Example Usage

```hcl
variable "keypair_name" { }

data "opentelekomcloud_compute_bms_keypairs_v2" "query_bms_keypair" {
  name = var.keypair_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `name` - (Required) It is the key pair name.

## Attributes Reference

All of the argument attributes are also exported as result attributes.

* `public_key` - It gives the information about the public key in the key pair.

* `fingerprint` - It is the fingerprint information about the key pair.
