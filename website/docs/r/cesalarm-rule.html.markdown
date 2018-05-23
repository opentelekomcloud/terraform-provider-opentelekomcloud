---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces-alarm-rule"
sidebar_current: "docs-opentelekomcloud-resource-ces-alarm-rule"
description: |-
  Manages a V2 topic resource within OpenTelekomCloud.
---

# opentelekomcloud\_ces\_alarm\_rule

Manages a V2 topic resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_ces_alarmrule" "alarm_rule" {
  "alarm_name" = "alarm_rule"
  "metric" {
    "namespace" = "SYS.ECS"
    "metric_name" = "network_outgoing_bytes_rate_inband"
    "dimensions" {
        "name" = "instance_id"
        "value" = "${opentelekomcloud_compute_instance_v2.webserver.id}"
    }
  }
  "condition"  {
    "period" = 300
    "filter" = "average"
    "comparison_operator" = ">"
    "value" = 6
    "unit" = "B/s"
    "count" = 1
  }
  "alarm_actions" {
    "type" = "notification"
    "notification_list" = [
      "${opentelekomcloud_smn_topic_v2.topic.id}"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_name` - (Required) The name of the alarm rule to be created.

* `matric` - (Optional) 

* `condition` - (Optional) 

* `alarm_actions` - (Optional) 


## Attributes Reference

The following attributes are exported:

* `alarm_name` - See Argument Reference above.
* `matric` - See Argument Reference above.
* `condition` - See Argument Reference above.
* `alarm_actions` - See Argument Reference above.
