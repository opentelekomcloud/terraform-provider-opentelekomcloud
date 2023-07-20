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

### Basic SAML example

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

### Basic OIDC example

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
  protocol    = "oidc"
  provider_id = opentelekomcloud_identity_provider_v3.provider.id
  mapping_id  = opentelekomcloud_identity_mapping_v3.mapping.id
  access_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "your_client_id"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
    response_type          = "id_token"
    response_mode          = "fragment"
    signing_key = jsonencode(
      {
        keys = [
          {
            alg = "RS256"
            e   = "AQAB"
            kid = "..."
            kty = "RSA"
            n   = "..."
            use = "sig"
          },
        ]
      }
    )
  }
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

* `access_config` - (Optional, List) Specifies the description of the identity provider.
  This field is required only if the protocol is set to *oidc*.

    * `access_type` - (Required) Specifies the access type of the identity provider.
      Available options are:
      + `program`: programmatic access only.
      + `program_console`: programmatic access and management console access.

    * `provider_url` - (Required) Specifies the URL of the identity provider.
      This field corresponds to the iss field in the ID token.

    * `client_id` - (Required) Specifies the ID of a client registered with the OpenID Connect identity provider.

    * `signing_key` - (Required) Public key used to sign the ID token of the OpenID Connect identity provider.
      This field is required only if the protocol is set to *oidc*.

    * `authorization_endpoint` - (Optional) Specifies the authorization endpoint of the OpenID Connect identity
      provider. This field is required only if the access type is set to `program_console`.

    * `scopes` - (Optional) Specifies the scopes of authorization requests. It is an array of one or more scopes.
      Valid values are *openid*, *email*, *profile* and other values defined by you.
      This field is required only if the access type is set to `program_console`.

    -> **NOTE:** 1. *openid* must be specified for this field.
    <br/>2. A maximum of 10 values can be specified, and they must be separated with spaces.
    <br/>Example: openid email host.

    * `response_type` - (Optional) Response type. Valid values is *id_token*, default value is *id_token*.
      This field is required only if the access type is set to `program_console`.

    * `response_mode` - (Optional) Response mode.
      Valid values is *form_post* and *fragment*, default value is *form_post*.
      This field is required only if the access type is set to `program_console`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `links` - Resource links of an identity protocol, including `identity_provider` and `self`.

## Import

Protocols can be imported using the `provider_id/protocol`, e.g.

```shell
terraform import opentelekomcloud_identity_protocol_v3.protocol ACME/saml
```

