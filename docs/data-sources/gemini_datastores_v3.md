---
subcategory: "GeminiDB"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gemini_datastores_v3"
sidebar_current: "docs-opentelekomcloud-datasource-gemini-datastores-v3"
description: |-
  Use this data source to get the GeminiDB datastore information within OpenTelekomCloud.
---

Up-to-date reference of API arguments for GeminiDB datastores you can get at
[documentation portal](https://docs.otc.t-systems.com/geminidb/api-ref/apis_v3/versions_and_specifications/querying_version_information.html).

# opentelekomcloud_gemini_datastores_v3

Use this data source to get the versions, storage engines and DB instance types supported by a GeminiDB
compatible API.

## Example Usage

### Query the GeminiDB Cassandra datastore
```hcl
data "opentelekomcloud_gemini_datastores_v3" "cassandra" {
  engine_name = "cassandra"
}
```

### Create an instance with the latest supported version
```hcl
data "opentelekomcloud_gemini_datastores_v3" "influx" {
  engine_name = "influxdb"
}

resource "opentelekomcloud_gemini_instance_v3" "influx" {
  name              = "gemini_influx_instance"
  password          = var.password
  flavor            = "geminidb.influxdb-geminifs.xlarge.4"
  volume_size       = 100
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  datastore {
    engine         = data.opentelekomcloud_gemini_datastores_v3.influx.engine_name
    version        = data.opentelekomcloud_gemini_datastores_v3.influx.versions[0]
    storage_engine = data.opentelekomcloud_gemini_datastores_v3.influx.storage_engines[0]
  }
}
```

## Argument Reference

* `engine_name` - (Required, String) Specifies the compatible API to query. The value can be `cassandra`
  (GeminiDB Cassandra) or `influxdb` (GeminiDB Influx).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `region` - Indicates the region in which the datastore was queried.

* `versions` - Indicates the list of database versions supported by the compatible API. Currently `3.11` for
  GeminiDB Cassandra and `1.7` for GeminiDB Influx.

  -> The version API doesn't report the GeminiDB Influx versions, it returns an empty list for `influxdb`.
  In that case the versions advertised by the specification API are returned instead.

* `storage_engines` - Indicates the list of storage engines supported by the compatible API. Currently
  `rocksDB` is the only storage engine offered for both GeminiDB Cassandra and GeminiDB Influx.

* `modes` - Indicates the list of DB instance types supported by the compatible API. GeminiDB Cassandra is
  offered as a classic storage cluster (`Cluster`), while GeminiDB Influx is offered as a performance-enhanced
  cluster with cloud native storage (`CloudNativeCluster`).

  -> `storage_engines` and `modes` are not exposed by the GeminiDB API, they are fixed per compatible API and
  reflect what the service currently offers.
