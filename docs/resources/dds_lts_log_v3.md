---
subcategory: "Document Database Service (DDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dds_lts_log_v3"
sidebar_current: "docs-opentelekomcloud-resource-dds-lts-log-v3"
description: |-
  Manages a DDS LTS Log resource within OpenTelekomCloud.
---

# opentelekomcloud_dds_lts_log_v3

Manages a DDS LTS log resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "lts_group_id" {}
variable "lts_stream_id" {}

resource "opentelekomcloud_dds_lts_log_v3" "log" {
  instance_id   = var.instance_id
  log_type      = "audit_log"
  lts_group_id  = var.lts_group_id
  lts_stream_id = var.lts_stream_id
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DDS instance.

  Changing this creates a new resource.

* `log_type` - (Required, String, ForceNew) Specifies the type of the LTS log. The value can be `audit_log`.

  Changing this creates a new resource.

* `lts_group_id` - (Required, String) Specifies the ID of the LTS log group.

* `lts_stream_id` - (Required, String) Specifies the ID of the LTS log stream.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is DDS instance ID.

* `region` - The region in which created the LTS log resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

The DDS LTS log config can be imported using DDS instance ID, e.g.

```bash
$ terraform import opentelekomcloud_dds_lts_log_v3.log <instance_id>
```
