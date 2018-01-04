---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_elb_health"
sidebar_current: "docs-opentelekomcloud-resource-elb-health"
description: |-
  Manages an ELB health resource within OpenTelekomCloud.
---

# opentelekomcloud\_elb\_health

Manages an ELB health resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_elb_health" "health_1" {
  listener_id = "}"
  healthy_threshold = 3

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Optional)

## Attributes Reference

The following attributes are exported:

* `listener_id` - 
