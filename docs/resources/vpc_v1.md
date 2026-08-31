---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_v1"
sidebar_current: "docs-opentelekomcloud-resource-vpc-v1"
description: |-
  Manages a VPC resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC service you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/vpc_apis_v1_v2/virtual_private_cloud/index.html)

# opentelekomcloud_vpc_v1

Manages a VPC v1 resource within OpenTelekomCloud.

## Example Usage

### Basic Usage

```hcl
variable "vpc_name" {
  default = "opentelekomcloud_vpc"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_v1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}
```

### VPC with tags

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_with_tags" {
  name = var.vpc_name
  cidr = var.vpc_cidr

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The range of available subnets in the VPC. The value ranges from
  `10.0.0.0/8` to `10.255.255.0/24`, `172.16.0.0/12` to `172.31.255.0/24`,
  or `192.168.0.0/16` to `192.168.255.0/24`.

* `secondary_cidr` - **DEPRECATED** (Optional) Secondary CIDR block that can be added to VPCs.
  The value cannot contain the following: `100.64.0.0/1`, `214.0.0.0/7`, `198.18.0.0/15`, `169.254.0.0/16`,
  `0.0.0.0/8`, `127.0.0.0/8`, `240.0.0.0/4`, `172.31.0.0/16`, `192.168.0.0/16`.
  Please use resource `opentelekomcloud_vpc_secondary_cidr_v3` instead.

* `name` - (Required) The name of the VPC. The name must be unique for a tenant. The value is a string of
  no more than `64` characters and can contain digits, letters, underscores (`_`), hyphens (`-`), and periods (`.`).

* `description` - (Optional) A description of the VPC.

* `enterprise_project_id` - (Optional) The enterprise project ID associated with the VPC.
  Changing this creates a new VPC. If omitted, the provider-level enterprise project ID is used,
  or `0` for the default enterprise project.

* `shared` - (Optional) Specifies whether the shared SNAT should be used or not. Is also
  required for cross-tenant sharing. Shared SNAT only avadilable in eu-de region.
  Deprecated, VPC Shared SNAT End of Life from `01.03.2024`, please do not use.

* `tags` - (Optional) The key/value pairs to associate with the VPC.


## Attributes Reference

All above argument parameters can be exported as attribute parameters.

* `status` - The current status of the desired VPC. Can be either `CREATING`,
  `OK`, `DOWN`, `PENDING_UPDATE`, `PENDING_DELETE` or `ERROR`.

* `routes` - The VPC routes. Each entry contains `destination` and `nexthop`.
  Manage custom routes with the dedicated VPC route resources.

* `tenant_id` - The project ID that owns the VPC.

* `created_at` - The UTC creation timestamp.

* `updated_at` - The UTC last-update timestamp.

## Import

VPCs can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_v1.vpc_v1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
