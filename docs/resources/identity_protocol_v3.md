---
subcategory: "Identity and Access Management (IAM)"
---

Up-to-date reference of API arguments for IAM protocol you can get at
`https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/federated_identity_authentication_management/protocol`.

# opentelekomcloud_identity_protocol_v3

Manages identity protocol resource providing binding between identity provider and identity mappings.

-> You _must_ have security admin privileges in your OpenTelekomCloud cloud to use this resource. Please refer
to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).

## Example Usage

### Basic usage

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

### Import SAML metadata file

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


  metadata {
    domain_id = var.domain_id
    metadata  = file("saml-metadata.xml")
  }
}
```

## Argument Reference

The following arguments are supported:

* `protocol` - (Required) ID of a protocol. Changing this creates a new protocol.

* `provider_id` - (Required) ID of an identity provider. Changing this creates a new protocol.

* `mapping_id` - (Required) ID of an identity mapping.

* `metadata` - (Optional) Metadata file configuration.

    * `xaccount_type` - (Optional) Source of a domain. Blank by the default.

    * `metadata` - (Required) Content of the metadata file on the IdP server.

    * `domain_id` - (Required) ID of the domain that a user belongs to.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `links` - Resource links of an identity protocol, including `identity_provider` and `self`.

## Import

Protocols can be imported using the `provider_id/protocol`, e.g.

```shell
terraform import opentelekomcloud_identity_protocol_v3.protocol ACME/saml
```

