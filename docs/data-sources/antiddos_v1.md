---
subcategory: "Anti-DDoS"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_antiddos_v1"
sidebar_current: "docs-opentelekomcloud-datasource-antiddos-v1"
description: |-
  Get a status of EIP from OpenTelekomCloud
---

Up-to-date reference of API arguments for EIP status you can get at
[documentation portal](https://docs.otc.t-systems.com/anti-ddos/api-ref/api/anti-ddos_apis/querying_the_list_of_defense_statuses_of_eips.html#antiddos-02-0023)

# opentelekomcloud_antiddos_v1

Use this data source to query the status of EIP, regardless whether an EIP has been bound to an Elastic Cloud Server (ECS) or not.

## Example Usage

```hcl
variable "id" {}

data "opentelekomcloud_antiddos_v1" "antiddos" {
  floating_ip_id = var.eip_id
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip_id` - (Optional) The Elastic IP ID.

* `floating_ip_address` - (Optional) The Elastic IP address.

* `status` - (Optional) The defense status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_type` - The EIP type.

* `period_start` - The Start time.

* `bps_attack` - The Attack traffic in (bit/s).

* `bps_in` - The inbound traffic in (bit/s).

* `total_bps` - The total traffic.

* `pps_in` - The inbound packet rate (number of packets per second).

* `pps_attack` - The attack packet rate (number of packets per second).

* `total_pps` - The total packet rate.

* `start_time` - The start time of cleaning and blackhole event.

* `end_time` - The end time of cleaning and blackhole event.

* `traffic_cleaning_status` - The traffic cleaning status.

* `trigger_bps` - The traffic at the triggering point.

* `trigger_pps` - The packet rate at the triggering point.

* `trigger_http_pps` - The HTTP request rate at the triggering point.
