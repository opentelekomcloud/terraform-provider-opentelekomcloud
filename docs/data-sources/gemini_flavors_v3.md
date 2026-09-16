---
subcategory: "GeminiDB"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gemini_flavors_v3"
sidebar_current: "docs-opentelekomcloud-datasource-gemini-flavors-v3"
description: |-
  Use this data source to get the list of GeminiDB flavors within OpenTelekomCloud.
---

Up-to-date reference of API arguments for GeminiDB flavors you can get at
[documentation portal](https://docs.otc.t-systems.com/geminidb/api-ref/apis_v3.1/api_versions_and_specifications/querying_instance_specifications.html).

# opentelekomcloud_gemini_flavors_v3

Use this data source to get available OpenTelekomCloud GeminiDB flavors (instance specifications).

## Example Usage

### List all GeminiDB Cassandra flavors
```hcl
data "opentelekomcloud_gemini_flavors_v3" "cassandra" {
  engine_name = "cassandra"
}
```

### Find a GeminiDB Influx flavor and use it to create an instance
```hcl
data "opentelekomcloud_gemini_flavors_v3" "influx" {
  engine_name = "influxdb"
  vcpus       = 4
}

resource "opentelekomcloud_gemini_instance_v3" "influx" {
  name              = "gemini_influx_instance"
  password          = var.password
  flavor            = data.opentelekomcloud_gemini_flavors_v3.influx.flavors[0].spec_code
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

### Filter flavors by availability zone
```hcl
data "opentelekomcloud_gemini_flavors_v3" "this" {
  engine_name       = "cassandra"
  availability_zone = "eu-de-01"
}
```

## Argument Reference

* `engine_name` - (Optional, String) Specifies the compatible API of the flavors to query. The value can be
  `cassandra` or `influxdb`. Defaults to `cassandra`.

* `mode` - (Optional, String) Specifies the DB instance type. The only supported value is `CloudNativeCluster`,
  which returns the specifications of instances with cloud native storage. If this parameter is not set, the
  specifications of all instances with classic storage are returned.

  -> GeminiDB Influx is only offered as a performance-enhanced cluster with cloud native storage. When
  `engine_name` is `influxdb` and `mode` is not set, `CloudNativeCluster` is used, so that its flavors are
  returned without having to set `mode` explicitly.

* `engine_version` - (Optional, String) Specifies the database version used to filter the flavors, for example
  `3.11` for GeminiDB Cassandra or `1.7` for GeminiDB Influx.

* `vcpus` - (Optional, Int) Specifies the number of vCPUs used to filter the flavors.

* `ram` - (Optional, Int) Specifies the memory size in GB used to filter the flavors.

* `availability_zone` - (Optional, String) Specifies the AZ name used to filter the flavors. Only the flavors
  that are on sale in this AZ are returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `region` - Indicates the region in which the flavors were queried.

* `flavors` - An array of available flavors. Only the flavors that are on sale in at least one AZ are returned.
  Structure is documented below.

The `flavors` block contains:

* `spec_code` - Indicates the resource specification code, for example `geminidb.cassandra.xlarge.8`. Use this
  value for the `flavor` argument of `opentelekomcloud_gemini_instance_v3`.

* `engine_name` - Indicates the compatible API of the flavor.

* `engine_version` - Indicates the database version of the flavor.

* `vcpus` - Indicates the number of vCPUs.

* `ram` - Indicates the memory size in GB.

* `availability_zones` - Indicates the list of AZs where the flavor is on sale.

* `az_status` - Indicates the status of the flavor per AZ. The value can be `normal`, `unsupported`
  or `sellout`.
