---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rts_stack_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rts-stack-v1"
description: |-
  Provides metadata of an RTS stack (e.g. outputs).
---

# Data Source: opentelekomcloud_rts_stack_v1

The OpenTelekomCloud RTS Stack data source allows access to stack outputs and other useful data including the template body.

## Example Usage


```hcl

data "opentelekomcloud_rts_stack_v1" "mystack" {
  name = "rts-stack"
}

resource "opentelekomcloud_compute_instance_v2" "basic" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"

  network {
    uuid = "${data.opentelekomcloud_rts_stack_v1.mystack.outputs["SubnetId"]}"
  }
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) The name of the stack.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id` - A unique identifier of the stack.

* `capabilities` - List of stack capabilities for stack.

* `notification_topics` - List of notification topics for stack.

* `status` - Specifies the stack status.

* `disable_rollback` - Whether the rollback of the stack is disabled when stack creation fails.

* `outputs` - A list of stack outputs.

* `parameters` - A map of parameters that specify input parameters for the stack.

* `template_body` - Structure containing the template body.

* `timeout_mins` - Specifies the timeout duration.
