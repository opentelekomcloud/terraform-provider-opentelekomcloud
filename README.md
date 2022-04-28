<a href="https://terraform.io">
    <img src=".github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

Terraform Open Telekom Cloud Provider
=====================================
[![Documentation](https://img.shields.io/badge/documentation-blue)](https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs)

Quick Start
-----------
> When using the OpenTelekomCloud Provider with Terraform 0.13 and later, the recommended approach is to declare Provider versions in the root module Terraform configuration, using a `required_providers` block as per the following example. For previous versions, please continue to pin the version within the provider block.

1. Add [opentelekomcloud/opentelekomcloud](https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs) to your `required_providers`.
```hcl
# provider.tf
terraform {
   required_providers {
      opentelekomcloud = {
         source = "opentelekomcloud/opentelekomcloud"
         version = ">= 1.23.2"
      }
   }
}
```
2. Run `terraform init -upgrade` to download the provider.
3. Add the provider and supply your `tenant_name` and `domain_name` for minimum configuration.
```hcl
# provider.tf
provider "opentelekomcloud" {
   # OpenTelekomCloud Provider Documentation:
   # https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs
   # domain_name = "..."
   # tenant_name = "..."
   # auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
   # user_name   = "..."
   # password    = "..."
}
```
5. [Authenticate](https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs#authentication) either by providing `user_name` and `password` in the previous file or setting them as environment variables.
```bash
# Linux
OS_USERNAME="<your_username>"
OS_PASSWORD="<your_password"
# Windows
$env:OS_USERNAME="<your_username>"
$env:OS_PASSWORD="<your_password"
```
7. Create your first resource.
```hcl
# main.tf

# Create an Elastic Cloud Server resource
resource "opentelekomcloud_compute_instance_v2" "debian_ecs" {
   name        = "debian_ecs"
   image_name  = "Standard_Debian_11_latest"
   flavor_name = "s3.medium.1"

   key_pair        = "kp_ecs"
   security_groups = ["default"]
   network {
      name = "network_ecs"
   }
}
```

### Full Examples

 - [Have a look here for basic examples](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/tree/devel/examples/basic-examples/modules)
 - [and here for more advanced examples](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/tree/master/examples)

Don't forget to fill in the required variables.

## Developing the Provider

See [Contribution Guide](.github/CONTRIBUTING.md) for the details.

### Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.13+ (but 1.x is recommended)
- [Go](https://golang.org/doc/install) 1.16.x (to build the provider plugin)


### Building The Provider

Clone repository to: `$GOPATH/src/github.com/opentelekomcloud/terraform-provider-opentelekomcloud`

```sh
$ export GO111MODULE=on
$ go get github.com/opentelekomcloud/terraform-provider-opentelekomcloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/opentelekomcloud/terraform-provider-opentelekomcloud
$ make build
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-opentelekomcloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
