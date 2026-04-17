---
subcategory: "Private NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_private_nat_snat_rule_v3"
sidebar_current: "docs-opentelekomcloud-resource-private-nat-snat-rule-v3"
description: |-
  Manages a Private NAT SNAT rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Private NAT SNAT rule you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/apis_for_private_nat_gateways_v3.0/snat_rules/index.html)

# opentelekomcloud_private_nat_snat_rule_v3

Manages a V3 Private NAT SNAT rule resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "transit_ip_id" {}
variable "network_id" {}
variable "gateway_id" {}

resource "opentelekomcloud_private_nat_snat_rule_v3" "rule_1" {
  description    = "created"
  transit_ip_ids = ["var.transit_ip_id"]
  virsubnet_id   = var.network_id
  gateway_id     = var.gateway_id
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the private NAT gateway ID.

* `cidr` - (Optional, String, ForceNew) Specifies the CIDR block that matches the SNAT rule. Either this parameter or `virsubnet_id` must be specified.

* `virsubnet_id` - (Optional, String, ForceNew) Specifies the ID of the subnet that matches the SNAT rule. Either this parameter or `cidr` must be specified.

* `description` - (Optional, String) Provides supplementary information about the SNAT rule. The description can contain up to 255 characters and cannot contain angle brackets (<>).

* `transit_ip_ids` - (Required, List[String]) Specifies the IDs of the transit IP addresses. A maximum number of `20` IDs is allowed.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Private NAT SNAT rule ID.

* `project_id` - Indicates the project ID.

* `transit_ip_associations` - Indicates the list of details of associated transit IP addresses. The structure is described below.
    * `transit_ip_id` - Indicates the ID of the transit IP address.
    * `transit_ip_address` - Indicates the transit IP address.

* `created_at` - Indicates the time when the private NAT SNAT rule was created. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `updated_at` - Indicates the time when the private NAT SNAT rule was updated. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `enterprise_project_id` - Indicates the ID of the enterprise project that is associated with the SNAT rule when the SNAT rule is created.

* `status` - Indicates the private NAT SNAT rule status. The value can be: `ACTIVE` (The SNAT rule is running properly) or `FROZEN` (The SNAT rule is frozen).

## Import

Private NAT SNAT rule V3 resource can be imported using the SNAT rule ID, `id`.

```sh
terraform import opentelekomcloud_private_nat_snat_rule_v3.rule_1 <id>
```
