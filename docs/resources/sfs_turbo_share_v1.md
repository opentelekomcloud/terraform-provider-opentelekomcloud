---
subcategory: "Scalable File Service (SFS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sfs_turbo_share_v1"
sidebar_current: "docs-opentelekomcloud-resource-sfs-turbo_share-v1"
description: |-
  Manages an SFS Turbo Share resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SFS turbo share you can get at
[documentation portal](https://docs.otc.t-systems.com/scalable-file-service/api-ref/sfs_turbo_apis/lifecycle_management)

# opentelekomcloud_sfs_turbo_share_v1

Provides a Shared File System (SFS) Turbo resource.

## Example Usage

### Basic SFS Turbo share

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "sg_id" {}
variable "az" {}

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name              = "sfs-turbo"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.sg_id
  availability_zone = var.az
}
```

### Enhanced SFS Turbo Performance share

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "sg_id" {}
variable "az" {}

resource "opentelekomcloud_sfs_turbo_share_v1" "sfs-turbo" {
  name              = "sfs-turbo"
  size              = 500
  share_proto       = "NFS"
  share_type        = "PERFORMANCE"
  expand_type       = "bandwidth"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.sg_id
  availability_zone = var.az
}
```

## Argument Reference
The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo resource. If omitted, the
  provider-level region will be used. Changing this creates a new SFS Turbo resource.

* `name` - (Required, String, ForceNew) Specifies the name of an SFS Turbo file system. The value contains 4 to 64
  characters and must start with a letter. Changing this will create a new resource.

* `size` - (Required, Int) Specifies the capacity of a common file system, in GB. The value ranges
  from `500` to `32768`.For an SFS Turbo Enhanced file system, if expand_type is set to bandwidth in the metadata field, the capacity ranges from 10240 to 327680, in GiB. If `expand_type` is set to `hpc` and `hpc_bw` is set to `20M` (20MB/s/TiB), `40M` (40MB/s/TiB:), `125M` (125MB/s/TiB), or `250M` (250MB/s/TiB), the capacity ranges from 3686 to 1048576, in GiB. The capacity must be a **multiple of 1.2 TiB**. The value must be rounded down after being converted to GiB. For example, 3.6TiB->3686GiB, 4.8TiB->4915GiB, 8.4TiB->8601GiB.


* `share_proto` - (Optional, String, ForceNew) Specifies the protocol for sharing file systems. The valid value is `NFS`.
  Changing this will create a new resource.

* `share_type` - (Optional, String, ForceNew) Specifies the file system type. The valid values are `STANDARD` and `PERFORMANCE`.
  Changing this will create a new resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone where the file system is located.
  Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of the subnet. Changing this will create a new resource.

* `security_group_id` - (Required, String, ForceNew) Specifies the security group ID.

* `crypt_key_id` - (Optional, String, ForceNew) Specifies the ID of a KMS key to encrypt the file system.
  Changing this will create a new resource.

* `expand_type` - (Optional, String, ForceNew) Specifies the extension type. This parameter is mandatory when you are creating an SFS Turbo 250 MB/s/TiB, 125 MB/s/TiB, 40 MB/s/TiB, 20 MB/s/TiB, or Enhanced file system. Accepted values: 
    * `bandwidth`: To create a `Standard - Enhanced` or `Performance - Enhanced` file system.
    * `hpc`: To create a 250 MB/s/TiB, 125 MB/s/TiB, 40 MB/s/TiB, or 20 MB/s/TiB file system.

* `hpc_bw` - (Optional, String, ForceNew) Specifies the file system bandwidth. This parameter is mandatory when you are creating an SFS Turbo 250 MB/s/TiB, 125 MB/s/TiB, 40 MB/s/TiB, 20 MB/s/TiB file system. Accepted values: `20M` (for a 20 MB/s/TiB file system), `40M` (for a 40 MB/s/TiB file system), `125M` (for a 125 MB/s/TiB file system), and `250M` (for a 250 MB/s/TiB file system).

-> SFS Turbo will create two private IP addresses and one virtual IP address under the subnet you specified.
To ensure normal use, SFS Turbo will enable the inbound rules for ports `111`, `445`, `2049`, `2051`, `2052`,
and `20048` in the security group you specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the SFS Turbo file system.

* `region` - The region of the SFS Turbo file system.

* `version` - The version ID of the SFS Turbo file system.

* `export_location` - The mount point of the SFS Turbo file system.

* `available_capacity` - The available capacity of the SFS Turbo file system in the unit of GB.

## Timeouts

This resource provides the following timeouts configuration options:
  - `create` - Default is 10 minute.
  - `delete` - Default is 10 minute.

## Import

SFS Turbo can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_sfs_turbo_share_v1 9e3dd316-64g9-0245-8456-71e9387d71ab
```
