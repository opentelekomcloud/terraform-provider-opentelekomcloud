Terraform Open Telekom Cloud Provider
=====================================
[![Documentation](https://camo.githubusercontent.com/3f4485d8adcf4893453f5aacd375ff6e32d9c1b79648335fec37a8a997aa9355/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f646f63756d656e746174696f6e2d626c7565)](https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs) | [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
| [Mailing list](http://groups.google.com/group/terraform-tool)

[![alt text](https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg)](https://www.terraform.io/)

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
     resource "opentelekomcloud_compute_instance_v2" "test-server" {
       name        = "test-server"
       image_name  = "Standard_CentOS_8_latest"
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

Developing the Provider
-----------------------

### Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.15 (to build the provider plugin)


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

If you wish to work on the provider, you'll first need [Go](https://golang.org) installed on your machine (version 1.15+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```sh
make build
...
$GOPATH/bin/terraform-provider-opentelekomcloud
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
