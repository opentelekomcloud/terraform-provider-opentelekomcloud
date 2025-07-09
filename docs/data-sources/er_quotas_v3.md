---
subcategory: "Enterprise Router (ER)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_er_quotas_v3"
sidebar_current: "docs-opentelekomcloud-datasource-er-quotas-v3"
description: |-
  Query Quotas for enterprise routers within OpenTelekomCloud
---
# opentelekomcloud_er_quotas_v3

Using this data source to query the list of available resource quotas within OpenTelekomCloud.

~> Using an invalid ID to filter the results will not report an error or return an empty list, but will return a quota
   list with all usage equal to 0.

## Example Usage

```hcl
data "opentelekomcloud_er_quotas_v3" "test" {}
```

## Argument Reference

The following arguments are supported:

* `type` - (Optional, String) The quota type to be queried.
  The valid values are as follows:
  + `er_instance`: Quotas and usage for enterprise router instances.
  + `dc_attachment`: Quotas and usage for DC attachment.
  + `vpc_attachment`: Quotas and usage for VPC attachment.
  + `vpn_attachment`: Quotas and usage for VPN attachment.
  + `peering_attachment`: Quotas and usage for peering attachment.
  + `can_attachment`: Quotas and usage for can attachment.
  + `route_table`: Quotas and usage for route table.
  + `static_route`: Quotas and usage for static route.
  + `vpc_er`: The number of enterprise routers that each VPC can access and the current usage.
  + `flow_log`: The number of flow logs that can be created per attachment.

* `instance_id` - (Optional, String) The instance ID.

* `route_table_id` - (Optional, String) The route table ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - All quotas that match the filter parameters.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `used` - The number of quota used.

* `unit` - The unit of usage.

* `type` - The quota type.

* `available` - The number of available quotas, `-1` means unlimited.

* `region` - The region where resources are located.

