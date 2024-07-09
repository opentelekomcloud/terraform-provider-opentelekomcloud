---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_flavors_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lb-flavors-v3"
description: |-
Get ELBv3 flavors names from OpenTelekomCloud
---

Up-to-date reference of API arguments for ELBv3 flavor you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/load_balancer_flavor/querying_flavors.html#listflavors)

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
