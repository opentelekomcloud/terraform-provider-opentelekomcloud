---
subcategory: "Direct Connect (DCaaS)"
---
# opentelekomcloud_dc_virtual_interface_v2 (Resource)

Up-to-date reference of API arguments for Direct Connect Virtual Interface you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/direct-connect/api-ref/apis/virtual_interface/index.html).

## Example Usage

```hcl
variable direct_connect_id {}

data "opentelekomcloud_identity_project_v3" "project" {
  name = "eu-de_project_1"
}

resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group" {
  name        = "tf_acc_eg_1"
  type        = "cidr"
  endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
  description = "first"
  project_id  = data.opentelekomcloud_identity_project_v3.project.id
}

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name              = "my_virtual_gateway"
  description       = "acc test"
  local_ep_group_id = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group.id
}

resource "opentelekomcloud_dc_virtual_interface_v2" "int_1" {
  direct_connect_id = var.direct_connect_id
  vgw_id            = opentelekomcloud_dc_virtual_gateway_v2.vgw_1.id
  name              = "vi_1"
  description       = "description"
  type              = "private"
  route_mode        = "static"
  vlan              = 100
  bandwidth         = 5

  remote_ep_group_id   = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group.id
  local_gateway_v4_ip  = "180.1.1.1/24"
  remote_gateway_v4_ip = "180.1.1.2/24"
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_id` (String, Required, ForceNew) - Specifies the connection ID.
* `virtual_gateway_id` (String, Required, ForceNew) - Specifies the virtual gateway ID.
* `name` (String, Required) - Specifies the virtual interface name.
* `type` (String, Required, ForceNew) - Specifies the virtual interface type. The value can only be `private`.
* `route_mode` (String, Required, ForceNew) - Specifies the routing mode. The value can be `static` or `bgp`.
* `vlan` (Int, Required, ForceNew) - Specifies the VLAN used by the local gateway to communicate with the remote gateway.
* `bandwidth` (Int, Required) - Specifies the virtual interface bandwidth.
* `remote_ep_group_id` (String, Required) - Specifies the ID of the remote endpoint group that records the CIDR blocks used by the on-premises network.
* `description` (String, Optional) - Provides supplementary information about the virtual interface.
* `service_type` (String, Optional, ForceNew) - Specifies what is to be accessed over the connection. The value can only be `vpc`.
* `local_gateway_v4_ip` (String, Optional, ForceNew) - Specifies the IPv4 address of the local gateway.
* `remote_gateway_v4_ip` (String, Optional, ForceNew) - Specifies the IPv4 address of the remote gateway.
* `asn` (Int, Optional, ForceNew) - Specifies the AS number of the BGP peer.
* `bgp_md5` (String, Optional, ForceNew) - Specifies the MD5 password of the BGP peer.
* `project_id` (String, Optional, ForceNew) - Specifies the project ID.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the virtual interface.
* `enable_bfd` - Bidirectional Forwarding Detection (BFD) function status.
* `enable_nqa` -  Network Quality Analysis (NQA) function status.
* `lag_id` -  The ID of the link aggregation group (LAG) associated with the virtual interface.
* `status` -  The current status of the virtual interface.
* `created_at` -  The creation time of the virtual interface.

## Import

Direct Connect Virtual Interface can be imported using `id`, e.g.

```sh
$ terraform import opentelekomcloud_dc_virtual_interface_v2.vgw <id>
```
