---
subcategory: "FunctionGraph"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fgs_event_v2"
sidebar_current: "docs-opentelekomcloud-resource-fgs-event-v2"
description: |-
Manages an FGS Event resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for FGS you can get at
[documentation portal](https://docs.otc.t-systems.com/function-graph/api-ref/apis/index.html)

# opentelekomcloud_fgs_event_v2

Manages an event for testing specified function within OpenTelekomCloud.

## Example Usage

### Create a simple event

```hcl
variable "function_urn" {}
variable "event_name" {}
variable "event_content" {}

resource "opentelekomcloud_fgs_event_v2" "test" {
  function_urn = var.function_urn
  name         = var.event_name
  content      = base64encode(var.event_content)
}
```

## Argument Reference

* `function_urn` - (Required, String, ForceNew) Specifies the URN of the function to which the event belongs.

* `name` - (Required, String) Specifies the function event name.
  The name can contain a maximum of `25`, only letters, digits, underscores (_) and hyphens (-) are allowed.

* `content` - (Required, String) Specifies the function event content.
  The value is the base64 encoding of the JSON string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `updated_at` - The latest update time of the function event.

* `region` - The region in which function graph resource is created.

## Import

Function event can be imported using the `function_urn` and `id`, separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_fgs_event_v2.test <function_urn>/<id>
```
