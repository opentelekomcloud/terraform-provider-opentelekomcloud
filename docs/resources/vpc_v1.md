---
subcategory: "Virtual Private Cloud (VPC)"
---

Up-to-date reference of API arguments for VPC service you can get at
`https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/virtual_private_cloud`.

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

### VPC with secondary cidr block

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_sec_cidr" {
  name           = "tf_vpc"
  description    = "description"
  cidr           = "192.168.0.0/16"
  secondary_cidr = "23.9.0.0/16"
  shared         = true

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

* `secondary_cidr` - (Optional) Secondary CIDR block that can be added to VPCs.
  The value cannot contain the following: `100.64.0.0/1`, `214.0.0.0/7`, `198.18.0.0/15`, `169.254.0.0/16`,
  `0.0.0.0/8`, `127.0.0.0/8`, `240.0.0.0/4`, `172.31.0.0/16`, `192.168.0.0/16`.
  Currently, only one secondary CIDR block can be added to each VPC.

* `name` - (Required) The name of the VPC. The name must be unique for a tenant. The value is a string of
  no more than `64` characters and can contain digits, letters, underscores (`_`), and hyphens (`-`).

* `description` - (Optional) A description of the VPC.

* `shared` - (Optional) Specifies whether the shared SNAT should be used or not. Is also
  required for cross-tenant sharing. Shared SNAT only avadilable in eu-de region.

* `tags` - (Optional) The key/value pairs to associate with the VPC.


## Attributes Reference

All above argument parameters can be exported as attribute parameters.

* `status` - The current status of the desired VPC. Can be either `CREATING`,
  `OK`, `DOWN`, `PENDING_UPDATE`, `PENDING_DELETE` or `ERROR`.

## Import

VPCs can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_v1.vpc_v1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
