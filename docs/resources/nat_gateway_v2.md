---
subcategory: "NAT"
---

Up-to-date reference of API arguments for NAT gateway you can get at
`https://docs.otc.t-systems.com/nat-gateway/api-ref/api_v2.0/nat_gateway_service`.

# opentelekomcloud_nat_gateway_v2

Manages a V2 NAT Gateway resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "router_id" {}
variable "internal_network_id" {}

resource "opentelekomcloud_nat_gateway_v2" "this" {
  name                = "tf_nat"
  description         = "NAT GW created by terraform"
  spec                = "0"
  router_id           = var.router_id
  internal_network_id = var.internal_network_id

  tags = {
    muh = "kuh"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NAT Gateway.

* `description` - (Optional) The description of the NAT Gateway.

* `spec` - (Required) The specification of the NAT Gateway, valid values are `"0"`,`"1"`, `"2"`, `"3"`, `"4"`.

* `tenant_id` - (Optional) The target tenant ID in which to allocate the NAT
  Gateway. Changing this creates a new NAT Gateway.

* `router_id` - (Required) ID of the router (or VPC) this NAT Gateway belongs to. Changing
  this creates a new NAT Gateway.

* `internal_network_id` - (Required) ID of the network this NAT Gateway connects to.
  Changing this creates a new NAT Gateway.

* `tags` - (Optional) Tags key/value pairs to associate with the NAT Gateway.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `spec` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `router_id` - See Argument Reference above.

* `internal_network_id` - See Argument Reference above.

Gateway can be imported using the following format:

```sh
terraform import opentelekomcloud_nat_gateway_v2.gw_1 e4f783a7-b908-4215-b018-724960e5g34t
```
