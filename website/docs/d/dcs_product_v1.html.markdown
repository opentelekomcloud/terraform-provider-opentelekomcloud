---
layout: "opentelekomcloud"
page_title: "Opentelekomcloud: opentelekomcloud_dcs_product_v1"
sidebar_current: "docs-opentelekomcloud-datasource-dcs-product-v1"
description: |-
  Get information on an Opentelekomcloud dcs product.
---

# opentelekomcloud\_dcs\_product_v1

Use this data source to get the ID of an available Opentelekomcloud dcs product.

## Example Usage

```hcl

data "opentelekomcloud_dcs_product_v1" "product1" {
  spec_code = "dcs.single_node"
}
```

## Argument Reference

* `spec_code` - (Optional) Indicates an I/O specification.

## Attributes Reference

`id` is set to the ID of the found product. In addition, the following attributes
are exported:

* `spec_code` - See Argument Reference above.
