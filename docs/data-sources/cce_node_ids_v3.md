---
subcategory: "Cloud Container Engine (CCE)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_node_ids_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cce-node-ids-v3"
description: |-
Get a list of node ids for a CCE cluster from OpenTelekomCloud
---

Up-to-date reference of API arguments for CCE nodes you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-engine/api-ref/apis/cluster_management/listing_all_nodes_in_a_cluster.html)

# opentelekomcloud_cce_node_ids_v3

Use this data source to get a list of node ids for a CCE cluster from OpenTelekomCloud.
This data source can be useful for getting back a list of node ids for a CCE cluster.

## Example Usage

```hcl
data "opentelekomcloud_cce_node_ids_v3" "node_ids" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) Specifies the CCE cluster ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of all the node ids found. This data source will fail if none are found.
