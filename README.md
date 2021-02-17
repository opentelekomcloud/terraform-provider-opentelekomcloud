Terraform Open Telekom Cloud Provider
=====================================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.15 (to build the provider plugin)


Building The Provider
---------------------

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

Quick Start
-----------
> When using the OpenTelekomCloud Provider with Terraform 0.13 and later, the recommended approach is to declare Provider versions in the root module Terraform configuration, using a `required_providers` block as per the following example. For previous versions, please continue to pin the version within the provider block.

```hcl
# We strongly recommend using the required_providers block to set the
# OpenTelekomCloud Provider source and version being used
terraform {
  required_providers {
    opentelekomcloud = {
      source = "opentelekomcloud/opentelekomcloud"
      version = ">= 1.22.0"
    }
  }
}

provider "opentelekomcloud" {
  # More information on the authentication methods supported by
  # the OpenTelekomCloud Provider can be found here:
  # https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs

  # user_name   = "..."
  # password    = "..."
  # domain_name = "..."
  # tenant_name = "..."
  # auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
}

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

Full Example
------------
Please see full example at https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/tree/master/examples,
you must fill in the required variables in variables.tf.

Using the provider
------------------
Please see the documentation at [provider usage](docs/index.md).

Developing the Provider
-----------------------

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
