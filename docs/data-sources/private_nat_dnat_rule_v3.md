---
subcategory: "Private NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_private_nat_dnat_rule_v3"
sidebar_current: "docs-opentelekomcloud-data-source-private-nat-dnat-rule-v3"
description: |-
  Manages a Private NAT DNAT rule v3 data source within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Private NAT DNAT rule you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/apis_for_private_nat_gateways_v3.0/dnat_rules/index.html)

# opentelekomcloud_private_nat_dnat_rule_v3

Manages a V3 Private NAT DNAT rule v3 data source within OpenTelekomCloud.

## Example Usage

### List all Private NAT DNAT rules

```hcl
data "opentelekomcloud_private_nat_dnat_rule_v3" "rule_1" {}
```

### Get Private NAT DNAT rule using ID

```hcl
variable "id" {}

data source "opentelekomcloud_private_nat_dnat_rule_v3" "rule_1" {
  id = var.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - Specifies the private NAT DNAT rule ID.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Private NAT DNAT rule ID.

* `description` - Provides supplementary information about the DNAT rule.

* `transit_ip_id` - Indicates the ID of the transit IP address.

* `network_interface_id` - Indicates the port ID of the resource that the NAT gateway is bound to. The resource can be a compute instance, load balancer (v2 or v3), or virtual IP address.

* `gateway_id` - Indicates the private NAT gateway ID.

* `protocol` - Indicates the protocol type. Supported values: `tcp`, `udp`, `any`.

* `private_ip_address` - Indicates the port IP address that the NAT gateway uses.

* `internal_service_port` - Indicates the port number of the resource, which can be a compute instance, load balancer (v2 or v3), or virtual IP address.

* `transit_service_port` - Indicates the port number of the transit IP address.

* `project_id` - Indicates the project ID.

* `type` - Indicates the backend resource type of the DNAT rule. The type can be: `COMPUTE`, `VIP`, `ELB`, `ELBv3`, `CUSTOMIZE`.

* `enterprise_project_id` - Indicates the ID of the enterprise project that is associated with the DNAT rule when the DNAT rule is created.

* `created_at` - Indicates the time when the private NAT DNAT rule was created. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `updated_at` - Indicates the time when the private NAT DNAT rule was updated. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `status` - Indicates the private NAT DNAT rule status. The value can be: `ACTIVE` (The DNAT rule is running properly) or `FROZEN` (The DNAT rule is frozen).
