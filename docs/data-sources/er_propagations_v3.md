---
subcategory: "Enterprise Router (ER)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_er_propagations_v3"
sidebar_current: "docs-opentelekomcloud-datasource-er-propagations-v3"
description: |-
  Get the list of ER propagations from OpenTelekomCloud
---
# opentelekomcloud_er_propagations_v3

Use this data source to get the list of propagations.

## Example Usage

```hcl
variable instance_id {}
variable route_table_id {}
variable attachment_id {}

data "opentelekomcloud_er_propagations_v3" "test" {
  instance_id    = var.instance_id
  route_table_id = var.route_table_id
  attachment_id  = var.attachment_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the ER instance ID to which the propagation belongs.

* `route_table_id` - (Required, String) Specifies the route table ID to which the propagation belongs.

* `attachment_id` - (Optional, String) Specifies the attachment ID to which the propagation belongs.

* `attachment_type` - (Optional, String) Specifies the attachment type of corresponding to the propagation.
  The valid values are as follows:
  + **vpc**: Virtual private cloud.
  + **vpn**: VPN gateway.
  + **vgw**: Virtual gateway of cloud private line.
  + **peering**: Peering connection, through the cloud connection (CC) to load ERs in different regions to create a
    peering connection.
  + **enc**: Enterprise connect network in EC.
  + **cfw**: VPC border firewall.

* `status` - (Optional, String) Specifies the status of the propagation. Default value is `available`.
  The valid values are as follows:
  + **available**
  + **failed**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `propagations` - All propagations that match the filter parameters.
  The [propagations](#route_propagations) structure is documented below.

<a name="route_propagations"></a>
The `propagations` block supports:

* `id` - The propagation ID.

* `instance_id` - The ER instance ID to which the propagation belongs.

* `route_table_id` - The route table ID of corresponding to the propagation.

* `attachment_id` - The attachment ID corresponding to the propagation.

* `attachment_type` - The attachment type corresponding to the propagation.

* `resource_id` - The resource ID of the attachment associated with the propagation.

* `status` - The current status of the propagation.

* `created_at` - The creation time of the propagation.

* `updated_at` - The latest update time of the propagation.
