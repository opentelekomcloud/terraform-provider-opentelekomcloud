---
subcategory: "Scalable File Service (SFS)"
---

# opentelekomcloud_sfs_share_access_rule_v2

Provides a Scalable File System resource.

## Example Usage

```hcl
variable "share_name" {}

variable "share_description" {}

resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name   = "sfs_share_vpc_1"
  cidr   = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name   = "sfs_share_vpc_2"
  cidr   = "192.168.0.0/16"
}

resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  name         = var.share_name
  size         = 50
  description  = var.share_description
  share_proto  = "NFS"
}

resource "opentelekomcloud_sfs_share_access_rule_v2" "sfs_rules" {
  share_id = opentelekomcloud_sfs_file_system_v2.sfs_1.id

  access_rules {
    access_to    = opentelekomcloud_vpc_v1.vpc_1.id
    access_type  = "cert"
    access_level = "rw"
  }
  access_rules {
    access_to    = opentelekomcloud_vpc_v1.vpc_2.id
    access_type  = "cert"
    access_level = "rw"
  }
}
```

## Argument Reference

The following arguments are supported:

* `share_id` - (Required) The UUID of the shared file system.

* `access_rules` - (Required) Specifies the access rules of SFS file share. Structure is documented below.

The `access_rules` block supports:

* `access_level` - (Optional) The access level of the shared file system.

* `access_type` - (Optional) The type of the share access rule.

* `access_to` - (Optional) The access that the back end grants or denies.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the shared file system.

* `share_id` - The UUID of the shared file system.

* `access_rules` - See Argument Reference above. The `access_rules` block also contains:

* `share_access_id` - The UUID of the share access rule.

* `access_rule_status` - The status of the share access rule.

## Import

SFS access rules can be imported using the `id` of the file share, e.g.

```shell
terraform import opentelekomcloud_sfs_share_access_rule_v2 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
