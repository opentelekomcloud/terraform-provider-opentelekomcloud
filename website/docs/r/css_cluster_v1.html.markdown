---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_css_cluster_v1"
sidebar_current: "docs-opentelekomcloud-resource-css-cluster-v1"
description: |-
  cluster management
---

# opentelekomcloud\_css\_cluster\_v1

cluster management

## Example Usage

### Cluster

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_css_cluster_v1" "cluster" {
  initial_node_num = 1
  name = "terraform_test_cluster"
  node_config = {
    flavor = "css.medium.8"
    network_info = {
      security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup.id}"
      subnet_id = "{{ network_id }}"
      vpc_id = "{{ vpc_id }}"
    }
    volume = {
      volume_type = "COMMON"
      size = 40
    }
    availability_zone = "{{ availability_zone }}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `initial_node_num` -
  (Required)
  Number of cluster instances. The value range is 1 to 32.

* `name` -
  (Required)
  Cluster name. It contains 4 to 32 characters. Only letters, digits,
  hyphens (-), and underscores (_) are allowed. The value must start
  with a letter.

* `node_config` -
  (Required)
  Instance object. Structure is documented below.

The `node_config` block supports:

* `availability_zone` -
  (Optional)
  Availability zone (AZ).

* `flavor` -
  (Required)
  Instance flavor name. Value range of flavor css.medium.8: 40 GB
  to 640 GB Value range of flavor css.large.8: 40 GB to 1280 GB
  Value range of flavor css.xlarge.8: 40 GB to 2560 GB Value range
  of flavor css.2xlarge.8: 80 GB to 5120 GB Value range of flavor
  css.4xlarge.8: 160 GB to 10240 GB

* `network_info` -
  (Required)
  Subnet information. Structure is documented below.

* `volume` -
  (Required)
  Information about the volume. Structure is documented below.

The `network_info` block supports:

* `security_group_id` -
  (Required)
  Security group ID. All instances in a cluster must have the
  same subnets and security groups.

* `subnet_id` -
  (Required)
  Subnet ID. All instances in a cluster must have the same
  subnets and security groups.

* `vpc_id` -
  (Required)
  VPC ID, which is used for configuring cluster network.

The `volume` block supports:

* `encryption_key` -
  (Required)
  Key ID. The Default Master Keys cannot be used to create
  grants. Specifically, you cannot use Default Master Keys
  whose aliases end with /default in KMS to create clusters.
  After a cluster is created, do not delete the key used by the
  cluster. Otherwise, the cluster will become unavailable.

* `size` -
  (Required)
  Volume size, which must be a multiple of 4 and 10.

* `volume_type` -
  (Required)
  COMMON: Common I/O. The SATA disk is used. HIGH: High I/O.
  The SAS disk is used. ULTRAHIGH: Ultra-high I/O. The
  solid-state drive (SSD) is used.

- - -

* `add_node_num` -
  (Optional)
  Number of instances to be added. NOTE: The total number of existing
  instances and newly added instances in a cluster cannot exceed 32.

* `enable_https` -
  (Optional)
  Whether communication encryption is performed on the cluster.
  Available values include true and false. By default, communication
  encryption is enabled. Value true indicates that communication
  encryption is performed on the cluster. Value false indicates that
  communication encryption is not performed on the cluster.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `actions` -
  Current behavior on a cluster. Value REBOOTING indicates that the
  cluster is being restarted, GROWING indicates that capacity expansion
  is being performed on the cluster, RESTORING indicates that the
  cluster is being restored, and SNAPSHOTTINGindicates that the
  snapshot is being created.

* `created` -
  Time when a cluster is created. The format is ISO8601:
  CCYY-MM-DDThh:mm:ss.

* `datastore` -
  Type of the data search engine. Structure is documented below.

* `endpoint` -
  Indicates the IP address and port number of the user used to access
  the VPC.

* `nodes` -
  List of node objects. Structure is documented below.

* `updated` -
  Last modification time of a cluster. The format is ISO8601:
  CCYY-MM-DDThh:mm:ss.

The `datastore` block contains:

* `type` -
  (Optional)
  Supported type: elasticsearch

* `version` -
  (Optional)
  Engine version number.

The `nodes` block contains:

* `id` -
  (Optional)
  Instance ID.

* `name` -
  (Optional)
  Instance name.

* `type` -
  (Optional)
  Supported type: ess (indicating the Elasticsearch node)

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.

## Import

Cluster can be imported using the following format:

```
$ terraform import opentelekomcloud_css_cluster_v1.default {{ resource id}}
```
