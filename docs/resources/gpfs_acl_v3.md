---
subcategory: "General Purpose File System (GPFS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gpfs_acl_v3"
sidebar_current: "docs-opentelekomcloud-resource-gpfs-acl-v3"
description: |-
  Manages the complete GPFS file-system ACL.
---

# opentelekomcloud_gpfs_acl_v3

Manages the complete access control list of a GPFS (General Purpose File System).
Adding, removing, or changing a statement updates the ACL in place.

-> **Note:** Use exactly one ACL resource per file system. Do not manage the same
ACL with `opentelekomcloud_gpfs_access_rule_v3`, another Terraform state, or another
writer. Every apply replaces the complete ACL; omitted grants are revoked.
Destroying this resource removes **all** file-system access rules.

The [ACL API reference](https://docs.otc.t-systems.com/scalable-file-service/api-ref/general_purpose_file_system_apis/file_system_access_rules/index.html)
describes the underlying operations. ACL changes take approximately 30 seconds
to become effective. A VPC endpoint is required in each authorized VPC before
clients can mount the file system.

## Example Usage

```hcl
variable "file_system_name" {}
variable "first_vpc_id" {}
variable "second_vpc_id" {}

resource "opentelekomcloud_gpfs_acl_v3" "example" {
  fs_name = var.file_system_name

  statement {
    vpc_id = var.first_vpc_id
    action = "FullControl"
  }

  statement {
    vpc_id = var.second_vpc_id
    action = "Read"
  }
}
```

## Argument Reference

* `fs_name` - (Required, String, ForceNew) Name of the file system.
* `statement` - (Required, List) Complete desired ACL. Specify one statement per
  VPC, with at least one statement and at most 20. The provider preserves
  configuration order when the API returns statements in a different order.
  To revoke every grant, destroy the ACL resource.

Each `statement` supports:

* `vpc_id` - (Required, String) ID of the authorized VPC.
* `action` - (Optional, String) `FullControl` (read/write, default) or `Read` (read-only).
* `effect` - (Optional, String) Permission effect. Only `Allow` is supported and is the default.
* `sid` - (Optional, String) Statement ID. When omitted, the value returned by the
  service is retained.
* `vpc_source_ip` - (Optional, List of Strings) Source IP addresses or CIDR ranges.
  The API documents this field but **does not currently support it**. Do not rely
  on it to restrict access. When omitted, the service-returned value is retained.

## Attribute Reference

* `id` - File-system name.

## Import

Import the complete existing ACL using the file-system name:

```sh
terraform import opentelekomcloud_gpfs_acl_v3.example example-fs-01
```
