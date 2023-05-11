---
subcategory: "Virtual Private Network (VPN)"
---

Up-to-date reference of API arguments for VPNAAS endpoint group service you can get at
`https://docs.otc.t-systems.com/virtual-private-network/api-ref/native_openstack_apis/vpn_endpoint_group_management`.

# opentelekomcloud_vpnaas_endpoint_group_v2

Manages a V2 Endpoint Group resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_vpnaas_endpoint_group_v2" "group_1" {
  name      = "Group 1"
  type      = "cidr"
  endpoints = ["10.2.0.0/24", "10.3.0.0/24"]

  lifecycle {
    create_before_destroy = true
  }
}
```

~>
  Endpoint group can't be deleted when used, `create_before_destroy` makes it possible to make
  changes which require endpoint group recreation.

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create an endpoint group. If omitted, the
  `region` argument of the provider is used. Changing this creates a new group.

* `name` - (Optional) The name of the group. Changing this updates the name of
  the existing group.

* `tenant_id` - (Optional) The owner of the group. Required if admin wants to
  create an endpoint group for another project. Changing this creates a new group.

* `description` - (Optional) The human-readable description for the group.
  Changing this updates the description of the existing group.

* `type` -  The type of the endpoints in the group. A valid value is subnet, cidr, network, router, or vlan.
  Changing this creates a new group.

* `endpoints` - List of endpoints of the same type, for the endpoint group. The values will depend on the type.
  Changing this creates a new group.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.

* `name` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `description` - See Argument Reference above.

* `type` - See Argument Reference above.

* `endpoints` - See Argument Reference above.

* `value_specs` - See Argument Reference above.

## Import

Groups can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpnaas_endpoint_group_v2.group_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
