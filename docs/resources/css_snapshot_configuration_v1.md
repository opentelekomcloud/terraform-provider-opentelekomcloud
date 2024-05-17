---
subcategory: "Cloud Search Service (CSS)"
---

Up-to-date reference of API arguments for CSS snapshot you can get at
`https://docs.otc.t-systems.com/cloud-search-service/api-ref/snapshot_management_apis`.

# opentelekomcloud_css_snapshot_configuration_v1

Manages a CSS configuration of automatic snapshot creation.

## Example Usage

```hcl
data "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = var.security_group
}

resource "opentelekomcloud_css_cluster_v1" "cluster" {
  name            = "terraform_test_cluster"
  expect_node_num = 1
  node_config {
    flavor = "css.medium.8"
    network_info {
      security_group_id = data.opentelekomcloud_networking_secgroup_v2.secgroup.id
      network_id        = var.network_id
      vpc_id            = var.vpc_id
    }
    volume {
      volume_type = "COMMON"
      size        = 40
    }

    availability_zone = var.availability_zone
  }
}

resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket        = "tf-snap-testing"
  force_destroy = true
}

resource "opentelekomcloud_css_snapshot_configuration_v1" "config" {
  cluster_id = opentelekomcloud_css_cluster_v1.cluster.id
  configuration {
    bucket    = opentelekomcloud_obs_bucket.bucket.bucket
    agency    = "css_obs_agency"
    base_path = "css/snapshot"
  }
  creation_policy {
    prefix      = "snapshot"
    period      = "00:00 GMT+03:00"
    keepday     = 2
    enable      = true
    delete_auto = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) ID of the CSS cluster.

* `automatic` - (Optional) Use automatic configuration for CCS cluster screenshots.
  Mutually exclusive with `configuration`/`creation_policy`.

* `configuration` - (Optional) The basic configurations of a cluster snapshot. Structure is documented below.
  Mutually exclusive with `automatic`.

* `creation_policy` - (Optional) Parameters related to automatic snapshot creation. Structure is documented below.
  Mutually exclusive with `automatic`.

The `configuration` block supports:

* `bucket` - (Required) The bucket which will be used for storing snapshots.

* `agency` - (Required) The agency used by CSS to access OBS.

* `base_path` - (Required) Storage path of the snapshot in the OBS bucket.

* `kms_id` - (Options) Key ID used for snapshot encryption.

~>
  If the key used for encryption is in the Pending deletion or disable state,
  you cannot perform backup and restoration operations on the cluster.
  Specifically, new snapshots cannot be created for the cluster, and existing snapshots cannot be used for restoration.
  In this case, switch to the KMS management console and change the state of the target key to enable so that backup
  and restore operations are allowed on the cluster. For more details
  see https://docs.otc.t-systems.com/api/css/css_03_0030.html

The `creation_policy` block supports:

* `prefix` - (Required) Prefix of the snapshot name that is automatically created.

* `period` - (Required) Time when a snapshot is created every day. Snapshots can only be created on the hour.
  The time format is the time followed by the time zone, specifically, `HH:mm z`.
  In the format, `HH:mm` refers to the hour time and `z` refers to the time zone, for example,
  `00:00 GMT+08:00` and `01:00 GMT+08:00`.

* `keepday` - (Required) Number of days that a snapshot can be retained. The value ranges from `1` to `90`.
  The system automatically deletes snapshots that have been retained for the allowed maximum duration on the half hour.

* `enable` - (Required) Value `true` indicates that the automatic snapshot creation policy is enabled,
  and value `false` indicates that the automatic snapshot creation policy is disabled.

* `delete_auto` - (Optional) Whether to delete all automatically created snapshots when the automatic
  snapshot creation policy is disabled. The default value is `false`, indicating that snapshots that have been
  automatically created are not deleted when the automatic snapshot creation function is disabled.
  If this parameter is set to `true`, all automatically created snapshots are deleted when the automatic snapshot
  creation policy is disabled.
