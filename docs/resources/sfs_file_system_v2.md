---
subcategory: "Scalable File Service (SFS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sfs_file_system_v2"
sidebar_current: "docs-opentelekomcloud-resource-sfs-file-system-v2"
description: |-
  Manages an SFS File System resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SFS file system you can get at
[documentation portal](https://docs.otc.t-systems.com/scalable-file-service/api-ref/sfs_capacity-oriented_apis/file_systems)

# opentelekomcloud_sfs_file_system_v2

Provides a Scalable File System resource.

## Example Usage

```hcl
variable "share_name" {}

variable "share_description" {}

resource "opentelekomcloud_sfs_file_system_v2" "share-file" {
  name        = var.share_name
  size        = 50
  description = var.share_description
  share_proto = "NFS"

  tags = {
    muh = "kuh"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 SFS client. If omitted, the
  `region` argument of the provider is used. Changing this creates a new share.

* `size` - (Required) The size (GB) of the shared file system.

* `share_proto` - (Optional) The protocol for sharing file systems. The default value is `NFS`.

* `name` - (Optional) The name of the shared file system.

* `description` - (Optional) Describes the shared file system.

* `is_public` - (Optional) The level of visibility for the shared file system.

* `metadata` - (Optional) Metadata key/value pairs as a dictionary of strings. Changing this will
  create a new resource.

* `availability_zone` - (Optional) The availability zone name. Changing this parameter will create
  a new resource.

* `access_level` - (Optional) The access level of the shared file system. Changing this will create
  a new access rule. Deprecated, please use the `opentelekomcloud_sfs_share_access_rule_v2`
  resource instead.

* `access_type` - (Optional) The type of the share access rule. Changing this will create a new
  access rule. Deprecated, please use the `opentelekomcloud_sfs_share_access_rule_v2` resource instead.

* `access_to` - (Optional) The access that the back end grants or denies. Changing this will
  create new access rule. Deprecated, please use the `opentelekomcloud_sfs_share_access_rule_v2`
  resource instead.

* `tags` - (Optional) Tags key/value pairs to associate with the SFS File System.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the shared file system.

* `status` - The status of the shared file system.

* `share_type` - The storage service type assigned for the shared file system, such as
  high-performance storage (composed of SSDs) and large-capacity storage (composed of SATA disks).

* `volume_type` - The volume type.

* `export_location` - The address for accessing the shared file system.

* `host` - The host name of the shared file system.

* `share_access_id` - The UUID of the share access rule.

* `access_rule_status` - The status of the share access rule.

* `tags` - See Argument Reference above.

## Import

SFS can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_sfs_file_system_v2 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
