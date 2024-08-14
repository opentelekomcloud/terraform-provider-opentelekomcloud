---
subcategory: "Enterprise Router (ER)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_er_association_v3"
sidebar_current: "docs-opentelekomcloud-resource-er-association-v3"
description: |-
  Manages an Enterprise Router Association resource within OpenTelekomCloud.
---

# opentelekomcloud_er_association_v3

Manages an association resource under the route table for ER service within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "route_table_id" {}
variable "attachment_id" {}

resource "opentelekomcloud_er_association_v3" "test" {
  instance_id    = var.instance_id
  route_table_id = var.route_table_id
  attachment_id  = var.attachment_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the ER instance to which the route table and the
  attachment belongs.

* `route_table_id` - (Required, String, ForceNew) Specifies the ID of the route table to which the association
  belongs.

* `attachment_id` - (Required, String, ForceNew) Specifies the ID of the attachment corresponding to the association.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attachment_type` - The type of the attachment corresponding to the association.

* `status` - The current status of the association.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

* `region` - The region where the ER instance and route table are located.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 2 minutes.

## Import

Associations can be imported using their `id` and the related `instance_id` and `route_table_id`, separated by
slashes (/), e.g.

```
$ terraform import opentelekomcloud_er_association_v3.test instance_id/route_table_id/id
```
