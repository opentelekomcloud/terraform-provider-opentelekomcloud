---
subcategory: "Private NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_private_nat_snat_rule_v3"
sidebar_current: "docs-opentelekomcloud-data-source-private-nat-snat-rule-v3"
description: |-
  Manages a Private NAT SNAT rule data source within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Private NAT SNAT rule you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/apis_for_private_nat_gateways_v3.0/snat_rules/index.html)

# opentelekomcloud_private_nat_snat_rule_v3

Manages a V3 Private NAT SNAT rule data source within OpenTelekomCloud.

## Example Usage

### List all Private NAT SNAT rules

```hcl
data "opentelekomcloud_private_nat_snat_rule_v3" "rule_1" {}
```

### Get Private NAT SNAT rule using ID

```hcl
variable "id" {}

data source "opentelekomcloud_private_nat_snat_rule_v3" "rule_1" {
  id = var.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) Specifies the private NAT SNAT rule ID.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `snat_rules` - Indicates the private NAT SNAT rules. The structure is defined below.

The `snat_rules` block supports:

* `id` - Private NAT SNAT rule ID.

* `gateway_id` - Indicates the private NAT gateway ID.

* `cidr` - Indicates the CIDR block that matches the SNAT rule.

* `virsubnet_id` - Indicates the ID of the subnet that matches the SNAT rule.

* `description` - Provides supplementary information about the SNAT rule.

* `transit_ip_ids` - Indicates the IDs of the transit IP addresses.

* `project_id` - Indicates the project ID.

* `transit_ip_associations` - Indicates the list of details of associated transit IP addresses. The structure is described below.
    * `transit_ip_id` - Indicates the ID of the transit IP address.
    * `transit_ip_address` - Indicates the transit IP address.

* `created_at` - Indicates the time when the private NAT SNAT rule was created. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `updated_at` - Indicates the time when the private NAT SNAT rule was updated. It is a UTC time in yyyy-mm-ddThh:mm:ssZ format.

* `enterprise_project_id` - Indicates the ID of the enterprise project that is associated with the SNAT rule when the SNAT rule is created.

* `status` - Indicates the private NAT SNAT rule status. The value can be: `ACTIVE` (The SNAT rule is running properly) or `FROZEN` (The SNAT rule is frozen).
