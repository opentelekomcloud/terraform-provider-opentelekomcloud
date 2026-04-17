---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_multiple_metric_data_v1"
sidebar_current: "docs-opentelekomcloud-datasource-ces-multiple-metric-data-v1"
description: |-
  Get details about CES metric data in a batch within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES multiple metric data v1 you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/monitoring_data/index.html)

# opentelekomcloud_ces_multiple_metric_data_v1

Get details about CES metric data in a batch within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_ces_multiple_metric_data_v1" "metric_data_1" {
  metrics {
    namespace   = "TEST.TF_ACC"
    metric_name = "cpu_util"
    dimensions {
      name  = "instance_id"
      value = "72d1377e-09e4-47bd-8ea4-71a815d4815d"
    }
  }
  from   = 1257892000000
  to     = 1257895000000
  period = "1"
  filter = "average"
}
```

## Argument Reference

The following arguments are supported:

* `metrics` - (Required, List) Specifies the metric data. The maximum length of the array is 10. The structure is documented below.

* `from` - (Required, Integer) Specifies the start time of the query. The time is a UNIX timestamp and the unit is ms. Rollup aggregates the raw data generated within a period to the start time of the period. If `from` and `to` are within a period, the query result will be empty due to the rollup failure. Set `from` to at least one period earlier than the current time. Take the 5-minute period as an example. If it is 10:35 now, the raw data generated between 10:30 and 10:35 will be aggregated to 10:30. In this example, if period is 5 minutes, from should be 10:30.

* `to` - (Required, Integer) Specifies the end time of the query. The time is a UNIX timestamp and the unit is ms. `from` must be earlier than `to`.

* `period` - (Required, String) Specifies how often Cloud Eye aggregates data. Possible values are:
    * `1`: Cloud Eye performs no aggregation and displays raw data.
    * `300`: Cloud Eye aggregates data every 5 minutes.
    * `1200`: Cloud Eye aggregates data every 20 minutes.
    * `3600`: Cloud Eye aggregates data every hour.
    * `14400`: Cloud Eye aggregates data every 4 hours.
    * `86400`: Cloud Eye aggregates data every 24 hours.

* `filter` - (Required, String) Specifies the data rollup method. Possible values are:
    * `average`: Cloud Eye calculates the average value of metric data within a rollup period.
    * `max`: Cloud Eye calculates the maximum value of metric data within a rollup period.
    * `min`: Cloud Eye calculates the minimum value of metric data within a rollup period.
    * `sum`: Cloud Eye calculates the sum of metric data within a rollup period.
    * `variance`: Cloud Eye calculates the variance value of metric data within a rollup period.

  -> **NOTE:**
  Rollup uses a rollup method to aggregate raw data generated within a specific period. Take the 5-minute period as an example. If it is 10:35 now, the raw data generated between 10:30 and 10:35 will be aggregated to 10:30. `filter` does not affect the query result of raw data (when `period` is `1`).

The `metrics` block supports:

* `namespace` - (Required, String) Specifies the namespace of a service.
* `metric_name` - (Required, String) Specifies the metric name.
* `dimensions` - (Required, List) Specifies the list of metric dimensions. Maximum number of supported dimensions is `4`. The structure is described below.
    * `name` - (Required, String) Specifies the dimension name.
    * `value` - (Required, String) Specifies the dimension value.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `metrics` - Specifies the metric data. The structure is documented below.

The `metrics` block supports:
* `namespace` - See Argument Reference above.
* `dimensions` - See Argument Reference above.
* `metric_name` - See Argument Reference above.
* `unit` - Specifies the metric unit.
* `datapoints` - Specifies the metric data list. Since Cloud Eye rounds up from based on the level of granularity for data query, datapoints may contain more data points than expected. The structure is described below. 
    * `average` - Specifies the average value of metric data within a rollup period.
    * `max` - Specifies the maximum value of metric data within a rollup period.
    * `min` - Specifies the minimum value of metric data within a rollup period.
    * `sum` - Specifies the sum of metric data within a rollup period.
    * `variance` - Specifies the variance of metric data within a rollup period.
    * `timestamp` - Specifies when the metric is collected. It is a UNIX timestamp in milliseconds.
