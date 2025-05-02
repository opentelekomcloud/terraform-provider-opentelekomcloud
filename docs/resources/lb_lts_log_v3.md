---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_lts_log_v3"
sidebar_current: "docs-opentelekomcloud-resource-lb-lts-log-v3"
description: |-
  Manages a LB Access Log resource within OpenTelekomCloud.
---

# opentelekomcloud_lb_lts_log_v3

Manage a LB Access Log resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "loadbalancer_id" {}
variable "group_id" {}
variable "stream_id" {}

resource "opentelekomcloud_lb_lts_log_v3" "test" {
  loadbalancer_id = var.loadbalancer_id
  log_group_id    = var.group_id
  log_stream_id   = var.stream_id
}
```

## Argument Reference

The following arguments are supported:

* `loadbalancer_id` - (Required, String, ForceNew) Specifies the ID of a load balancer.

  Changing this creates a new resource.

* `log_group_id` - (Required, String) Specifies the ID of a log group.

* `log_stream_id` - (Required, String) Specifies the ID of the subscribe stream.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log ID.

* `region` - The region where resource created.

## Import

ELB Access Log resource can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_lb_lts_log_v3.log 7b80e108-1636-44e5-aece-986b0052b7dd
```
