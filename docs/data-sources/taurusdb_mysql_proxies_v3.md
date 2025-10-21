---
subcategory: "TaurusDB(for MySQL)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_taurusdb_mysql_proxies_v3"
sidebar_current: "docs-opentelekomcloud-datasource-taurusdb-mysql-proxies-v3"
description: |-
  Get the list of TaurusDB MySQL proxies from OpenTelekomCloud
---

# opentelekomcloud_taurusdb_mysql_proxies_v3

Use this data source to get the list of TaurusDB MySQL proxies.

## Example Usage

```hcl
variable "instance_id" {}

data "opentelekomcloud_taurusdb_mysql_proxies_v3" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB MySQL instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which to query the resource.

* `proxy_list` - Indicates the list of proxies.
  The [proxy_list](#proxy_list) structure is documented below.

<a name="proxy_list"></a>
The `proxy_list` block supports:

* `id` - Indicates the proxy ID.

* `name` - Indicates the proxy name.

* `flavor` - Indicates the flavor of the proxy.

* `port` - Indicates the proxy port.

* `status` - Indicates the status of the proxy instance.

* `delay_threshold_in_seconds` - Indicates the delay threshold in seconds.

* `node_num` - Indicates the number of proxy nodes.

* `ram` - Indicates the memory size of the proxy.

* `mode` - Indicates the proxy mode.

* `elb_vip` - Indicates the virtual IP address in ELB mode.

* `vcpus` - Indicates the number of vCPUs of the proxy.

* `transaction_split` - Indicates whether the proxy transaction splitting is enabled.

* `address` - Indicates the address of the proxy.

* `nodes` - Indicates the node information of the proxy.
  The [nodes](#proxy_nodes) structure is documented below.

* `master_node_weight` - Indicates the read weight of the master node.
  The [master_node_weight](#master_node_weight) structure is documented below.

* `readonly_nodes_weight` - Indicates the read weight of the read-only node.
  The [readonly_nodes_weight](#readonly_nodes_weight) structure is documented below.

<a name="proxy_nodes"></a>
The `nodes` block supports:

* `id` - Indicates the proxy node ID.

* `name` - Indicates the proxy node name.

* `role` - Indicates the proxy node role.

* `az_code` - Indicates the proxy node AZ.

* `status` - Indicates the proxy node status.

* `frozen_flag` - Indicates whether the proxy node is frozen.

<a name="master_node_weight"></a>
The `master_node_weight` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `weight` - Indicates the weight assigned to the node.

<a name="readonly_nodes_weight"></a>
The `readonly_nodes_weight` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `weight` - Indicates the weight assigned to the node.
