---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_protocol_v3

Manages identity protocol resource providing binding between identity provider and identity mappings.

-> You _must_ have security admin privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).


## Example Usage

```hcl
resource "opentelekomcloud_identity_provider_v3" "provider" {
  name        = "ACME"
  description = "This is simple identity provider"
  enabled     = true
}

resource "opentelekomcloud_identity_mapping_v3" "mapping" {
  mapping_id = "ACME"
  rules      = file("./rules.json")
}

resource "opentelekomcloud_identity_protocol_v3" "saml" {
  protocol    = "saml"
  provider_id = opentelekomcloud_identity_provider_v3.provider.id
  mapping_id  = opentelekomcloud_identity_mapping_v3.mapping.id
}
```

## Argument Reference

The following arguments are supported:

`protocol` - (Required) ID of a protocol. Changing this creates a new protocol.

`provider_id` - (Required) ID of an identity provider. Changing this creates a new protocol.

`mapping_id` - (Required) ID of an identity mapping.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

`links` - Resource links of an identity protocol, including `identity_provider` and `self`.

## Import

Protocols can be imported using the `provider_id/protocol`, e.g.

```shell
terraform import opentelekomcloud_identity_protocol_v3.protocol ACME/saml
```

