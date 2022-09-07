---
subcategory: "Cloud Container Engine (CCE)"
---

# opentelekomcloud_cce_cluster_v3

Provides a cluster resource management.

~>
  Before starting working with CCE, you need to authorize it via _console_ or [creating agency](#creating-agency).
  Otherwise, you will face the following error during cluster creation:
  `CCE is not authorized, see `cce_cluster_v3` documentation for details`.

~>
  You need to authorize CCE for the default (`eu-de`) project for CCE to be able to pull SWR images.

## Example Usage

### Simple cluster

```hcl
variable "flavor_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name        = "cluster"
  description = "Create cluster"

  cluster_type           = "VirtualMachine"
  flavor_id              = var.flavor_id
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  kube_proxy_mode        = "ipvs"
}
```

### Creating agency

You can create agency for CCE authorization using `opentelekomcloud_identity_agency_v3` resource.
For agency creation your user need to have corresponding permissions, which are not required for authorizing CCE via console

```hcl
resource "opentelekomcloud_identity_agency_v3" "enable_cce_auto_creation" {
  name                  = "cce_admin_trust"
  description           = "Created by Terraform to auto create cce"
  delegated_domain_name = "op_svc_cce"
  dynamic "project_role" {
    for_each = var.projects
    content {
      project = project_role.value
      roles = [
        "Tenant Administrator"
      ]
    }
  }
}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name        = "cluster"
  description = "Create cluster"

  cluster_type           = "VirtualMachine"
  flavor_id              = var.flavor_id
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"

  depends_on = [opentelekomcloud_identity_agency_v3.enable_cce_auto_creation]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Cluster name. Changing this parameter will create a new cluster resource.

* `labels` - (Optional) Cluster tag, key/value pair format. Changing this parameter will create a new cluster resource.

* `annotations` - (Optional) Cluster annotation, key/value pair format. Changing this parameter will create a new cluster resource.

* `flavor_id` - (Required) Cluster specifications. Changing this parameter will create a new cluster resource.
  * `cce.s1.small` - small-scale single cluster (up to 50 nodes).
  * `cce.s1.medium` - medium-scale single cluster (up to 200 nodes).
  * `cce.s1.large` - large-scale single cluster (up to 1000 nodes).
  * `cce.s2.small` - small-scale HA cluster (up to 50 nodes).
  * `cce.s2.medium` - medium-scale HA cluster (up to 200 nodes).
  * `cce.s2.large` - large-scale HA cluster (up to 1000 nodes).
  * `cce.t1.small` - small-scale single physical machine cluster (up to 10 nodes).
  * `cce.t1.medium` - medium-scale single physical machine cluster (up to 100 nodes).
  * `cce.t1.large` - large-scale single physical machine cluster (up to 500 nodes).
  * `cce.t2.small` - small-scale HA physical machine cluster (up to 10 nodes).
  * `cce.t2.medium` - medium-scale HA physical machine cluster (up to 100 nodes).
  * `cce.t2.large` - large-scale HA physical machine cluster (up to 500 nodes).

* `cluster_version` - (Optional) For the cluster version, possible values are `v1.17.9-r0`, `v1.19.8-r0`.
  Changing this parameter will create a new cluster resource. [OTC-API](https://docs.otc.t-systems.com/en-us/api2/cce/cce_02_0236.html)

* `cluster_type` - (Required) Cluster Type, possible values are `VirtualMachine` and `BareMetal`. Changing this parameter will create a new cluster resource.

* `description` - (Optional) Cluster description.

* `billing_mode` - (Optional) Charging mode of the cluster, which is 0 (on demand). Changing this parameter will create a new cluster resource.

* `extend_param` - (Optional) Extended parameter. Changing this parameter will create a new cluster resource.

* `vpc_id` - (Required) The ID of the VPC used to create the node. Changing this parameter will create a new cluster resource.

* `subnet_id` - (Required) The Network ID of the subnet used to create the node. Changing this parameter will create a new cluster resource.

* `highway_subnet_id` - (Optional) The ID of the high speed network used to create bare metal nodes. Changing this parameter will create a new cluster resource.

* `container_network_type` - (Required) Container network type.
  * `overlay_l2` - An overlay_l2 network built for containers by using Open vSwitch(OVS)
  * `underlay_ipvlan` - An underlay_ipvlan network built for bare metal servers by using ipvlan.
  * `vpc-router` - An vpc-router network built for containers by using ipvlan and custom VPC routes.

* `container_network_cidr` - (Optional) Container network segment. Changing this parameter will create a new cluster resource.

* `eni_subnet_id` -  - (Optional) Specifies the ENI subnet ID. Specified when creating a CCE Turbo cluster. Changing this parameter will create a new cluster resource.

* `eni_subnet_cidr` - (Optional) Specifies the ENI network segment. Specified when creating a CCE Turbo cluster. Changing this parameter will create a new cluster resource.

* `authentication_mode` - (Optional) Authentication mode of the cluster, possible values are `rbac` and `authenticating_proxy`.
  Defaults to `rbac`. Changing this parameter will create a new cluster resource.

* `authenticating_proxy_ca` - (Optional) CA root certificate provided in the `authenticating_proxy` mode.
  Deprecated, use `authenticating_proxy` instead.

* `authenticating_proxy` - (Optional) Authenticating proxy configuration. Required if `authentication_mode` is set to `authenticating_proxy`.
  * `ca` - X509 CA certificate configured in `authenticating_proxy` mode. The maximum size of the certificate is 1 MB.
  * `cert` - Client certificate issued by the X509 CA certificate configured in `authenticating_proxy` mode.
  This certificate is used for authentication from kube-apiserver to the extended API server.
  * `private_key` - Private key of the client certificate issued by the X509 CA certificate configured in `authenticating_proxy` mode.
  This key is used for authentication from kube-apiserver to the extended API server.

~>
  The private key used by the Kubernetes cluster does not support password encryption. Use an unencrypted private key.

* `multi_az` - (Optional) Enable multiple AZs for the cluster, only when using HA flavors. Changing this parameter will create a new cluster resource.

* `eip` - (Optional) EIP address of the cluster.

* `kubernetes_svc_ip_range` - (Optional) Service CIDR block, or the IP address range which the kubernetes
  clusterIp must fall within. This parameter is available only for clusters of v1.11.7 and later.

* `no_addons` - (Optional) Remove addons installed by the default after the cluster creation.

* `ignore_addons` - (Optional) Skip all cluster addons operations.

* `kube_proxy_mode` - Service forwarding mode. Two modes are available:
  * `iptables`: Traditional kube-proxy uses iptables rules to implement service load balancing.
    In this mode, too many iptables rules will be generated when many services are deployed.
    In addition, non-incremental updates will cause a latency and even obvious performance issues
    in the case of heavy service traffic.
  * `ipvs`: Optimized kube-proxy mode with higher throughput and faster speed.
    This mode supports incremental updates and can keep connections uninterrupted during service updates.
    It is suitable for large-sized clusters.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `id` - ID of the cluster resource.

* `status` - Cluster status information.

* `internal` - The internal network address.

* `external` - The external network address.

* `external_otc` - The endpoint of the cluster to be accessed through API Gateway.

* `certificate_clusters/name` - The cluster name.

* `certificate_clusters/server` - The server IP address.

* `certificate_clusters/certificate_authority_data` - The certificate data.

* `certificate_users/name` - The user name.

* `certificate_users/client_certificate_data` - The client certificate data.

* `certificate_users/client_key_data` - The client key data.

* `installed_addons` - List of installed addon IDs. Empty if `ignore_addons` is `true`.

* `security_group_control` - ID of the autogenerated security group for the CCE master port.

* `security_group_node` - ID of the autogenerated security group for the CCE nodes.

## Timeouts

This resource provides the following timeouts configuration options:

- `create` - Default is 30 minutes.

- `delete` - Default is 30 minutes.

## Import

Cluster can be imported using the cluster id, e.g.

```shell
terraform import opentelekomcloud_cce_cluster_v3.cluster_1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
