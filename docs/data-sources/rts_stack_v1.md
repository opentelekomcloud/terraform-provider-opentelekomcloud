---
subcategory: "Resource Template Service (RTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rts_stack_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rts-stack-v1"
description: |-
Get details about a specific RTS resource from OpenTelekomCloud
---

# opentelekomcloud_rts_stack_v1

Use this data source to get details about RTS outputs and other useful data including the template body.

## Example Usage

```hcl
variable "stack_name" {}

data "opentelekomcloud_rts_stack_v1" "mystack" {
  name = var.stack_name
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
