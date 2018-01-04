---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_elb_backend"
sidebar_current: "docs-opentelekomcloud-resource-elb-backend"
description: |-
  Manages an ELB backend resource within OpenTelekomCloud.
---

# opentelekomcloud\_elb\_backend

Manages an Elastic Load Balancer Backend resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_elb_backend" "backend_1" {
  address = "1.1.1.1s"
  listener_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
  server_id = "%s"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
```

## Argument Reference

The following arguments are supported:

    Note:  One of LoadbalancerID or ListenerID must be provided.

* `listener_id` - (Optional) The Listener on which the members of the pool
    will be associated with. Changing this creates a new pool.
	Note:  One of LoadbalancerID or ListenerID must be provided.


## Attributes Reference

The following attributes are exported:

