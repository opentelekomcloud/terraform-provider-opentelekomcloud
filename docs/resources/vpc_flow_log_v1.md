---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_flow_log_v1"
sidebar_current: "docs-opentelekomcloud-resource-vpc-flow-log-v1"
description: |-
  Manages a VPC Flow Log resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC flow log you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/vpc_flow_log)

# opentelekomcloud_vpc_flow_log_v1

Manages a VPC flow log resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_logtank_group_v2" "log_group1" {
  group_name = var.log_group_name
}

resource "opentelekomcloud_logtank_topic_v2" "log_topic1" {
  group_id   = opentelekomcloud_logtank_group_v2.log_group1.id
  topic_name = var.log_topic_name
}

resource "opentelekomcloud_vpc_v1" "vpc_v1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "opentelekomcloud_vpc_flow_log_v1" "flowlog1" {
  name          = var.flow_log_name
  description   = var.flow_log_desc
  resource_type = "vpc"
  resource_id   = opentelekomcloud_vpc_v1.vpc_v1.id
  traffic_type  = "all"
  log_group_id  = opentelekomcloud_logtank_group_v2.log_group1.id
  log_topic_id  = opentelekomcloud_logtank_topic_v2.log_topic1.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specifies the flow log name.
  The value is a string of 1 to 64 characters that can contain letters, digits, underscores (_), hyphens (-) and periods (.).

* `description` - (Optinal) Provides supplementary information about the VPC flow log.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).

* `resource_type` - (Required) Specifies the type of resource on which to create the VPC flow log.
  The value can be `port`, `vpc` and `network`.
  Changing this creates a new VPC flow log.

* `resource_id` - (Required) Specifies the ID of resource type.
  Changing this creates a new VPC flow log.

* `traffic_type` - (Required) Specifies the type of traffic to log. The value can be `all`, `accept` and `reject`.
  Changing this creates a new VPC flow log.

* `log_group_id` - (Required) Specifies the log group ID.
  Changing this creates a new VPC flow log.

* `log_topic_id` - (Required) Specifies the log topic ID.
  Changing this creates a new VPC flow log.

## Attributes Reference

The following attributes are exported:

* `id` - The VPC flow log ID in UUID format.

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `resource_type` - See Argument Reference above.

* `resource_id` - See Argument Reference above.

* `traffic_type` - See Argument Reference above.

* `log_group_id` - See Argument Reference above.

* `log_topic_id` - See Argument Reference above.

* `status` - The status of the flow log. The value can be `ACTIVE`, `DOWN` or `ERROR`.

* `admin_state` - Whether to enable the VPC flow log function.

## Import

VPC flow logs can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_flow log_v1 ab76d479-9ef8-4034-88c4-4ab82fc87572
```
