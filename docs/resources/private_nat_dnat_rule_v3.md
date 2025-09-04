---
subcategory: "Private NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_private_nat_dnat_rule_v3"
sidebar_current: "docs-opentelekomcloud-resource-private-nat-dnat-rule-v3"
description: |-
  Manages a Private NAT DNAT rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Private NAT DNAT rule you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/apis_for_private_nat_gateways_v3.0/dnat_rules/index.html)

# opentelekomcloud_private_nat_dnat_rule_v3

Manages a V3 Private NAT DNAT rule resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "transit_ip_id" {}
variable "network_interface_id" {}
variable "gateway_id" {}

resource "opentelekomcloud_private_nat_dnat_rule_v3" "rule_1" {
  description          = "created"
  transit_ip_id        = var.transit_ip_id
  network_interface_id = var.network_interface_id
  gateway_id           = var.gateway_id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, String) Provides supplementary information about the DNAT rule. The description can contain up to 255 characters and cannot contain angle brackets (<>).

* `transit_ip_id` - (Required, String) Specifies the ID of the transit IP address.

* `network_interface_id` - (Optional, String) Specifies the port ID of the resource that the NAT gateway is bound to. The resource can be a compute instance, load balancer (v2 or v3), or virtual IP address. Either this parameter or `private_ip_address` must be specified. Otherwise, an error will be reported.

* `gateway_id` - (Required, String, ForceNew) Specifies the private NAT gateway ID.

* `protocol` - (Optional, String) Specifies the protocol type. Supported values: `tcp`, `udp`, `any`.

* `private_ip_address` - (Optional, String) Specifies the port IP address that the NAT gateway uses. The resource can be a compute instance, load balancer (v2 or v3), or virtual IP address. Either this parameter or `network_interface_id` must be specified. Otherwise, an error will be reported.

* `internal_service_port` - (Optional, String) Specifies the port number of the resource, which can be a compute instance, load balancer (v2 or v3), or virtual IP address. Value range: `0-65535`. Default value: `0`.

* `transit_service_port` - (Optional, String) Specifies the port number of the transit IP address. Value range: `0-65535`. Default value: `0`.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Private NAT DNAT rule ID.

* `project_id` - Indicates the project ID.

* `type` - Indicates the backend resource type of the DNAT rule. The type can be: `COMPUTE`, `VIP`, `ELB`, `ELBv3`, `CUSTOMIZE`.

* `enterprise_project_id` - Indicates the ID of the enterprise project that is associated with the DNAT rule when the DNAT rule is created.

* `created_at` - Indicates the time when the private NAT DNAT rule was created. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `updated_at` - Indicates the time when the private NAT DNAT rule was updated. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `status` - Indicates the private NAT DNAT rule status. The value can be: `ACTIVE` (The DNAT rule is running properly) or `FROZEN` (The DNAT rule is frozen).


## Import

Private NAT DNAT rule V3 resource can be imported using the DNAT rule ID, `id`.

```sh
terraform import opentelekomcloud_private_nat_dnat_rule_v3.rule_1 <id>
```
