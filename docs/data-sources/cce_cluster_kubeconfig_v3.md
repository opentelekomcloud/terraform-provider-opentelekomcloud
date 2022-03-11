---
subcategory: "Cloud Container Engine (CCE)"
---

# opentelekomcloud_cce_cluster_kubeconfig_v3

Use this data source to get details about a cluster kubeconfig file from OpenTelekomCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "opentelekomcloud_cce_cluster_kubeconfig_v3" "this" {
  name   = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` -  (Optional) The Name of the cluster resource.

* `duration` - (Optional) Period during which a cluster certificate is valid, in days. A cluster certificate can
  be valid for `1` to `1825` days. If this parameter is set to `-1`, the validity period is `1825` days (about 5 years).
  Default vault `-1`.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `id` - The ID of the cluster.

* `kubeconfig` - The kubeconfig file of the cluster.
