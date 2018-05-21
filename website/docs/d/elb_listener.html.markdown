---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_elb_listener"
sidebar_current: "docs-opentelekomcloud-resource-elb-listener"
description: |-
  Manages a ELB listener resource within OpenTelekomCloud.
---

# opentelekomcloud\_elb\_listener

Manages an ELB listener resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "TCP"
  protocol_port = 8080
  backend_protocol = "TCP"
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
    timeouts {
	create = "5m"
	update = "5m"
	delete = "5m"
    }
}
```

## Argument Reference

The following arguments are supported:

* `protocol` - (Required) The protocol - can either be TCP, HTTP, HTTPS or TERMINATED_HTTPS.
    Changing this creates a new Listener.

* `protocol_port` - (Required) The port on which to listen for client traffic.
    Changing this creates a new Listener.

* `loadbalancer_id` - (Required) The load balancer on which to provision this
    Listener. Changing this creates a new Listener.

* `name` - (Optional) Human-readable name for the Listener. Does not have
    to be unique.

* `backend_port` - port

* `lb_algorithm` - 

## Attributes Reference

The following attributes are exported:

* `protocol` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `name` - See Argument Reference above.
