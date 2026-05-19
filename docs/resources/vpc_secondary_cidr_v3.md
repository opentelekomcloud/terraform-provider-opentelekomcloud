---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_secondary_cidr_v3"
sidebar_current: "docs-opentelekomcloud-resource-vpc-secondary-cidr-v3"
description: |-
  Manages secondary CIDR blocks for a VPC within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC secondary CIDR block management you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/virtual_private_cloud)

# opentelekomcloud_vpc_secondary_cidr_v3

Manages up to 5 secondary CIDR blocks attached to an existing VPC via the VPC v3 API.

~> **Warning:** Each VPC has a single shared `extend_cidrs` list in the cloud.
Define **at most one** `opentelekomcloud_vpc_secondary_cidr_v3` resource per
VPC, and do **not** combine it with the `secondary_cidr` attribute on
[`opentelekomcloud_vpc_v1`](vpc_v1.md) against the same VPC which will be soon
deprecated. Multiple resources targeting the same `vpc_id` will fight on every
plan, because each resource reads back the full cloud-side list and tries to
reconcile it to its own subset.

## Example Usage

```hcl
resource "opentelekomcloud_vpc_v1" "vpc" {
  name = "tf_vpc"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_secondary_cidr_v3" "secondary" {
  vpc_id = opentelekomcloud_vpc_v1.vpc.id
  cidrs = [
    "23.9.0.0/16",
    "23.10.0.0/16",
    "23.11.0.0/16",
  ]
}
```

## Argument Reference

The following arguments are supported:

- `vpc_id` - (Required, ForceNew) The ID of the VPC to which the secondary CIDR
  blocks are attached. Changing this forces a new resource to be created.

- `cidrs` - (Required) Set of secondary CIDR blocks to attach to the VPC.
  At least 1 and at most 5 entries are supported. Each entry must be a valid
  IPv4 CIDR. The cloud rejects reserved ranges (for example
  `100.64.0.0/10`, `214.0.0.0/7`, `198.18.0.0/15`, `169.254.0.0/16`,
  `0.0.0.0/8`, `127.0.0.0/8`, `240.0.0.0/4`, `172.31.0.0/16`,
  `192.168.0.0/16`); these checks are enforced server-side.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the resource. Equal to the `vpc_id` because each VPC has a
  single set of secondary CIDR blocks.

- `region` - The region in which the secondary CIDRs are managed. Inherited
  from the provider configuration.

## Import

VPC secondary CIDR resources can be imported using the parent VPC `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_secondary_cidr_v3.secondary 7117d38e-4c8f-4624-a505-bd96b97d024c
```
