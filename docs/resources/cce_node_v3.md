---
subcategory: "Cloud Container Engine (CCE)"
---

# opentelekomcloud_cce_node_v3

Add a node to a container cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "ssh_key" {}
variable "availability_zone" {}

resource "opentelekomcloud_cce_node_v3" "node_1" {
  name              = "node1"
  cluster_id        = var.cluster_id
  availability_zone = var.availability_zone

  os        = "EulerOS 2.5"
  flavor_id = "s2.large.2"
  key_pair  = var.ssh_key

  bandwidth_size = 100

  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    extend_params = {
      "useType" = "docker"
    }
  }
}
```

## Argument Reference
The following arguments are supported:

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `flavor_id` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.

* `availability_zone` - (Required) specify the name of the available partition (AZ). Changing this parameter will create a new resource.

* `key_pair` - (Required) Key pair name when logging in to select the key pair mode. Changing this parameter will create a new resource.

* `os` - (Optional) Node OS. Changing this parameter will create a new resource.

  Supported OS depends on kubernetes version of the cluster.
  * Clusters of Kubernetes `v1.13` or later support `EulerOS 2.5`.
  * Clusters of Kubernetes `v1.17` or later support `EulerOS 2.5` and `CentOS 7.7`.
  * Clusters of Kubernetes `v1.21` or later support `EulerOS 2.5`, `EulerOS 2.9` and `CentOS 7.7`.

* `billing_mode` - (Optional) Node's billing mode: The value is `0` (on demand). Changing this parameter will create a new resource.

* `name` - (Optional) Node Name.

* `subnet_id` - (Optional) The ID of the subnet to which the NIC belongs. Changing this parameter will create a new resource.

* `labels` - (Optional) Node tag, key/value pair format. Changing this parameter will create a new resource.

* `tags` - (Optional) The field is alternative to `labels`, key/value pair format.

* `k8s_tags` - (Optional) Tags of a Kubernetes node, key/value pair format.

* `annotations` - (Optional) Node annotation, key/value pair format. Changing this parameter will create a new resource

* `runtime` - (Optional) Container runtime. Changing this parameter will create a new resource. Options are:
              `docker` - Docker
              `containerd` - Containerd

* `taints` - (Optional) Taints to created nodes to configure anti-affinity.
  * `key` - (Required) A key must contain 1 to 63 characters starting with a letter or digit. Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed. A DNS subdomain name can be used as the prefix of a key.
  * `value` - (Required) A value must start with a letter or digit and can contain a maximum of 63 characters, including letters, digits, hyphens (-), underscores (_), and periods (.).
  * `effect` - (Optional) Available options are `NoSchedule`, `PreferNoSchedule`, and `NoExecute`.

* `eip_ids` - (Optional) List of existing elastic IP IDs.

-> If the `eip_ids` parameter is configured, you do not need to configure the `eip_count` and `bandwidth` parameters:
`iptype`, `bandwidth_charge_mode`, `bandwidth_size` and `share_type`.

* `eip_count` - (Optional) Number of elastic IPs to be dynamically created.

* `iptype` - (Optional) Elastic IP type.

* `bandwidth_size` - (Optional) Bandwidth size.

-> If the `bandwidth_size` parameter is configured, you do not need to configure the
  `eip_count`, `bandwidth_charge_mode`, `sharetype` and `iptype` parameters.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type.

* `sharetype` - (Optional) Bandwidth sharing type.

* `extend_param_charging_mode` - (Optional) Node charging mode, 0 is on-demand charging. Changing this parameter will create a new cluster resource.

* `ecs_performance_type` - (Optional) Classification of cloud server specifications. Changing this parameter will create a new cluster resource.

* `order_id` - (Optional) Order ID, mandatory when the node payment type is the automatic payment package period type.
  Changing this parameter will create a new cluster resource.

* `product_id` - (Optional) The Product ID. Changing this parameter will create a new cluster resource.

* `max_pods` - (Optional) The maximum number of instances a node is allowed to create. Changing this parameter will create a new node resource.

* `public_key` - (Optional) The Public key. Changing this parameter will create a new cluster resource.

* `private_ip` - (Optional) Private IP of the CCE node. Changing this parameter will create a new resource.

* `preinstall` - (Optional) Script required before installation. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new resource.

* `postinstall` - (Optional) Script required after installation. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new resource.

* `docker_base_size` - (Optional) Available disk space of a single Docker container on the node using the device mapper.
  Changing this parameter will create a new node.

* `docker_lvm_config_override` - (Optional) `ConfigMap` of the Docker data disk.
  Changing this parameter will create a new node.

  Example:

  `dockerThinpool=vgpaas/90%VG;kubernetesLV=vgpaas/10%VG;diskType=evs;lvType=linear`

  In this example:

  - `userLV`: size of the user space, for example, vgpaas/20%VG.
  - `userPath`: mount path of the user space, for example, /home/wqt-test.
  - `diskType`: disk type. Currently, only the evs, hdd, and ssd are supported.
  - `lvType`: type of a logic volume. Currently, the value can be linear or striped.
  - `dockerThinpool`: Docker space size, for example, vgpaas/60%VG.
  - `kubernetesLV`: kubelet space size, for example, vgpaas/20%VG.

* `root_volume` - (Required) It corresponds to the system disk related configuration. Changing this parameter will create a new resource.
  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_params` - (Optional) Disk expansion parameters. A list of strings which describes additional disk parameters.
  * `extend_param` **DEPRECATED** - (Optional) Disk expansion parameters.
  Please use alternative parameter `extend_params`.
  * `kms_id` - (Optional) The Encryption KMS ID of the system volume. By default, it tries to get from env by `OS_KMS_ID`.

* `data_volumes` - (Required) Represents the data disk to be created. Changing this parameter will create a new resource.
  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_params` - (Optional) Disk expansion parameters. A list of strings which describes additional disk parameters.
  * `extend_param` **DEPRECATED** - (Optional) Disk expansion parameters.
  Please use alternative parameter `extend_params`.
  * `kms_id` - (Optional) The Encryption KMS ID of the data volume. By default, it tries to get from env by `OS_KMS_ID`.

-> To enable encryption with the KMS. Firstly, you need to create the agency to grant KMS rights to EVS.
The agency has to be created for a new project first with a user who has security `admin` permissions.
It is created automatically with the first encrypted EVS disk via UI.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `status` - Node status information.

* `server_id` - ID of the ECS where the node resides.

* `public_ip` - Public IP of the CCE node.

## Timeouts

This resource provides the following timeouts configuration options:

- `create` - Default is 10 minutes.

- `delete` - Default is 10 minutes.
