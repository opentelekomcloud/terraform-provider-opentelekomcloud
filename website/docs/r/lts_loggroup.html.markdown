---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_logtank_group_v2"
sidebar_current: "docs-opentelekomcloud-resource-logtank-group-v2"
description: |-
  log group management
---

# opentelekomcloud\_logtank\_group\_v2

Manages a log group resource within OpenTelekomCloud.

## Example Usage

### create a log group

```hcl
resource "opentelekomcloud_logtank_group_v2" "log_group1" {
  group_name  = "log_group1"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required)
  Specifies the log group name.
  Changing this parameter will create a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The log group ID.

* `group_name` - See Argument Reference above.

* `ttl_in_days` -
  Specifies the log expiration time. The value is fixed to 7 days.

## Import

Log group can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_logtank_group_v2.group_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
