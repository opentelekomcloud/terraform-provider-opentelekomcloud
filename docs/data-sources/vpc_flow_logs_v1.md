---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_flow_logs_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpc-flow-logs-v1"
description: |-
  Queries VPC flow logs in OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC flow logs is available at the
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/vpc_apis_v1_v2/vpc_flow_log/querying_vpc_flow_logs.html).

# opentelekomcloud_vpc_flow_logs_v1

Queries VPC flow logs using optional filters.

## Example Usage

```hcl
data "opentelekomcloud_vpc_flow_logs_v1" "flow_logs" {
  resource_type = "vpc"
  resource_id   = var.vpc_id
  status        = "ACTIVE"
}
```

## Argument Reference

* `id` - (Optional) Flow log ID.

* `name` - (Optional) Flow log name.

* `tenant_id` - (Optional) Project ID.

* `description` - (Optional) Flow log description.

* `resource_type` - (Optional) Resource type. Valid values are `port`, `vpc`, and `network`.

* `resource_id` - (Optional) ID of the resource whose traffic is logged.

* `traffic_type` - (Optional) Traffic type. Valid values are `all`, `accept`, and `reject`.

* `log_group_id` - (Optional) LTS log group ID.

* `log_topic_id` - (Optional) LTS log topic ID.

* `status` - (Optional) Flow log status. Valid values are `ACTIVE`, `DOWN`, and `ERROR`.

* `limit` - (Optional) Number of records requested per page. Valid values are from `0` to `2000`.
  Defaults to `2000`.

* `marker` - (Optional) Resource ID after which listing starts.

## Attributes Reference

* `flow_logs` - Matching VPC flow logs. Each entry contains:
  * `id` - Flow log ID.
  * `name` - Flow log name.
  * `tenant_id` - Project ID.
  * `description` - Flow log description.
  * `resource_type` - Resource type.
  * `resource_id` - Resource ID.
  * `traffic_type` - Traffic type.
  * `log_group_id` - LTS log group ID.
  * `log_topic_id` - LTS log topic ID.
  * `enabled` - Whether the flow log function is enabled.
  * `status` - Flow log status.
  * `created_at` - Creation time.
  * `updated_at` - Last update time.
