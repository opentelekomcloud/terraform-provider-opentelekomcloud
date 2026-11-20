---
subcategory: "GeminiDB"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gemini_instance_v3"
sidebar_current: "docs-opentelekomcloud-resource-gemini-instance-v3"
description: |-
  Manages a GeminiDB resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for GeminiDB you can get at
[documentation portal](https://docs.otc.t-systems.com/geminidb/api-ref/apis_v3/instances/index.html).

# opentelekomcloud_gemini_instance_v3

GeminiDB instance management within OpenTelekomCloud.

Both GeminiDB variants are managed by this resource, the `datastore` block selects the compatible API:

| Compatible API | `datastore.engine` | `datastore.version` | DB instance type                                | Flavor prefix                  |
|----------------|--------------------|---------------------|-------------------------------------------------|--------------------------------|
| Cassandra      | `cassandra`        | `3.11`              | Cluster (classic storage)                       | `geminidb.cassandra.`          |
| InfluxDB       | `influxdb`         | `1.7`               | Cluster with cloud native storage (performance-enhanced) | `geminidb.influxdb-geminifs.` |

If the `datastore` block is omitted, a GeminiDB Cassandra instance is created. The DB instance type is derived
from the engine, as each compatible API is offered with a single storage type.

Use the `opentelekomcloud_gemini_datastores_v3` and `opentelekomcloud_gemini_flavors_v3` data sources to look up
the versions and flavors that are actually available in your region.

## Example Usage

### Create a basic instance with tags
```hcl
resource "opentelekomcloud_gemini_instance_v3" "instance_1" {
  name              = "gaussdb_cassandra_instance_1"
  password          = var.password
  flavor            = "geminidb.cassandra.xlarge.8"
  volume_size       = 100
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create an instance with backup strategy
```hcl
resource "opentelekomcloud_gemini_instance_v3" "instance_1" {
  name              = "gaussdb_cassandra_instance_1"
  password          = var.password
  flavor            = "geminidb.cassandra.xlarge.8"
  volume_size       = 100
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }
}
```

### Create a GeminiDB Influx instance
```hcl
resource "opentelekomcloud_gemini_instance_v3" "influx_instance" {
  name              = "gemini_influx_instance_1"
  password          = var.password
  flavor            = "geminidb.influxdb-geminifs.xlarge.4"
  volume_size       = 100
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  datastore {
    engine         = "influxdb"
    version        = "1.7"
    storage_engine = "rocksDB"
  }
}
```

### Create a GeminiDB Cassandra instance in an enterprise project
```hcl
resource "opentelekomcloud_gemini_instance_v3" "instance_1" {
  name              = "gaussdb_cassandra_instance_1"
  password          = var.password
  flavor            = "geminidb.cassandra.xlarge.8"
  volume_size       = 100
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  enterprise_project_id = var.enterprise_project_id

  datastore {
    engine         = "cassandra"
    version        = "3.11"
    storage_engine = "rocksDB"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) Specifies the AZ name. For a three-AZ deployment instance,
  use commas (,) to separate the AZs, for example, `eu-de-01,eu-de-02,eu-de-03`.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name. The value
  must be 4 to 64 characters in length and start with a letter. It is case-sensitive and can contain only letters,
  digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String) Specifies the instance specifications, for example
  `geminidb.cassandra.xlarge.8` or `geminidb.influxdb-geminifs.xlarge.4`. The flavor must match the engine
  selected in the `datastore` block. The available values can be queried with the
  `opentelekomcloud_gemini_flavors_v3` data source, for details see
  [DB Instance Specifications](https://docs.otc.t-systems.com/geminidb/api-ref/apis_v3.1/api_versions_and_specifications/querying_instance_specifications.html).

* `node_num` - (Optional, Int) Specifies the number of nodes, ranges from 3 to 200. Defaults to 3.

  -> Nodes can only be removed from a GeminiDB Cassandra instance. Decreasing `node_num` of a GeminiDB Influx
  instance is rejected by the service.

* `volume_size` - (Required, Int) Specifies the storage space in GB. The value must be a multiple of 10. For both
  GeminiDB Cassandra and GeminiDB Influx DB instances, the minimum storage space is 100 GB, and the maximum storage
  space is related to the instance performance specifications.

* `password` - (Required, String) Specifies the database password. The value must be 8 to 32 characters in length,
  including uppercase and lowercase letters, digits, and special characters, such as ~!@#%^*-_=+? You are advised to
  enter a strong password to improve security, preventing security risks such as brute force cracking.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet. Changing this parameter will create a
  new resource.

* `security_group_id` - (Optional, String) Specifies the security group ID. Required if the selected subnet doesn't
  enable network ACL.

* `configuration_id` - (Optional, String) Specifies the Parameter Template ID.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID the instance belongs
  to. Leave it empty if the Enterprise Project Management Service is not enabled, otherwise the default enterprise
  project is used when the value is not set. Changing this parameter will create a new resource.

* `ssl` - (Optional, Bool, ForceNew) Specifies whether to enable or disable SSL. Defaults to `false`. Changing this
  parameter will create a new resource.

* `datastore` - (Optional, List, ForceNew) Specifies the database information. Structure is documented below. Changing
  this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below.

* `tags` - (Optional, Map) The key/value pairs to associate with the instance.

The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the compatible API. The value can be `cassandra` for a
  GeminiDB Cassandra instance or `influxdb` for a GeminiDB Influx instance. Defaults to `cassandra` when the
  `datastore` block is omitted. Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the database version. Use `3.11` for GeminiDB Cassandra and
  `1.7` for GeminiDB Influx. The versions actually available in your region can be queried with the
  `opentelekomcloud_gemini_datastores_v3` data source. Defaults to `3.11` when the `datastore` block is omitted.
  Changing this parameter will create a new resource.

* `storage_engine` - (Required, String, ForceNew) Specifies the storage engine. Only `rocksDB` is supported now,
  for both compatible APIs. Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the "hh:mm-HH:MM" format. The current time is in the UTC format. The
  HH value must be 1 greater than the hh value. The values of mm and MM must be the same and must be set to 00. Example
  value: 08:00-09:00, 03:00-04:00.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files. The value ranges from
  0 to 35. If this parameter is set to 0, the automated backup policy is not set. If this parameter is not transferred,
  the automated backup policy is enabled by default. Backup files are stored for seven days by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.

* `region` - Indicates the DB instance region.

* `status` - Indicates the DB instance status.

* `port` - Indicates the database port.

* `mode` - Indicates the DB instance type. `Cluster` for GeminiDB Cassandra and `CloudNativeCluster` for
  GeminiDB Influx.

* `db_user_name` - Indicates the default username.

* `private_ips` - Indicates the IP address list of the db.

* `nodes` - Indicates the instance nodes information. Structure is documented below.

The `nodes` block contains:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `status` - Indicates the node status.

* `support_reduce` - Indicates whether the node support reduce or not.

* `private_ip` - Indicates the private IP address of a node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 120 minutes.
* `delete` - Default is 30 minutes.

## Import

GeminiDB instance can be imported using the `id`, e.g.
```
$ terraform import opentelekomcloud_gemini_instance_v3.instance_1 749112383d5342e9acb6c7825801b452in06
```

Due to the security reasons, `password` can not be imported. It can be ignored as shown below.
```hcl
resource "opentelekomcloud_gemini_instance_v3" "instance_1" {
  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
