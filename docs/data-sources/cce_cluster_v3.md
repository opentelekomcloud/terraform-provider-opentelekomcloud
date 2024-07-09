---
subcategory: "Cloud Container Engine (CCE)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cce_cluster_v3"
sidebar_current: "docs-opentelekomcloud-datasource-cce-cluster-v3"
description: |-
Get CCE cluster details from OpenTelekomCloud
---

Up-to-date reference of API arguments for CCE cluster you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-engine/api-ref/apis/cluster_management/listing_clusters_in_a_specified_project.html#cce-02-0239)

# opentelekomcloud_cce_cluster_v3

Use this data source to get details about all clusters and obtains the certificate for accessing cluster information.

## Example Usage

```hcl
variable "cluster_name" {}
variable "cluster_id" {}
variable "vpc_id" {}

data "opentelekomcloud_cce_cluster_v3" "cluster" {
  name   = var.cluster_name
  status = "Available"
}
```

## Argument Reference

The following arguments are supported:

* `name` -  (Optional) The Name of the cluster resource.

* `status` - (Optional) The state of the cluster.

* `cluster_type` - (Optional) Type of the cluster. Possible values: `VirtualMachine`, `BareMetal` or `Windows`.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `billingMode` - Charging mode of the cluster.

* `description` - Cluster description.

* `name` - The name of the cluster in string format.

* `id` - The ID of the cluster.

* `flavor_id` - The cluster specification in string format.

* `cluster_version` - The version of cluster in string format.

* `container_network_cidr` - The container network segment.

* `container_network_type` - The container network type: overlay_l2 , underlay_ipvlan or vpc-router.

* `eni_subnet_id` - ENI subnet ID.

* `eni_subnet_cidr` - ENI network segment.

* `authentication_mode` - (Optional) Authentication mode of the cluster, possible values are `rbac` and `authenticating_proxy`.

* `subnet_id` - The ID of the subnet used to create the node.

* `highway_subnet_id` - The ID of the high speed network used to create bare metal nodes.

* `internal` - The internal network address.

* `external` - The external network address.

* `external_otc` - The endpoint of the cluster to be accessed through API Gateway.

* `certificate_clusters/name` - The cluster name.

* `certificate_clusters/server` - The server IP address.

* `certificate_clusters/certificate_authority_data` - The certificate data.

* `certificate_users/name` - The user name.

* `certificate_users/client_certificate_data` - The client certificate data.

* `certificate_users/client_key_data` - The client key data.
