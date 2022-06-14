---
subcategory: "Elastic Load Balancer (ELB)"
---

# opentelekomcloud_lb_member_ids_v2

Use this data source to get a list of member IDs for a ELBv2 pool from OpenTelekomCloud.
This data source can be useful for getting back a list of member IDs for a ELBv2 pool.

## Example Usage

```hcl
data "opentelekomcloud_lb_member_ids_v2" "this" {
  pool_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `pool_id` - (Required) Specifies the ELBv2 pool ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of all the member IDs found. This data source will fail if none are found.

