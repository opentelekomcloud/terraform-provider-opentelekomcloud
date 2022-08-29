---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_flavors_v3

Use this data source to get list of ELBv3 flavors names.

## Example Usage

```hcl
data "opentelekomcloud_lb_flavors_v3" "flavors_names" {}
```

## Argument Reference

* `id` - (Optional) Specifies the flavor ID.

* `name` - (Optional) Specifies the flavor name.

## Attributes Reference

In addition, the following attributes are exported:

* `flavors` - A list of all the flavors names found. This data source will fail if none are found.
