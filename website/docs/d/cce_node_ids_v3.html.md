---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_node_ids_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cce-node-ids-v3"
description: |-
  Provides a list of node Ids for a CCE cluster.
---

# Data Source: opentelekomcloud_cce_node_ids_v3

`opentelekomcloud_cce_node_ids_v3` provides a list of node ids for a CCE cluster. This resource can be useful for getting back a list of node ids for a CCE cluster.

## Example Usage

```hcl
data "opentelekomcloud_cce_node_ids_v3" "node_ids" {
  cluster_id = "${var.cluster_id}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` (Required) - Specifies the CCE cluster ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of all the node ids found. This data source will fail if none are found.
