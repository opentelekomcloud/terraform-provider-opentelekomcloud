---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_metric_data_v1"
sidebar_current: "docs-opentelekomcloud-resource-ces-metric-data-v1"
description: |-
  Manages a v1 CES Metric Data resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES Metric Data v1 you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_description/monitoring_data/index.html)

# opentelekomcloud_ces_metric_data_v1

Manages a V1 CES Metric Data (used for reporting custom monitoring data) resource within OpenTelekomCloud.

## Example Usage

### Basic metric data usage

```hcl
variable instance_id {}

resource "opentelekomcloud_ces_metric_data_v1" "metric_1" {
  metric {
    namespace   = "TEST.TF_ACC"
    metric_name = "cpu_util"
    dimensions {
      name  = "instance_id"
      value = var.instance_id
    }
  }
  ttl          = 172800
  collect_time = 1257894000000
  value        = 0.09
  unit         = "%"
  type         = "float"
}
```

## Argument Reference

The following arguments are supported:

* `metric` - (Required, List, ForceNew) Specifies the metric data. The structure is described below.

* `ttl` - (Required, Integer, ForceNew) Specifies the data validity period. The unit is `second`. Supported range: `1` to `604800`. If the validity period expires, the data will be automatically deleted.

* `collect_time` - (Required, Integer, ForceNew) Specifies when the data was collected, which is a UNIX timestamp (ms). 
  
  -> **NOTE:**
  Since there is a latency between the client and the server, the data timestamp to be inserted should be within the period that starts from three days before the current time plus 20s to 10 minutes after the current time minus 20s. In this way, the timestamp will be inserted to the database without being affected by the latency.

* `value` - (Required, Float, ForceNew) Specifies the monitoring metric data to be added, which can be an integer or a floating point number.

* `unit` - (Optional, String, ForceNew) Specifies the data unit. Supported range: `1` to `32`

* `type` - (Optional, String, ForceNew) Specifies the enumerated type. Possible values: `int`, `float`


The `metric` block supports:

* `namespace` - (Required, String, ForceNew) Specifies the namespace in `service.item` format. `service.item` can be a string of `3` to `32` characters that must start with a letter and can consists of uppercase letters, lowercase letters, numbers, or underscores (_). In addition, service **cannot** start with `SYS`, `AGT`, or `SRE`, and namespace **cannot** be `SERVICE.BMS` because this namespace has been used by the system.

* `metric_name` - (Required, String, ForceNew) Specifies the metric ID. 
  [Available metrics.](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/tree/devel/opentelekomcloud/services/ces/interconnected_services.md)

* `dimensions` - (Required, List, ForceNew) Specifies the metric dimension. A maximum of three dimensions are supported.

The `dimensions` block supports:

* `name` - (Required, String, ForceNew) Specifies the dimension name. The value can be a string
  of `1` to `32` characters that must start with a letter and can consists of uppercase
  letters, lowercase letters, numbers, underscores (_), or hyphens (-).

* `value` - (Required, String, ForceNew) Specifies the dimension value. The value can be a string
  of `1` to `64` characters that must start with a letter or a number and can consists
  of uppercase letters, lowercase letters, numbers, underscores (_), or hyphens (-).
