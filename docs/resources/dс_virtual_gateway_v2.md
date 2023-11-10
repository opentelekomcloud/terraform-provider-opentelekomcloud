---
subcategory: "Direct Connect (DCaaS)"
---
# opentelekomcloud_dc_virtual_gateway_v2 (Resource)

Up-to-date reference of API arguments for Direct Connect Virtual Gateway you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/direct-connect/api-ref/apis/virtual_gateway/index.html).

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `vpc_id` (String, Required, ForceNew) - Specifies the ID of the VPC to be accessed.
* `local_ep_group_id` (String, Required) - Specifies the ID of the local endpoint group that records CIDR blocks of the VPC subnets.
* `name` (String, Required) - Specifies the virtual gateway name.
* `description` (String, Optional) - Provides supplementary information about the virtual gateway.
* `asn` (Int, Optional, ForceNew) - Specifies the BGP ASN of the virtual gateway.
* `device_id` (String, Optional) - Specifies the ID of the physical device used by the virtual gateway.
* `redundant_device_id` (String, Optional) - Specifies the ID of the redundant physical device used by the virtual gateway.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the virtual gateway.
* `status` -  Virtual gateway status.
* `project_id` -  Project id.

## Import

Direct Connect Virtual Gateway can be imported using `id`, e.g.

```sh
$ terraform import opentelekomcloud_dc_virtual_gateway_v2.vgw <id>
```
