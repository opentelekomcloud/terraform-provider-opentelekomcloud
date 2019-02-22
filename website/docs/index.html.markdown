---
layout: "opentelekomcloud"
page_title: "Provider: OpenTelekomCloud"
sidebar_current: "docs-opentelekomcloud-index"
description: |-
  The OpenTelekomCloud provider is used to interact with the many resources supported by OpenTelekomCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# OpenTelekomCloud Provider

The OpenTelekomCloud provider is used to interact with the
many resources supported by OpenTelekomCloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the OpenTelekomCloud Provider
provider "opentelekomcloud" {
  user_name   = "${var.user_name}"
  password    = "${var.password}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}

# Create a web server
resource "opentelekomcloud_compute_instance_v2" "test-server" {
  # ...
}
```

## Authentication

This provider offers 3 means for authentication.

- User name + Password
- AKSK
- Token

### User name + Password

```hcl
provider "opentelekomcloud" {
  user_name   = "${var.user_name}"
  password    = "${var.password}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

### AKSK

```hcl
provider "opentelekomcloud" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}
```

### Token

```hcl
provider "opentelekomcloud" {
  token       = "${var.token}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}
```
Note: if token, aksk and password are set simultaneously, then it will authenticate in the order of Token, AKSK and Password.

## Configuration Reference

The following arguments are supported:

* `access_key` - (Optional) The access key of the OpenTelekomCloud cloud to use.
  If omitted, the `OS_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the OpenTelekomCloud cloud to use.
  If omitted, the `OS_SECRET_KEY` environment variable is used.

* `auth_url` - (Required) The Identity authentication URL. If omitted, the
  `OS_AUTH_URL` environment variable is used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `OS_USERNAME` environment variable is used.

* `tenant_name` - (Required) The Name of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_NAME` or
  `OS_PROJECT_NAME` environment variable are used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `token` - (Optional; Required if not using `user_name` and `password`)
  A token is an expiring, temporary means of access issued via the Keystone
  service. By specifying a token, you do not have to specify a username/password
  combination, since the token was already created by a username/password out of
  band of Terraform. If omitted, the `OS_AUTH_TOKEN` environment variable is used.

* `domain_name` - (Required) The Name of the Domain to scope to (Identity v3).
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
  service catalog. It can be set using the OS_ENDPOINT_TYPE environment
  variable. If not set, public endpoints is used.

* `swauth` - (Optional) Set to `true` to authenticate against Swauth, a
  Swift-native authentication system. If omitted, the `OS_SWAUTH` environment
  variable is used. You must also set `username` to the Swauth/Swift username
  such as `username:project`. Set the `password` to the Swauth/Swift key.
  Finally, set `auth_url` as the location of the Swift service. Note that this
  will only work when used with the OpenTelekomCloud Object Storage resources.

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

[Issues](https://github.com/terraform-providers/terraform-provider-opentelekomcloud/issues)
can be used to keep track of bugs, enhancements, or other requests.
See the github help [here](https://help.github.com/articles/creating-an-issue/)

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OS_IMAGE_ID` or `OS_IMAGE_NAME` - a UUID or name of an existing image in
    Glance.

* `OS_FLAVOR_ID` or `OS_FLAVOR_NAME` - an ID or name of an existing flavor.

* `OS_POOL_NAME` - The name of a Floating IP pool.

* `OS_NETWORK_ID` - The UUID of a network in your test environment.

* `OS_EXTGW_ID` - The UUID of the external gateway.

You should be able to use any OpenTelekomCloud environment to develop on as long as the
above environment variables are set.

Most of Terraform's OpenTelekomCloud support is done in a standardized Packstack
all-in-one environment. You can find the scripts to build this environment
[here](https://github.com/jtopjian/terraform-devstack/tree/master/packstack-standard).
The included `main.tf` file will need to be modified for your specific
environment. Once it's up and running, you will have access to a standard,
up-to-date OpenTelekomCloud environment with the latest OpenTelekomCloud services.

If you require access to deprecated services, such as Keystone v2 and
LBaaS v1, you can use the "legacy" environment
[here](https://github.com/jtopjian/terraform-devstack/tree/master/packstack-legacy).
