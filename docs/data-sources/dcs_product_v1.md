---
subcategory: "Distributed Cache Service (DCS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dcs_product_v1"
sidebar_current: "docs-opentelekomcloud-datasource-dcs-product-v1"
description: |-
  Get DCS product from OpenTelekomCloud
---

Up-to-date reference of API arguments for DCS product you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-cache-service/api-ref/outdated_apis_v1/other_apis/querying_service_specifications.html#dcs-api-0312040)

# opentelekomcloud_dcs_product_v1

Use this data source to get the ID of an available DCS product.

## Example Usage

```hcl
data "opentelekomcloud_dcs_product_v1" "product1" {
  spec_code = "dcs.single_node"
}
```

## Argument Reference

* `spec_code` - (Optional) Indicates an I/O specification.

## Attributes Reference

`id` is set to the ID of the found product. In addition, the following attributes are exported:

* `spec_code` - See Argument Reference above.
