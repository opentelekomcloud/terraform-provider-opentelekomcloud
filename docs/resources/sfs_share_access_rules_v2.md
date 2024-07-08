---
subcategory: "Scalable File Service (SFS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_sfs_share_access_rules_v2"
sidebar_current: "docs-opentelekomcloud-resource-sfs-share-access-rules-v2"
description: |-
Manages an SFS Access Rules resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SFS access rules you can get at
[documentation portal](https://docs.otc.t-systems.com/scalable-file-service/api-ref/sfs_capacity-oriented_apis/file_system_access_rules)

# opentelekomcloud_sfs_share_access_rules_v2

Provides a possibility to manage access rules of Scalable File System resource.

## Example Usage

```hcl
variable "share_name" {}

variable "share_description" {}

resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "sfs_share_vpc_1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "sfs_share_vpc_2"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  name        = var.share_name
  size        = 50
  description = var.share_description
  share_proto = "NFS"
}

resource "opentelekomcloud_sfs_share_access_rules_v2" "sfs_rules" {
  share_id = opentelekomcloud_sfs_file_system_v2.sfs_1.id

  access_rule {
    access_to    = opentelekomcloud_vpc_v1.vpc_1.id
    access_type  = "cert"
    access_level = "rw"
  }

  access_rule {
    access_to    = opentelekomcloud_vpc_v1.vpc_2.id
    access_type  = "cert"
    access_level = "rw"
  }
}
```

## Argument Reference

The following arguments are supported:

* `share_id` - (Required) The UUID of the shared file system.

* `access_rule` - (Required) Specifies the access rules of SFS file share. Structure is documented below.

The `access_rule` block supports:

* `access_level` - (Required) The access level of the shared file system. Possible values are `ro` (read-only)
  and `rw` (read-write). The default value is `rw` (read/write).

* `access_type` - (Optional) The type of the share access rule. The value `cert` indicates
  that the certificate is used to access the storage.

* `access_to` - (Required) The access that the back end grants or denies.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_rule` - See Argument Reference above. The `access_rule` block also contains:

The `access_rule` block supports:

* `share_access_id` - The UUID of the share access rule.

* `access_rule_status` - The status of the share access rule.

## Import

SFS access rules can be imported using the `id` of the file share, e.g.

```shell
terraform import opentelekomcloud_sfs_share_access_rules_v2 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
