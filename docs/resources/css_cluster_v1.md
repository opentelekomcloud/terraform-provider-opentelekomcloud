---
subcategory: "Cloud Search Service (CSS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_css_cluster_v1"
sidebar_current: "docs-opentelekomcloud-resource-css-cluster-v1"
description: |-
Manages a CSS Cluster resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CSS cluster you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-search-service/api-ref/cluster_management_apis)

# opentelekomcloud_css_cluster_v1

Manages a CSS cluster resource.

## Example usage

### Creating ElasticSearch v7.6.2

```hcl
data "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = var.security_group
}

resource "opentelekomcloud_css_cluster_v1" "cluster" {
  name            = "terraform_test_cluster"
  expect_node_num = 1
  node_config {
    flavor = "css.medium.8"
    network_info {
      security_group_id = data.opentelekomcloud_networking_secgroup_v2.secgroup.id
      network_id        = var.network_id
      vpc_id            = var.vpc_id
    }
    volume {
      volume_type = "COMMON"
      size        = 40
    }

    availability_zone = var.availability_zone
  }
}
```

### Creating OpenSearch v1.3.6

```hcl
data "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name = var.security_group
}

resource "opentelekomcloud_css_cluster_v1" "cluster" {
  name            = "terraform_test_cluster"
  expect_node_num = 1
  datastore {
    version = "1.3.6"
    type    = "opensearch"
  }
  enable_https     = true
  enable_authority = true
  admin_pass       = "QwertyUI!"
  node_config {
    flavor = "css.medium.8"
    network_info {
      security_group_id = data.opentelekomcloud_networking_secgroup_v2.secgroup.id
      network_id        = var.network_id
      vpc_id            = var.vpc_id
    }
    volume {
      volume_type = "COMMON"
      size        = 40
    }

    availability_zone = var.availability_zone
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Cluster name. It contains `4` to `32` characters. Only letters, digits,
  hyphens (`-`), and underscores (`_`) are allowed. The value must start with a letter.
  Changing this parameter will create a new resource.

* `datastore` - (Optional) Type of the data search engine. Structure is documented below.
  Changing this parameter will create a new resource.

* `node_config` - (Required) Instance object. Structure is documented below.
  Changing this parameter will create a new resource.

* `enable_https` - (Optional) Whether communication encryption is performed on the cluster.
  By default, communication encryption is enabled.
  Value `true` indicates that communication encryption is performed on the cluster.
  Value `false` indicates that communication encryption is not performed on the cluster.
  Changing this parameter will create a new resource.

* `enable_authority` - (Optional) Whether to enable authentication.
  Authentication is disabled by default. When authentication is enabled, `enable_https` must be set to `true`.
  Changing this parameter will create a new resource.

* `admin_pass` - (Optional) Password of the cluster user admin in security mode.
  This parameter is mandatory only when `enable_authority` is set to `true`.

~>
The administrator password must meet the following requirements: contain `8` to `32` characters,
contain at least `3` of the following character types: uppercase letters,
lowercase letters, numbers, and special characters (`~!@#$%^&*()-_=+\\|[{}];:,<.>/?`).

* `expect_node_num` - (Optional) Number of cluster instances. The value range is `1` to `32`.

* `tags` - (Optional) Tags key/value pairs to associate with the cluster.

The `node_config` block supports:

* `availability_zone` - (Optional) Availability zone (AZ). Changing this parameter will create a new resource.

* `flavor` - (Required) Instance flavor name.
  - Value range of flavor `css.medium.8`: 40 GB to 640 GB
  - Value range of flavor `css.xlarge.8`: 40 GB to 2560 GB
  - Value range of flavor `css.2xlarge.8`: 80 GB to 5120 GB
  - Value range of flavor `css.4xlarge.8`: 160 GB to 10240 GB
  - Value range of flavor `css.xlarge.4`: 40 GB to 1,600 GB
  - Value range of flavor `css.2xlarge.4`: 80 GB to 3,200 GB
  - Value range of flavor `css.4xlarge.4`: 100 GB to 6,400 GB
  - Value range of flavor `css.8xlarge.4`: 160 GB to 10,240 GB
  - Value range of flavor `css.xlarge.2`: 40 GB to 800 GB
  - Value range of flavor `css.2xlarge.2`: 80 GB to 1,600 GB
  - Value range of flavor `css.4xlarge.2`: 100 GB to 3,200 GB
  - Value range of flavor `css.8xlarge.2`: 320 GB to 10,240 GB

  Changing this parameter will create a new resource.

* `network_info` - (Required) Network information. Structure is documented below.
  Changing this parameter will create a new resource.

* `volume` - (Required) Information about the volume. Structure is documented below.
  Changing this parameter will create a new resource.

The `network_info` block supports:

* `network_id` - (Required) Network ID. All instances in a cluster must have the same
  networks and security groups. Changing this parameter will create a new resource.

* `security_group_id` - (Required) Security group ID. All instances in a cluster must have the
  same subnets and security groups. Changing this parameter will create a new resource.

* `vpc_id` - (Required) VPC ID, which is used for configuring cluster network.
  Changing this parameter will create a new resource.

The `volume` block supports:

* `encryption_key` - (Optional) Key ID. The Default Master Keys cannot be used to create
  grants. Specifically, you cannot use Default Master Keys
  whose aliases end with /default in KMS to create clusters.
  After a cluster is created, do not delete the key used by the cluster.
  Otherwise, the cluster will become unavailable.
  Changing this parameter will create a new resource.

* `size` - (Required) Volume size, which must be a multiple of `4` and `10`.

* `volume_type` - (Required) `COMMON`: Common I/O. The SATA disk is used. `HIGH`: High I/O.
  The SAS disk is used. `ULTRAHIGH`: Ultra-high I/O. The solid-state drive (SSD) is used.
  Changing this parameter will create a new resource.

The `datastore` block contains:

* `type` - Engine type. The default value is `elasticsearch`. Currently, the value can be `elasticsearch` or `opensearch`.

* `version` - Engine version. The value can be `7.6.2`, `7.9.3`, `7.10.2` or `1.3.6` for `opensearch`.
  The default value is `7.6.2`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created` - Time when a cluster is created. The format is ISO8601: `CCYY-MM-DDThh:mm:ss`.

* `endpoint` - Indicates the IP address and port number of the user used to access the VPC.

* `nodes` - List of node objects. Structure is documented below.

* `updated` - Last modification time of a cluster. The format is ISO8601: `CCYY-MM-DDThh:mm:ss`.

The `nodes` block contains:

* `id` - Instance ID.

* `name` - Instance name.

* `type` - Supported type: `ess` (indicating the Elasticsearch node)

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.

* `update` - Default is 30 minutes.

## Import

Backup can be imported using  `cluster_id`, e.g.

```sh
terraform import opentelekomcloud_css_cluster_v1.cluster 5c77b71c-5b35-4f50-8984-76387e42451a
```
