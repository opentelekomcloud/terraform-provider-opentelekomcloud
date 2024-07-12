---
subcategory: "Scalable File Service (SFS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sfs_file_system_v2"
sidebar_current: "docs-opentelekomcloud-datasource-sfs-file-system-v2"
description: |-
  Get details about a Scalable File Service from OpenTelekomCloud
---

Up-to-date reference of API arguments for SFS you can get at
[documentation portal](https://docs.otc.t-systems.com/scalable-file-service/api-ref/sfs_capacity-oriented_apis/file_systems/querying_all_shared_file_systems.html#sfs-02-0022)

# opentelekomcloud_sfs_file_system_v2

Use this data source to get details about a Scalable File Service.

## Example Usage

```hcl
variable "share_name" {}
variable "share_id" {}

data "opentelekomcloud_sfs_file_system_v2" "shared_file" {
  name = var.share_name
  id   = var.share_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the shared file system.

* `id` - (Optional) The UUID of the shared file system.

* `status` - (Optional) The status of the shared file system.


## Attributes Reference

The following attributes are exported:

* `availability_zone` - The availability zone name.

* `size` - 	The size (GB) of the shared file system.

* `share_type` - The storage service type for the shared file system, such as high-performance storage (composed of SSDs) or large-capacity storage (composed of SATA disks).

* `status` - The status of the shared file system.

* `host` - The host name of the shared file system.

* `is_public` - The level of visibility for the shared file system.

* `share_proto` - The protocol for sharing file systems.

* `volume_type` - The volume type.

* `metadata` - Metadata key and value pairs as a dictionary of strings.

* `export_location` - The path for accessing the shared file system.

* `access_level` - The level of the access rule.

* `access_rules_status` - The status of the share access rule.

* `access_type` - The type of the share access rule.

* `access_to` - The access that the back end grants or denies.

* `share_access_id` - The UUID of the share access rule.

* `mount_id` - The UUID of the mount location of the shared file system.

* `share_instance_id` - The access that the back end grants or denies.

* `preferred` - Identifies which mount locations are most efficient and are used preferentially when multiple mount locations exist.
