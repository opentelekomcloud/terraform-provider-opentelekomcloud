---
subcategory: "Cloud Container Engine (CCE)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_node_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cce-node-v3"
description: |-
Get a list of node ids for a CCE cluster from OpenTelekomCloud
---

Up-to-date reference of API arguments for CCE nodes you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-engine/api-ref/apis/cluster_management/listing_all_nodes_in_a_cluster.html)

# opentelekomcloud_cce_node_v3

Use this data source to get the specified node in a cluster from OpenTelekomCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_id" {}

data "opentelekomcloud_cce_node_v3" "node" {
  cluster_id = var.cluster_id
  node_id    = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of container cluster.

* `name` - (Optional) Name of the node.

* `node_id` - (Optional) The id of the node.

* `status` - (Optional) The state of the node.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `flavor_id` - The flavor id to be used.

* `availability_zone` - Available partitions where the node is located.

* `key_pair` - Key pair name when logging in to select the key pair mode.

* `billing_mode` - Node's billing mode: The value is `0` (on demand).

* `charge_mode` - Bandwidth billing type.

* `bandwidth_size` - Bandwidth (Mbit/s), in the range of `[1, 2000]`.

* `disk_size` - Root volume disk size in GB.

* `volume_type` - Root volume disk type.

* `eip_ids` - List of existing elastic IP IDs.

* `server_id` - The node's virtual machine ID in ECS.

* `public_ip` - Elastic IP parameters of the node.

* `private_ip` - Private IP of the node

* `ip_type` - Elastic IP address type.

* `share_type` - The bandwidth sharing type.

* `data_volumes` - Represents the data disks configuration.
  * `size` - Disk size in GB.
  * `volumetype` - Disk type.
  * `extend_param` - Disk expansion parameters.
  * `kms_id` - The Encryption KMS ID of the data volume.

* `runtime` - The runtime of the node.
