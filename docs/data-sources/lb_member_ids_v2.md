---
subcategory: "Elastic Load Balancer (ELB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_member_ids_v2"
sidebar_current: "docs-opentelekomcloud-datasource-lb-member-ids-v2"
description: |-
Get details about ELBv2 pool members from OpenTelekomCloud
---

Up-to-date reference of API arguments for ELBv3 pool members you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v2.0/backend_server/querying_backend_servers.html#elb-zq-hd-0002)

# opentelekomcloud_lb_member_ids_v2

Use this data source to get a list of member IDs for a ELBv2 pool from OpenTelekomCloud.
This data source can be useful for getting back a list of member IDs for a ELBv2 pool.

## Example Usage

```hcl
data "opentelekomcloud_lb_member_ids_v2" "this" {
  pool_id = var.pool_id
}
```

## Argument Reference

The following arguments are supported:

* `pool_id` - (Required) Specifies the ELBv2 pool ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of all the member IDs found.

