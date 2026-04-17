---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_metrics_v1"
sidebar_current: "docs-opentelekomcloud-datasource-ces-metrics-v1"
description: |-
  Get details about CES metrics within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES metrics v1 you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/metrics/index.html)

# opentelekomcloud_ces_metrics_v1

Get details about CES metrics within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_ces_metrics_v1" "metrics_1" {
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional, String) Specifies the namespace of a service.

* `metric_name` - (Optional, String) Specifies the metric ID.

* `dim` - (Optional, String) Specifies the dimension. A maximum of three dimensions are supported, and the dimensions are numbered from `0` in `dim.{i}=key,value` format. key cannot exceed `32` characters and value cannot exceed 256 characters.
  Single dimension: `dim.0=instance_id,xxxx-xxxx-xxxx`,
  Multiple dimensions: `dim.0=key,value&dim.1=key,value`

* `start` - (Optional, String) Specifies the paging start value. The format is `namespace.metric_name.key:value`.
  Example: start=SYS.ECS.cpu_util.instance_id:d9112af5-6913-4f3b-bd0a-3f96711e004d.

* `limit` - (Optional, Int) Specifies the maximum number of records that can be queried at a time. Supported range: `1` to `1000` (default)

* `order` - (Optional, String) Specifies the result sorting method, which is sorted by timestamp. The default method is `desc`. Supported values:
  `asc`: The query results are displayed in the ascending order.
  `desc`: The query results are displayed in the descending order.

  -> **NOTE:**
  `namespaces` and `dimensions` are available on our [github link](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/tree/devel/opentelekomcloud/services/ces/interconnected_services.md) or [official documentation](https://docs.otc.t-systems.com/cloud-eye/api-ref/appendix/services_interconnected_with_cloud_eye.html).

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `metrics` - Specifies the list of metric objects. The structure is described below. 
* `meta_data` - Specifies the metadata of query results, including the pagination information. The structure is described below.
    * `count` - Specifies the number of returned results.
    * `marker` - Specifies the pagination marker.
    * `total` - Specifies the total number of metrics.


The `metrics` block supports:

* `namespace` - Specifies the metric namespace.
* `dimensions` - Specifies the list of metric dimensions. The structure is described below.
    * `name` - Specifies the dimension name.
    * `value` - Specifies the dimension value.
* `metric_name` - Specifies the metric name, such as cpu_util.
* `unit` - Specifies the metric unit.
