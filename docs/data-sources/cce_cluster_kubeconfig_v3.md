---
subcategory: "Cloud Container Engine (CCE)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_cluster_kubeconfig_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cce-cluster-kubeconfig-v3"
description: |-
Get CCE a cluster's kubeconfig file from OpenTelekomCloud
---

# opentelekomcloud_cce_cluster_kubeconfig_v3

Use this data source to get a cluster's kubeconfig file from OpenTelekomCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "opentelekomcloud_cce_cluster_kubeconfig_v3" "this" {
  cluster_id = var.cluster_id
}
```

## Example with expiration date

```hcl
variable "cluster_id" {}

data "opentelekomcloud_cce_cluster_kubeconfig_v3" "this" {
  cluster_id  = opentelekomcloud_cce_cluster_v3.cluster_1.id
  expiry_date = "2024-02-01"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` -  (Required, String) The Name of the cluster resource.

* `duration` - (Optional, Int) Period during which a cluster certificate is valid, in days. A cluster certificate can
  be valid for `1` to `1825` days. If this parameter is set to `-1`, the validity period is `1825` days (about 5 years).
  Default vault `-1`.

* `expiry_date` - (Optional, String) Specifies the date until which the certificate will be valid, in RFC3339 format, like `2023-02-01`.
  Conflicts with `duration` attribute.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `id` - The ID of the cluster.

* `kubeconfig` - The cluster's kubeconfig file contents.
