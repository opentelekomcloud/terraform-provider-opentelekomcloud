---
subcategory: "General Purpose File System (GPFS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gpfs_file_system_v3"
sidebar_current: "docs-opentelekomcloud-resource-gpfs-file-system-v3"
description: |-
  Manages a GPFS (General Purpose File System) resource within T-Cloud Public (formerly OpenTelekomCloud).
---

# opentelekomcloud_gpfs_file_system_v3

Manages a GPFS (General Purpose File System) resource within T-Cloud Public (formerly OpenTelekomCloud).

Up-to-date reference of API arguments for General Purpose File System service you can get at
[documentation portal](https://docs.otc.t-systems.com/scalable-file-service/api-ref/general_purpose_file_system_apis_beta/file_systems/index.html)

## Example Usage

```hcl
variable "file_system_name" {}

resource "opentelekomcloud_gpfs_file_system_v3" "test" {
  name       = var.file_system_name
  redundancy = "3az"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the file system.

  The name must be 3 to 63 characters in length, consist only of lowercase letters,
  digits and hyphens, and start and end with a lowercase letter or a digit.
  The file system name must be unique in the region.
  Changing this will create a new resource.

* `redundancy` - (Optional, String, ForceNew) Specifies the AZ redundancy of the file
  system. Only `3az` (multi-AZ) is supported. Defaults to `3az`.
  Changing this will create a new resource.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the file system (the file system name).

* `creation_date` - The creation time of the file system (in RFC3339 format).

## Import

GPFS file systems can be imported using the file system name, e.g.

```sh
terraform import opentelekomcloud_gpfs_file_system_v3.test example-fs-01
```
