---
subcategory: "Virtual Private Cloud (VPC)"
---

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
