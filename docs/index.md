# Open Telekom Cloud Provider

The Open Telekom Cloud provider is used to interact with the
many resources supported by OpenTelekomCloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the OpenTelekomCloud Provider
provider "opentelekomcloud" {
  user_name   = var.user_name
  password    = var.password
  domain_name = var.domain_name
  tenant_name = var.tenant_name
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}

# Create a web server
resource "opentelekomcloud_compute_instance_v2" "test-server" {
  # ...
}
```

## Authentication

This provider offers 5 means for authentication.

- User name + Password
- AK/SK
- Token
- Federated
- Assume Role
- OpenStack configuration file

### User name + Password

```hcl
provider "opentelekomcloud" {
  user_name   = var.user_name
  password    = var.password
  domain_name = var.domain_name
  tenant_name = var.tenant_name
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

### AKSK

```hcl
provider "opentelekomcloud" {
  access_key  = var.access_key
  secret_key  = var.secret_key
  domain_name = var.domain_name
  tenant_name = var.tenant_name
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

### Temporary AKSK

```hcl
provider "opentelekomcloud" {
  access_key     = var.access_key
  secret_key     = var.secret_key
  security_token = var.security_token
  domain_name    = var.domain_name
  tenant_name    = var.tenant_name
  auth_url       = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

### Token

```hcl
provider "opentelekomcloud" {
  token       = var.token
  domain_name = var.domain_name
  tenant_name = var.tenant_name
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

-> If token, AK/SK and password are set simultaneously, authentication will be done in the following order:
  Token, AK/SK, and Password.

### Federated

```hcl
provider "opentelekomcloud" {
  token          = var.token
  security_token = var.security_token
  access_key     = var.access_key
  secret_key     = var.secret_key
  domain_name    = var.domain_name
  tenant_name    = var.tenant_name
  auth_url       = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

### Assume Role

#### User name + Password

```hcl
provider "opentelekomcloud" {
  agency_name        = var.agency_name
  agency_domain_name = var.agency_domain_name
  delegated_project  = var.delegated_project
  user_name          = var.user_name
  password           = var.password
  domain_name        = var.domain_name
  auth_url           = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

#### User ID + Password + TOTP
```hcl
provider "opentelekomcloud" {
  agency_name        = var.agency_name
  agency_domain_name = var.agency_domain_name
  delegated_project  = var.delegated_project
  user_id            = var.user_id
  password           = var.password
  domain_name        = var.domain_name
  auth_url           = "https://iam.eu-de.otc.t-systems.com/v3"
  passcode           = var.passcode
}
```

#### AK/SK

```hcl
provider "opentelekomcloud" {
  agency_name        = var.agency_name
  agency_domain_name = var.agency_domain_name
  delegated_project  = var.delegated_project
  access_key         = var.access_key
  secret_key         = var.secret_key
  domain_name        = var.domain_name
  auth_url           = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

### Temporary AKSK

```hcl
provider "opentelekomcloud" {
  agency_name        = var.agency_name
  agency_domain_name = var.agency_domain_name
  delegated_project  = var.delegated_project
  access_key         = var.access_key
  secret_key         = var.secret_key
  security_token     = var.security_token
  domain_name        = var.domain_name
  auth_url           = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

#### Token

```hcl
provider "opentelekomcloud" {
  agency_name        = var.agency_name
  agency_domain_name = var.agency_domain_name
  delegated_project  = var.delegated_project
  token              = var.token
  auth_url           = "https://iam.eu-de.otc.t-systems.com/v3"
}
```
`token` specified is not the normal token, but must have the authority of `Agent Operator`.

### OpenStack configuration file

```hcl
provider "opentelekomcloud" {
  cloud = var.cloud_name
}
```

`cloud` should be the name of cloud in `clouds.yaml`

See [OpenStack configuration documentation](https://docs.openstack.org/python-openstackclient/latest/configuration/index.html) for details.


## Configuration Reference

The following arguments are supported:

* `access_key` - (Optional) The access key of the OpenTelekomCloud cloud to use.
  If omitted, the `OS_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the OpenTelekomCloud cloud to use.
  If omitted, the `OS_SECRET_KEY` environment variable is used.

* `auth_url` - (Optional; required if `cloud` is not specified) The Identity
  authentication URL. If omitted, the `OS_AUTH_URL` environment variable is used.

* `cloud` - (Optional; required if `auth_url` is not specified) An entry in a
  `clouds.yaml` file. See the OpenStack `os-client-config`
  [documentation](https://docs.openstack.org/os-client-config/latest/user/configuration.html)
  for more information about `clouds.yaml` files. If omitted, the `OS_CLOUD`
  environment variable is used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `OS_USERNAME` environment variable is used.

* `user_id` - (Optional) The ID of the user to login with. Required when TOTP is used (`passcode` is not empty).
  If `user_id` is set, `user_name` is ignored.

* `tenant_name` - (Optional) The Name of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_NAME` or
  `OS_PROJECT_NAME` environment variable are used.

* `region` - (Optional) The name of the region to be used. Required for some resources
  (e.g. `s3_bucket`) in case no tenant name provided and no region is defined in the
  resource. If omitted, the `OS_REGION` or `OS_REGION_NAME` environment variables are used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `token` - (Optional; Required if not using `user_name` and `password`)
  A token is an expiring, temporary means of access issued via the Keystone
  service. By specifying a token, you do not have to specify a username/password
  combination, since the token was already created by a username/password out of
  band of Terraform. If omitted, the `OS_AUTH_TOKEN` or `OS_TOKEN` environment
  variable is used.

* `security_token` - (Optional) Security token required to authenticate with temporary AK/SK.

* `passcode` - (Optional) One-time password provided by your authentication app.

->
  Please note that MFA requires `user_id` to be used. Setting `user_name` won't work.

* `domain_name` - (Optional) The Name of the Domain to scope to (Identity v3).
  If omitted, the following environment variables are checked (in this order):
  `OS_USER_DOMAIN_NAME`, `OS_PROJECT_DOMAIN_NAME`, `OS_DOMAIN_NAME`,
  `DEFAULT_DOMAIN`.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `cacert_file` - (Optional) Specify a custom CA certificate when communicating
  over SSL. You can specify either a path to the file or the contents of the
  certificate. If omitted, the `OS_CACERT` environment variable is used.

* `cert` - (Optional) Specify client certificate file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the certificate. If omitted the `OS_CERT` environment variable is used.

* `key` - (Optional) Specify client private key file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the key. If omitted the `OS_KEY` environment variable is used.

* `endpoint_type` - (Optional) Specify which type of endpoint to use from the
  service catalog. It can be set using the `OS_ENDPOINT_TYPE` environment
  variable. If not set, public endpoints is used.

* `swauth` - (Optional) Set to `true` to authenticate against Swauth, a
  Swift-native authentication system. If omitted, the `OS_SWAUTH` environment
  variable is used. You must also set `username` to the Swauth/Swift username
  such as `username:project`. Set the `password` to the Swauth/Swift key.
  Finally, set `auth_url` as the location of the Swift service.

-> This will only work when used with the OpenTelekomCloud Object Storage resources.

* `agency_name` - (Optional) if authorized by assume role, it must be set. The
  name of agency.

* `agency_domain_name` - (Optional) if authorized by assume role, it must be set.
  The name of domain who created the agency (Identity v3).

* `delegated_project` - (Optional) The name of delegated project (Identity v3).

* `max_retries` - (Optional) Maximum number of retries of HTTP requests failed
  due to connection issues.

* `max_backoff_retries` - (Optional) Maximum number of retries of HTTP requests failed
  due to reaching the rate limit. It can be set using the `OS_MAX_BACKOFF_RETRIES` environment
  variable. If not set, default value is used.
  Default: `5`

* `backoff_retry_timeout` - (Optional) Timeout in seconds for backoff retry due to reaching the rate limit.
  It can be set using the `OS_BACKOFF_RETRY_TIMEOUT` environment
  variable. If not set, default value is used.
  Default: `60` seconds.

## Additional Logging

This provider has the ability to log all HTTP requests and responses between
Terraform and the OpenTelekomCloud cloud which is useful for troubleshooting and
debugging.

To enable these logs, set the `OS_DEBUG` environment variable to `1` along
with the usual `TF_LOG=DEBUG` environment variable:

```shell
$ OS_DEBUG=1 TF_LOG=DEBUG terraform apply
```

If you submit these logs with a bug report, please ensure any sensitive
information has been scrubbed first!

## Creating an issue

[Issues](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues)
can be used to keep track of bugs, enhancements, or other requests.
See the github help [here](https://help.github.com/articles/creating-an-issue/)

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OS_IMAGE_ID` or `OS_IMAGE_NAME` - a UUID or name of an existing image in Glance.

* `OS_FLAVOR_ID` or `OS_FLAVOR_NAME` - an ID or name of an existing flavor.

* `OS_POOL_NAME` - The name of a Floating IP pool.

* `OS_NETWORK_ID` - The UUID of a network in your test environment.

* `OS_EXTGW_ID` - The UUID of the external gateway.

You should be able to use any OpenTelekomCloud environment to develop on as long as the
above environment variables are set.

<div style="visibility:hidden">

```{toctree}
:maxdepth: 2
:hidden:
data-sources/index
resources/index
```

</div>
