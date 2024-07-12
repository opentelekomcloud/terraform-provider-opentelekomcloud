---
subcategory: "Bare Metal Server (BMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_bms_nic_v2"
sidebar_current: "docs-opentelekomcloud-datasource-compute-bms-nic-v2"
description: |-
  Get details about BMS NIC from OpenTelekomCloud
---

Up-to-date reference of API arguments for BMS NIC you can get at
[documentation portal](https://docs.otc.t-systems.com/bare-metal-server/api-ref/native_openstack_nova_v2.1_apis/bms_nic_management/querying_information_about_bms_nics_native_openstack_api.html#en-us-topic-0053158678)

# opentelekomcloud_compute_bms_nic_v2

Use this data source to get details about a BMS NIC based on the NIC ID from OpenTelekomCloud.

## Example Usage

```hcl
variable "bms_id" {}
variable "nic_id" {}

data "opentelekomcloud_compute_bms_nic_v2" "query_bms_nic" {
  server_id = var.bms_id
  id        = var.nic_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `server_id` - (Required) This is the unique BMS id.

* `id` - (Optional) The ID of the NIC.

* `status` - (Optional) The NIC port status.

## Attributes Reference

All of the argument attributes are also exported as result attributes.

* `mac_address` - It is NIC's mac address.

* `fixed_ips` - The NIC IP address.

* `network_id` - The ID of the network to which the NIC port belongs.
