---
subcategory: "Identity and Access Management (IAM)"
---

Up-to-date reference of API arguments for IAM provider you can get at
`https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/federated_identity_authentication_management/identity_provider`.

# opentelekomcloud_identity_provider

-> You _must_ have security admin privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).


## Example Usage

### Create a SAML protocol provider

```hcl
resource "opentelekomcloud_identity_provider" "provider_1" {
  name     = "example_com_provider_saml"
  protocol = "saml"
}
```

### Create a OpenID Connect protocol provider

```hcl
resource "opentelekomcloud_identity_provider" "provider_2" {
  name     = "example_com_provider_oidc"
  protocol = "oidc"

  access_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "your_client_id"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the identity provider to be registered.
  The maximum length is 64 characters. Only letters, digits, underscores (_), and hyphens (-) are allowed.
  The name is unique, it is recommended to include domain name information.
  Changing this creates a new resource.

* `protocol` - (Required) Specifies the protocol of the identity provider.
  Valid values are *saml* and *oidc*.

* `status` - (Optional) Enabled status for the identity provider. Default: `true`.

* `description` - (Optional) Specifies the description of the identity provider.

* `metadata` - (Optional) Specifies the metadata of the IDP(Identity Provider) server.
  This field is used to import a metadata file to IAM to implement federated identity authentication.
  This field is required only if the protocol is set to *saml*.
  The maximum length is 30,000 characters and it stores in the state with SHA1 algorithm.

-> **NOTE:**
The metadata file specifies API addresses and certificate information in compliance with the SAML 2.0 standard.
It is usually stored in a file. In the TF script, you can import the metafile through the `file` function,
for example:
<br/>`metadata = file("/usr/local/data/files/metadata.txt")`

* `access_config` - (Optional, List) Specifies the description of the identity provider.
  This field is required only if the protocol is set to *oidc*.

The `access_config` block supports:

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

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `login_link` - The login link of the identity provider.

* `conversion_rules` - The identity conversion rules of the identity provider.
  The structure is documented below.

The `conversion_rules` block supports:

* `local` - The federated user information on the cloud platform.

* `remote` - The description of the identity provider.

The `local` block supports:

* `username` - The name of a federated user on the cloud platform.

* `group` - The user group to which the federated user belongs on the cloud platform.

The `remote` block supports:

* `attribute` - The attribute in the IDP assertion.

* `condition` - The condition of conversion rule.

* `value` - The rule is matched only if the specified strings appear in the attribute type.

## Import

Identity provider can be imported using the `name`, e.g.

```
$ terraform import opentelekomcloud_identity_provider.provider_1 example_provider_saml
```
