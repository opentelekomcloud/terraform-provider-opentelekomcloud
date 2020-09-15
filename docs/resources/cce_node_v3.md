---
subcategory: "Cloud Container Engine (CCE)"
---

# opentelekomcloud_cce_node_v3

Add a node to a container cluster.

## Example Usage

```hcl
variable "cluster_id" { }
variable "ssh_key" { }
variable "availability_zone" { }

resource "opentelekomcloud_cce_node_v3" "node_1" {
  name              = "node1"
  cluster_id        = var.cluster_id
  availability_zone = var.availability_zone

  flavor_id      = "s1.medium"
  key_pair       = var.ssh_key
  iptype         = "5_bgp"
  sharetype      = "PER"
  bandwidth_size = 100

  root_volume {
    size       = 40
    volumetype = "SATA"
  }

  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
```

## Argument Reference
The following arguments are supported:

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `billing_mode` - (Optional) Node's billing mode: The value is 0 (on demand). Changing this parameter will create a new resource.

* `name` - (Optional) Node Name.

* `labels` - (Optional) Node tag, key/value pair format. Changing this parameter will create a new resource.

* `tags` - (Optional) The field is alternative to `labels`, key/value pair format.

* `annotations` - (Optional) Node annotation, key/value pair format. Changing this parameter will create a new resource.

* `flavor_id` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.

* `availability_zone` - (Required) specify the name of the available partition (AZ). Changing this parameter will create a new resource.

* `key_pair` - (Required) Key pair name when logging in to select the key pair mode. Changing this parameter will create a new resource.

* `eip_ids` - (Optional) List of existing elastic IP IDs. Changing this parameter will create a new resource.

-> **Note:** If the `eip_ids` parameter is configured, you do not need to configure the `eip_count` and `bandwidth` parameters:
`iptype`, `charge_mode`, `bandwidth_size` and `share_type`.

* `eip_count` - (Optional) Number of elastic IPs to be dynamically created. Changing this parameter will create a new resource.

* `iptype` - (Optional) Elastic IP type. Default value is `5_bgp`.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Default value is `traffic`. Changing this parameter will create a new resource.

* `sharetype` - (Optional) Bandwidth sharing type. Default value is `PER` Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional) Bandwidth size. Changing this parameter will create a new resource.

* `extend_param_charging_mode` - (Optional) Node charging mode, 0 is on-demand charging. Changing this parameter will create a new cluster resource.

* `ecs_performance_type` - (Optional) Classification of cloud server specifications. Changing this parameter will create a new cluster resource.

* `order_id` - (Optional) Order ID, mandatory when the node payment type is the automatic payment package period type.
  Changing this parameter will create a new cluster resource.

* `product_id` - (Optional) The Product ID. Changing this parameter will create a new cluster resource.

* `max_pods` - (Optional) The maximum number of instances a node is allowed to create. Changing this parameter will create a new cluster resource.

* `public_key` - (Optional) The Public key. Changing this parameter will create a new cluster resource.

* `preinstall` - (Optional) Script required before installation. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new resource.

* `postinstall` - (Optional) Script required after installation. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new resource.

* `root_volume` - (Required) It corresponds to the system disk related configuration. Changing this parameter will create a new resource.
  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_param` - (Optional) Disk expansion parameters.

* `data_volumes` - (Required) Represents the data disk to be created. Changing this parameter will create a new resource.
  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_param` - (Optional) Disk expansion parameters.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `status` - Node status information.

* `server_id` - ID of the ECS where the node resides.

* `private_ip` - Private IP of the CCE node.

* `public_ip` - Public IP of the CCE node.
