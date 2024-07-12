---
subcategory: "Bare Metal Server (BMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_bms_tags_v2"
sidebar_current: "docs-opentelekomcloud-resource-compute-bms-tags-v2"
description: |-
  Manages a BMS Tags resource within OpenTelekomCloud.
---

# opentelekomcloud_compute_bms_tags_v2

Used to add tags to a BMS within OpenTelekomCloud.

## Example Usage

```hcl
variable "bms_id" {}

resource "opentelekomcloud_compute_bms_tags_v2" "add_tags" {
  server_id = var.bms_id
  tags      = ["tags_type_baremetal"]
}
```

## Argument Reference

The following arguments are supported:

* `server_id`- (Required) The unique id of bare metal server.

* `tags` - (Required) The tags of a BMS. Changing this parameter creates a new resource.

## Attributes Reference

All above argument parameters can be exported as attribute parameters.

## Import

BMS tags can be imported using the server_id, e.g.

```
terraform import opentelekomcloud_compute_bms_tags_v2.add_tags 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
