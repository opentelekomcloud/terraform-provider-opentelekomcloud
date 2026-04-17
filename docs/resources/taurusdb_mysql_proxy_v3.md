---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_proxy_v3"
sidebar_current: "docs-opentelekomcloud-resource-taurusdb-mysql-proxy-v3"
description: |-
  Manages a TaurusDB MySQL proxy resource within OpenTelekomCloud.
---

# opentelekomcloud_taurusdb_mysql_proxy_v3

Manages a TaurusDB MySQL proxy resource within OpenTelekomCloud.

## Example Usage

### Basic example
```hcl
variable "instance_id" {}

resource "opentelekomcloud_taurusdb_mysql_proxy_v3" "test" {
  instance_id = var.instance_id
  flavor      = "gaussdb.proxy.xlarge.x86.2"
  node_num    = 3
}
```

### TaurusDB proxy with weights
```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "opentelekomcloud_taurusdb_mysql_instance_v3" "instance" {
  name                     = "taurus-instance"
  password                 = "Test@12345678"
  flavor                   = "gaussdb.mysql.xlarge.x86.8"
  vpc_id                   = var.vpc_id
  subnet_id                = var.subnet_id
  security_group_id        = var.security_group_id
  availability_zone_mode   = "multi"
  master_availability_zone = "eu-de-01"
  read_replicas            = 2
}

locals {
  sorted_nodes = [
    for node in opentelekomcloud_taurusdb_mysql_instance_v3.instance.nodes :
    node if node.type == "master"
  ]

  readonly_nodes = [
    for node in opentelekomcloud_taurusdb_mysql_instance_v3.instance.nodes :
    node if node.type == "slave"
  ]
}

resource "opentelekomcloud_taurusdb_mysql_proxy_v3" "test" {
  instance_id = opentelekomcloud_taurusdb_mysql_instance_v3.instance.id
  flavor      = "gaussdb.proxy.large.x86.2"
  node_num    = 3
  proxy_name  = "taurus-proxy"
  proxy_mode  = "readwrite"

  master_node_weight {
    id     = local.sorted_nodes[0].id
    weight = 40
  }

  readonly_nodes_weight {
    id     = local.readonly_nodes[0].id
    weight = 30
  }

  readonly_nodes_weight {
    id     = local.readonly_nodes[1].id
    weight = 30
  }
}
```

## Argument Reference
The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the TaurusDB MySQL instance. Changing this parameter
will create a new resource.

* `flavor` - (Required, String) Specifies the flavor of the proxy.

* `node_num` - (Required, Int) Specifies the node count of the proxy.

* `proxy_name` - (Optional, String) Specifies the name of the proxy. The name consists of `4` to `64` characters and
  starts with a letter. It is case-sensitive and can contain only letters, digits, hyphens (-), and underscores (_).

* `proxy_mode` - (Optional, String, ForceNew) Specifies the type of the proxy. Changing this creates a new resource.
  Value options:
  **readwrite**: read and write.
  **readonly**: read-only.

  Defaults to readwrite.

* `master_node_weight` - (Optional, List) Specifies the read weight of the master node.
  The [master_node_weight](#node_weight_struct) structure is documented below.

* `readonly_nodes_weight` - (Optional, Set) Specifies the read weight of the read-only nodes.
  The [readonly_nodes_weight](#node_weight_struct) structure is documented below.

<a name="node_weight_struct"></a>
The `master_node_weight` and `readonly_nodes_weight` block supports:

* `id` - (Required, String) Specifies the ID of the node.

* `weight` - (Required, Int) Specifies the weight assigned to the node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `region` - Indicates the region in which to create the resource.

* `address` - Indicates the address of the proxy.

* `port` - Indicates the port of the proxy.

* `status` - Indicates the status of the proxy.

* `nodes` - Indicates the node information of the proxy.

The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The nodes block supports:

* `id` - Indicates the proxy node ID.

* `status` - Indicates the proxy node status. The values can be:
  + **ACTIVE**: The node is available.
  + **ABNORMAL**: The node is abnormal.
  + **FAILED**: The node fails.
  + **DELETED**: The node has been deleted.

* `name` - Indicates the proxy node name.

* `role` - Indicates the proxy node role. The values can be:
  + **master**: primary node.
  + **slave**: read replica.

* `az_code` - Indicates the proxy node availability zone.

* `frozen_flag` - Indicates whether the proxy node is frozen. The values can be:
  + **0**: unfrozen.
  + **1**: frozen.
  + **2**: deleted after being frozen.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The TaurusDB MySQL proxy can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_taurusdb_mysql_proxy_v3.test <instance_id>/<id>
````

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attributes are: `proxy_mode`, `master_node_weight` and `readonly_nodes_weight`. It is
generally recommended running `terraform plan` after importing a TaurusDB MySQL proxy. You can then decide if changes
should be applied to the TaurusDB MySQL proxy, or the resource definition should be updated to align with the TaurusDB
MySQL proxy. Also you can ignore changes as below.


```hcl
resource "opentelekomcloud_taurusdb_mysql_proxy_v3" "test" {
  # ...

  lifecycle {
    ignore_changes = [
      proxy_mode, master_node_weight, readonly_nodes_weight,
    ]
  }
}
```
