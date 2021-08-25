---
subcategory: "Cloud Container Engine (CCE)"
---

# opentelekomcloud_cce_cluster_v3

Use this data source to get details about all clusters and obtains the certificate for accessing cluster information.

## Example Usage

```hcl
variable "cluster_name" {}
variable "cluster_id" {}

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

* `router_id` - (Optional) ID of the router (VPC).

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

* `authentication_mode` - (Optional) Authentication mode of the cluster, possible values are `rbac` and `authenticating_proxy`.

* `network_id` - The ID of the subnet used to create the node.

* `highway_network_id` - The ID of the high speed network used to create bare metal nodes.

* `internal` - The internal network address.

* `external` - The external network address.

* `external_otc` - The endpoint of the cluster to be accessed through API Gateway.

* `certificate_clusters/name` - The cluster name.

* `certificate_clusters/server` - The server IP address.

* `certificate_clusters/certificate_authority_data` - The certificate data.

* `certificate_users/name` - The user name.

* `certificate_users/client_certificate_data` - The client certificate data.

* `certificate_users/client_key_data` - The client key data.
