---
subcategory: "Scalable File Service"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gpfs_access_rule_v3"
sidebar_current: "docs-opentelekomcloud-resource-gpfs-access-rule-v3"
description: |-
  Manages a GPFS (General Purpose File System) access rule resource within T-Cloud Public (formerly OpenTelekomCloud).
---

# opentelekomcloud_gpfs_access_rule_v3

Manages a GPFS (General Purpose File System) access rule resource within T-Cloud Public (formerly OpenTelekomCloud).

Each rule grants a VPC read/write or read-only access to a file system.
Configuring rules replaces the current file system access-rule configuration.

-> **Note:** Deleting this resource removes **all** access rules of the file system —
the API does not support deleting a single rule.

Up-to-date reference of API arguments for General Purpose File System service you can get at
[documentation portal](https://docs.otc.t-systems.com/scalable-file-service/api-ref/general_purpose_file_system_apis_beta/file_system_access_rules/index.html)

## Example Usage

```hcl
variable "file_system_name" {}
variable "vpc_id" {}

resource "opentelekomcloud_gpfs_access_rule_v3" "test" {
  fs_name = var.file_system_name
  vpc_id  = var.vpc_id
  action  = "FullControl"
}
```

## Argument Reference

The following arguments are supported:

* `fs_name` - (Required, String, ForceNew) Specifies the name of the file system
  the access rule is configured for.
  Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC that is granted
  access to the file system.
  Changing this will create a new resource.

* `action` - (Optional, String, ForceNew) Specifies the type of access granted to the
  VPC. Valid values are as follows:

  + `FullControl` — read/write access (default).
  + `Read` — read-only access.

  Changing this will create a new resource.

* `vpc_source_ip` - (Optional, List, ForceNew) Specifies the IP addresses or
  CIDR ranges in the source VPC. The service currently ignores this field; its value
  is read back from the API.
  Changing this will create a new resource.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the access rule in the format `<fs_name>/<vpc_id>`.

## Import

GPFS access rules can be imported using the file system name and the VPC ID, e.g.

```sh
terraform import opentelekomcloud_gpfs_access_rule_v3.test example-fs-01/vpc-12345678
```
