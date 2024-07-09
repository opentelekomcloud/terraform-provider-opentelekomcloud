---
subcategory: "Bare Metal Server (BMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_bms_server_v2"
sidebar_current: "docs-opentelekomcloud-datasource-compute-bms-server-v2"
description: |-
Get details about BMS from OpenTelekomCloud
---

Up-to-date reference of API arguments for BMS you can get at
[documentation portal](https://docs.otc.t-systems.com/bare-metal-server/api-ref/native_openstack_nova_v2.1_apis/bms_lifecycle_management/querying_details_about_bmss_native_openstack_api.html#en-us-topic-0053158679)

# opentelekomcloud_compute_bms_server_v2

Use this data source to get details about a BMS or BMSs from OpenTelekomCloud.

## Example Usage

```hcl
variable "bms_id" {}
variable "bms_name" {}

data "opentelekomcloud_compute_bms_server_v2" "query_bms" {
  id   = var.bms_id
  name = var.bms_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `id` - (Optional) The unique ID of the BMS.

* `user_id` - (Optional) The ID of the user to which the BMS belongs.

* `name` - (Optional) The name of BMS.

* `status` - (Optional) The BMS status.

* `host_status` - (Optional) The nova-compute status: `UP`, `UNKNOWN`, `DOWN`, `MAINTENANCE` and `Null`.

* `key_name` - (Optional) It is the SSH key name.

* `flavor_id` - (Optional) It gives the BMS flavor information.

* `image_id` - (Optional) The BMS image.


## Attributes Reference

All of the argument attributes are also exported as result attributes.

* `host_id` - It is the host ID of the BMS.

* `progress` - This is a reserved attribute.

* `metadata` -  The BMS metadata is specified.

* `access_ip_v4` -  This is a reserved attribute.

* `access_ip_v6` - This is a reserved attribute.

* `addresses` - It gives the BMS network address.

* `security_groups` - The list of security groups to which the BMS belongs.

* `tags` - Specifies the BMS tag.

* `locked` -  It specifies whether a BMS is locked, true: The BMS is locked, false: The BMS is not locked.

* `config_drive` -  This is a reserved attribute.

* `availability_zone` - Specifies the AZ ID.

* `description` -  Provides supplementary information about the pool.

* `kernel_id` - The UUID of the kernel image when the AMI image is used.

* `hypervisor_hostname` -  It is the name of a host on the hypervisor.

* `instance_name` - Instance name is specified.
