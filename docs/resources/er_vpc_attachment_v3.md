---
subcategory: "Enterprise Router (ER)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_er_vpc_attachment_v3"
sidebar_current: "docs-opentelekomcloud-resource-er-vpc-attachment-v3"
description: |-
  Manages an Enterprise Router VPC Attachment resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Enterprise Router you can get at
[documentation portal](https://docs.otc.t-systems.com/enterprise-router/api-ref/apis/vpc_attachments/index.html).

# opentelekomcloud_er_vpc_attachment_v3

Manages an ER vpc attachment resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "attachment_name" {}

resource "opentelekomcloud_er_vpc_attachment_v3" "test" {
  instance_id            = var.instance_id
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  name                   = var.attachment_name
  description            = "VPC attachment created by terraform"
  auto_create_vpc_routes = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the ER instance to which the VPC attachment
  belongs.
  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the VPC attachment belongs.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the VPC subnet to which the VPC attachment belongs.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the VPC attachment.
  The name can contain 1 to 64 characters, digits, underscore (_), hyphens (-) and
  dots (.) allowed.

* `description` - (Optional, String) Specifies the description of the VPC attachment.
  The description contain a maximum of 255 characters, and the angle brackets (< and >) are not allowed.

* `auto_create_vpc_routes` - (Optional, Bool, ForceNew) Specifies whether to automatically configure routes for the VPC
  which pointing to the ER instance.
  The destination CIDRs of the routes are fixed as follows:
    + **10.0.0.0/8**
    + **172.16.0.0/12**
    + **192.168.0.0/16**

  The default value is false. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The current status of the VPC attachment.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

* `region` - The region where the ER instance and the VPC attachment are.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 2 minutes.

## Import

VPC attachments can be imported using their `id` and the related `instance_id`, e.g.

```
$ terraform import opentelekomcloud_er_vpc_attachment_v3.test instance_id/id
```
