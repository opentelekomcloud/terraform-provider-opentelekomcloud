---
subcategory: "Bare Metal Server (BMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_bms_flavors_v2"
sidebar_current: "docs-opentelekomcloud-datasource-compute-bms-flavors-v2"
description: |-
  Get details about flavors of BMSs from OpenTelekomCloud
---

Up-to-date reference of API arguments for BMSs flavors you can get at
[documentation portal](https://docs.otc.t-systems.com/bare-metal-server/api-ref/native_openstack_nova_v2.1_apis/bms_flavor_query/querying_bms_flavors_native_openstack_api.html#en-us-topic-0053158684)

# opentelekomcloud_compute_bms_flavors_v2

Use this data source to get details about flavors of BMSs from OpenTelekomCloud.

## Example Usage

```hcl
variable "flavor_id" {}
variable "disk_size" {}

data "opentelekomcloud_compute_bms_flavors_v2" "query_bms_flavors" {
  id       = var.bms_id
  min_disk = var.disk_size
  sort_key = "id"
  sort_dir = "desc"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `name` - (Optional) The name of the BMS flavor.

* `id` - (Optional) The BMS flavor id.

* `min_ram` - (Optional) The minimum memory size in MB. Only the BMSs with the memory size greater than or equal to the minimum size can be queried.

* `min_disk` - (Optional) The minimum disk size in GB. Only the BMSs with a disk size greater than or equal to the minimum size can be queried.

* `sort_key` - (Optional) The sorting field. The default value is **flavorid**. The other values are **name**, **memory_mb**, **vcpus**, **root_gb**, or **flavorid**.

* `sort_dir` - (Optional) The sorting order, which can be **ascending** (**asc**) or **descending** (**desc**). The default value is **asc**.

## Attributes Reference

All of the argument attributes are also exported as result attributes.

* `ram` - It is the memory size (in MB) of the flavor.

* `vcpus` - It is the number of CPU cores in the BMS flavor.

* `disk` - Specifies the disk size (GB) in the BMS flavor.

* `swap` - This is a reserved attribute.

* `rx_tx_factor` - This is a reserved attribute.
