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

### create a cluster

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_css_cluster_v1" "cluster" {
  expect_node_num = 1
  name = "terraform_test_cluster"
  node_config {
    flavor = "css.medium.8"
    network_info {
      security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup.id}"
      network_id = "{{ network_id }}"
      vpc_id = "{{ vpc_id }}"
    }
    volume {
      volume_type = "COMMON"
      size = 40
    }
    availability_zone = "{{ availability_zone }}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Cluster name. It contains 4 to 32 characters. Only letters, digits,
  hyphens (-), and underscores (_) are allowed. The value must start
  with a letter.  Changing this parameter will create a new resource.

* `node_config` -
  (Required)
  Instance object. Structure is documented below. Changing this parameter will create a new resource.

The `node_config` block supports:

* `availability_zone` -
  (Optional)
  Availability zone (AZ).  Changing this parameter will create a new resource.

* `flavor` -
  (Required)
  Instance flavor name. Value range of flavor css.medium.8: 40 GB
  to 640 GB Value range of flavor css.large.8: 40 GB to 1280 GB
  Value range of flavor css.xlarge.8: 40 GB to 2560 GB Value range
  of flavor css.2xlarge.8: 80 GB to 5120 GB Value range of flavor
  css.4xlarge.8: 160 GB to 10240 GB  Changing this parameter will create a new resource.

* `network_info` -
  (Required)
  Network information. Structure is documented below. Changing this parameter will create a new resource.

* `volume` -
  (Required)
  Information about the volume. Structure is documented below. Changing this parameter will create a new resource.

The `network_info` block supports:

* `network_id` -
  (Required)
  Network ID. All instances in a cluster must have the same
  networks and security groups.  Changing this parameter will create a new resource.

* `security_group_id` -
  (Required)
  Security group ID. All instances in a cluster must have the
  same subnets and security groups.  Changing this parameter will create a new resource.

* `vpc_id` -
  (Required)
  VPC ID, which is used for configuring cluster network.  Changing this parameter will create a new resource.

The `volume` block supports:

* `encryption_key` -
  (Optional)
  Key ID. The Default Master Keys cannot be used to create
  grants. Specifically, you cannot use Default Master Keys
  whose aliases end with /default in KMS to create clusters.
  After a cluster is created, do not delete the key used by the
  cluster. Otherwise, the cluster will become unavailable.  Changing this parameter will create a new resource.

* `size` -
  (Required)
  Volume size, which must be a multiple of 4 and 10.  Changing this parameter will create a new resource.

* `volume_type` -
  (Required)
  COMMON: Common I/O. The SATA disk is used. HIGH: High I/O.
  The SAS disk is used. ULTRAHIGH: Ultra-high I/O. The
  solid-state drive (SSD) is used.  Changing this parameter will create a new resource.

- - -

* `enable_https` -
  (Optional)
  Whether communication encryption is performed on the cluster.
  Available values include true and false. By default, communication
  encryption is enabled. Value true indicates that communication
  encryption is performed on the cluster. Value false indicates that
  communication encryption is not performed on the cluster.  Changing this parameter will create a new resource.

* `expect_node_num` -
  (Optional)
  Number of cluster instances. The value range is 1 to 32.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

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
  Supported type: elasticsearch

* `version` -
  Engine version number.

The `nodes` block contains:

* `id` -
  Instance ID.

* `name` -
  Instance name.

* `type` -
  Supported type: ess (indicating the Elasticsearch node)

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 15 minute.
- `update` - Default is 30 minute.
