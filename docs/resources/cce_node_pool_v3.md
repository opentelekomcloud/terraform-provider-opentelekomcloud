---
subcategory: "Cloud Container Engine (CCE)"
---

# opentelekomcloud_cce_node_pool_v3

Provides a node pool resource management of a container cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "ssh_key" {}
variable "availability_zone" {}

resource "opentelekomcloud_cce_node_pool_v3" "node_pool_1" {
  cluster_id         = var.cluster_id
  name               = "opentelekomcloud-cce-node-pool-test"
  os                 = "EulerOS 2.5"
  flavor             = "s2.xlarge.2"
  initial_node_count = 2
  availability_zone  = var.availability_zone
  key_pair           = var.ssh_key

  scale_enable             = true
  min_node_count           = 2
  max_node_count           = 9
  scale_down_cooldown_time = 100
  priority                 = 1

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
```

## Argument Reference
The following arguments are supported:

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `flavor` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.

* `availability_zone` - (Required) Specify the name of the available partition (AZ). If zone is not
  specified than `node_pool` will be in randomly selected AZ. The default value is `random`. Changing
  this parameter will create a new resource.

->
If AZ is set to `random`, when you create a node pool or update the number of nodes in a node pool, a scaling task is
triggered. The system selects an AZ from all AZs where scaling is allowed to add nodes based on priorities. AZs with a
smaller the number of existing nodes have a higher priority. If AZs have the same number of nodes, the system selects
the AZ based on the AZ sequence. For more details see
[API documentation](https://docs.otc.t-systems.com/en-us/api2/cce/cce_02_0354.html#cce_02_0354__table620623542313)

* `key_pair` - (Optional) Key pair name when logging in to select the key pair mode.
  This parameter and password are alternative. Changing this parameter will create a new resource.

* `password` - (Optional) Key pair name when logging in to select the key pair mode.
  This parameter and password are alternative. Changing this parameter will create a new resource.

* `os` - (Optional) Node OS. Changing this parameter will create a new resource.
  Supported OS depends on kubernetes version of the cluster.
  * Clusters of Kubernetes `v1.13` or later support `EulerOS 2.5`.
  * Clusters of Kubernetes `v1.17` or later support `EulerOS 2.5` and `CentOS 7.7`.
  * Clusters of Kubernetes `v1.21` or later support `EulerOS 2.5`, `EulerOS 2.9` and `CentOS 7.7`.

* `name` - (Required) Node Pool Name.

* `initial_node_count` - (Required) Initial number of expected nodes in the node pool.

* `subnet_id` - (Optional) The ID of the subnet to which the NIC belongs. Changing this parameter will create a new resource.

* `preinstall` - (Optional) Script required before installation. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new resource.

* `postinstall` - (Optional) Script required after installation. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new resource.

* `max_pods` - (Optional) The maximum number of instances a node is allowed to create.
  Changing this parameter will create a new node pool.

* `docker_base_size` - (Optional) Available disk space of a single Docker container on the node using the device mapper.
  Changing this parameter will create a new node pool.

* `docker_lvm_config_override` - (Optional) `ConfigMap` of the Docker data disk.
  Changing this parameter will create a new node.

* `scale_enable` - (Optional) Whether to enable auto scaling. If Autoscaler is enabled, install the autoscaler add-on to use the auto scaling feature.

* `min_node_count` - (Optional) Minimum number of nodes allowed if auto scaling is enabled.

* `max_node_count` - (Optional) Maximum number of nodes allowed if auto scaling is enabled.

* `scale_down_cooldown_time` - (Optional) Interval between two scaling operations, in minutes.

* `server_group_reference` - (Optional) ECS group ID. If this parameter is specified, all nodes in the node pool will be created in this ECS group.

* `priority` - (Optional) Weight of a node pool. A node pool with a higher weight has a higher priority during scaling.

* `user_tags` - (Optional) Tag of a VM, key/value pair format. Changing this parameter will create a new resource.

* `k8s_tags` - (Optional) Tags of a Kubernetes node, key/value pair format.

* `taints` - (Optional) Taints to created nodes to configure anti-affinity.
  * `key` - (Required) A key must contain 1 to 63 characters starting with a letter or digit. Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed. A DNS subdomain name can be used as the prefix of a key.
  * `value` - (Required) A value must start with a letter or digit and can contain a maximum of 63 characters, including letters, digits, hyphens (-), underscores (_), and periods (.).
  * `effect` - (Optional) Available options are `NoSchedule`, `PreferNoSchedule`, and `NoExecute`.

* `root_volume` - (Required) It corresponds to the system disk related configuration. Changing this parameter will create a new resource.
  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_param` - (Optional) Disk expansion parameters.
  * `kms_id` - (Optional) The Encryption KMS ID of the system volume. By default, it tries to get from env by `OS_KMS_ID`.

* `data_volumes` - (Required) Represents the data disk to be created. Changing this parameter will create a new resource.
  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_param` - (Optional) Disk expansion parameters.
  * `kms_id` - (Optional) The Encryption KMS ID of the data volume. By default, it tries to get from env by `OS_KMS_ID`.

-> To enable encryption with the KMS. Firstly, you need to create the agency to grant KMS rights to EVS.
The agency has to be created for a new project first with a user who has security `admin` permissions.
It is created automatically with the first encrypted EVS disk via UI.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `status` - Node status information.

* `id` - Specifies a resource ID in UUID format.

* `billing_mode ` - Billing mode of a node.

## Timeouts

This resource provides the following timeouts configuration options:
  - `create` - Default is 30 minutes.
  - `update` - Default is 30 minutes.
  - `delete` - Default is 30 minutes.

## Import

CCE NodePool can be imported using the `cluster_id/node_pool_id`, e.g.

```sh
terraform import opentelekomcloud_cce_node_pool_v3.pool_1 14a80bc7-c12c-4fe0-a38a-cb77eeac9bd6/89c60255-9bd6-460c-822a-e2b959ede9d2
```
