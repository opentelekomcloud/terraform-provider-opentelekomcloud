---
subcategory: "Direct Connect (DCaaS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dc_virtual_gateway_v2"
sidebar_current: "docs-opentelekomcloud-resource-dc-virtual-gateway-v2"
description: |-
  Manages a Direct Connect Virtual Gateway resource within OpenTelekomCloud.
---

# opentelekomcloud_dc_virtual_gateway_v2 (Resource)

Up-to-date reference of API arguments for Direct Connect Virtual Gateway you can get at
[documentation portal](https://docs.otc.t-systems.com/direct-connect/api-ref/apis/virtual_gateway/index.html).

## Example Usage

```hcl
resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name        = "%s"
  description = "acc test"
  local_ep_group {
    name        = "tf_eg_1"
    endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
    description = "first"
  }
  local_ep_group_v6 {
    name        = "tf_eg_2"
    endpoints   = ["2a07:8700:2:4::/64", "2a07:8700:2:54::/64"]
    description = "first"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` (String, Required, ForceNew) - Specifies the ID of the VPC to be accessed.

* `local_ep_group` (String, Optional, List) - Specifies the local endpoint group that records CIDR IPV4 blocks of the VPC subnets.
  At least one of `local_ep_group` or `local_ep_group_v6` should be specified.
  The [object](#local_group) structure is documented below.

* `local_ep_group_v6` (String, Optional, List) - Specifies the local endpoint group that records CIDR IPV6 blocks of the VPC subnets.
  At least one of `local_ep_group` or `local_ep_group_v6` should be specified.
  The [object](#local_group) structure is documented below.

* `name` (String, Required) - Specifies the virtual gateway name.]()

* `description` (String, Optional) - Provides supplementary information about the virtual gateway.

* `asn` (Int, Optional, ForceNew) - Specifies the BGP ASN of the virtual gateway.

* `device_id` (String, Optional) - Specifies the ID of the physical device used by the virtual gateway.

* `project_id` (String, Optional) - Specifies the project ID.

* `redundant_device_id` (String, Optional) - Specifies the ID of the redundant physical device used by the virtual gateway.

  <a name="local_group"></a>
  The `local_group` block supports:

    * `name` (String, Optional) - Specifies the name of the Direct Connect endpoint group.

    * `description` (Str[ing, Optional) - Provides supplementary information about the Direct Connect endpoint group.

    * `endpoints` (List, Required) - Specifies the list of the endpoints in a Direct Connect endpoint group.

    * `type` (String, Required, ForceNew) - Specifies the type of the Direct Connect endpoints. The value can only be `cidr`. Default value: `cidr`.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the virtual gateway.

* `status` -  Virtual gateway status.

* `local_ep_group_id` - ID of the local IPV4 endpoint group that records CIDR blocks of the VPC subnets.

* `local_ep_group_ipv6_id` - ID of the local IPV6 endpoint group that records CIDR blocks of the VPC subnets.

## Import

Direct Connect Virtual Gateway can be imported using `id`, e.g.

```sh
$ terraform import opentelekomcloud_dc_virtual_gateway_v2.vgw <id>
```
