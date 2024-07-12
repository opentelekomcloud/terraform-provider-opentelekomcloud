---
subcategory: "Resource Template Service (RTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rts_stack_resource_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rts-stack-resource-v1"
description: |-
  Get details about a specific RTS resource metadata from OpenTelekomCloud
---

# opentelekomcloud_rts_stack_resource_v1

Use this data source to get details about RTS resource metadata.

## Example Usage

```hcl
variable "stack_name" {}
variable "resource_name" {}

data "opentelekomcloud_rts_stack_resource_v1" "stackresource" {
  stack_name    = var.stack_name
  resource_name = var.resource_name
}
```

## Argument Reference

The following arguments are supported:

* `stack_name` - (Required) The unique stack name.

* `resource_name` - (Optional) The name of a resource in the stack.

* `physical_resource_id` - (Optional) The physical resource ID.

* `resource_type` - (Optional) The resource type.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `logical_resource_id` - The logical resource ID.

* `resource_status` - The status of the resource.

* `resource_status_reason` - The resource operation reason.

* `required_by` - Specifies the resource dependency.
