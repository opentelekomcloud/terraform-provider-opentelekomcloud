---
subcategory: "Scalable File Service (SFS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sfs_turbo_share_v1"
sidebar_current: "docs-opentelekomcloud-datasource-sfs-turbo-share-v1"
description: |-
Get details about a SFS Turbo from OpenTelekomCloud
---

Up-to-date reference of API arguments for SFS you can get at
[documentation portal](https://docs.otc.t-systems.com/scalable-file-service/api-ref/sfs_turbo_apis/lifecycle_management/querying_details_about_all_file_systems.html#sfs-02-0053)

# opentelekomcloud_sfs_turbo_share_v1

Use this data source to get details about a Shared File System (SFS) Turbo resource.

## Example Usage

```hcl
data "opentelekomcloud_sfs_turbo_share_v1" "turbo" {
  name = "turbo-share-1"
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) The name of an SFS Turbo share.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the SFS Turbo file system.

* `region` - The region of the SFS Turbo file system.

* `version` - The version ID of the SFS Turbo file system.

* `export_location` - The mount point of the SFS Turbo file system.

* `available_capacity` - The available capacity of the SFS Turbo file system in the unit of GB.

* `region` - The region of SFS Turbo share.

* `size` - Capacity of the share common file system, in GB.

* `share_proto` - The protocol for sharing file systems.

* `share_type` - The file system type.

* `availability_zone` - Tthe availability zone where the file system is located.

* `vpc_id` - The share VPC ID.

* `subnet_id` - Specifies the share network ID of the subnet.

* `security_group_id` - The share security group ID.

* `crypt_key_id` - The ID of a KMS key to encrypt the file system.
