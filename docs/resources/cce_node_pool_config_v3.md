---
subcategory: "Cloud Container Engine (CCE)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_node_pool_config_v3"
sidebar_current: "docs-opentelekomcloud-resource-cce-node-pool-config-v3"
description: |-
  Manages a CCE node pool configuration resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCE node pool configuration you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-engine/api-ref/apis/configuration_management)

# opentelekomcloud_cce_node_pool_config_v3

Provides a node pool configuration resource of a CCE container cluster.

## Example Usage

### Basic example

```hcl
variable "cluster_id" {}
variable "nodepool_id" {}

resource "opentelekomcloud_cce_node_pool_config_v3" "node_pool_config" {
  cluster_id  = var.cluster_id
  nodepool_id = var.nodepool_id
  name        = "configuration"

  packages {
    name = "kubelet"
    configurations {
      name  = "system-reserved-mem"
      value = 600
    }
  }
}
```

## Argument Reference
The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the CCE cluster.

* `nodepool_id` - (Required, String, ForceNew) Specifies the node pool ID.

* `name` - (Requried, String) Specifies the configuration name.

* `labels` - (Optional, Map) Specifies the configuration labels in a key-value pair.

* `packages` - (Required, List) Specifies the vomponent configuration item details. The [packages](#packages) structure is documented below.

<!-- ### PACKAGES ### -->
<a name="packages"></a>
The `packages` block supports:

* `name` - (Required, String) Specifies the component/package name.

* `configurations` - (Required, List) Specifies configuration items. The structure is documented below:
    * `name` - (Required, String) Specifies the configuration item name one wishes to override. This configuration item will be ignored if unsupported component or parameter is specified.
    * `value` - (Required) Specifies the configuration item value. This configuration item will be ignored if unsupported component or parameter is specified.

-> **Note:** Supported values for package names and their corresponding configurable items can be found at this page in [CCE user guide](https://docs.otc.t-systems.com/cloud-container-engine/umn/node_pools/managing_node_pools/modifying_node_pool_configurations.html#containerd-available-only-for-containerd-node-pools)

-> **Note:** For `registry-mirrors` parameter, use the semicolon separated string in `value`, for example, `"example.com=https://my-mirror.example.com;"` or `"example.com=https://my-mirror.example.com;old-mirror.com=https://new-mirror.com;"`. Semicolon `;` is necessary to distinguish the array from strings. The `https` is required for new value after `=` and the trailing `/` must be avoided.


## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `id` - Specifies a resource ID in UUID format.
