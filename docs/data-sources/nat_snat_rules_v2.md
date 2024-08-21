---
subcategory: "NAT"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_nat_snat_rules_v2"
sidebar_current: "docs-opentelekomcloud-datasource-nat-snat-rules-v2"
description: |-
Get details about NAT Gateway SNAT rules resource from OpenTelekomCloud
---

Up-to-date reference of API arguments for NAT Gateway you can get at
[documentation portal](https://docs.otc.t-systems.com/nat-gateway/api-ref/api_v2.0/snat_rules/querying_snat_rules.html#nat-api-0007)

# opentelekomcloud_nat_snat_rules_v2

Use this data source to get the list of SNAT rules.

## Example Usage

```hcl
variable "rule_id" {}

data "opentelekomcloud_nat_snat_rules_v2" "rule" {
  rule_id = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the SNAT rules are located.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the ID of the SNAT rule.

* `gateway_id` - (Optional, String) Specifies the ID of the NAT gateway to which the SNAT rule belongs.

* `floating_ip_id` - (Optional, String) Specifies the ID of the EIP associated with SNAT rule.

* `floating_ip_address` - (Optional, String) Specifies the IP of the EIP associated with SNAT rule.

* `cidr` - (Optional, String) Specifies the CIDR block to which the SNAT rule belongs.

* `subnet_id` - (Optional, String) Specifies the ID of the subnet to which the SNAT rule belongs.

* `source_type` - (Optional, String) Specifies the source type of the SNAT rule.
  The value can be one of the following:
  + **0** : The use scenario is VPC.
  + **1** : The use scenario is DC.

* `status` - (Optional, String) Specifies the status of the SNAT rule.
  The value can be one of the following:
  + **ACTIVE**: The SNAT rule is available.
  + **EIP_FREEZED**: The global EIP is frozen associated with SNAT rule.
  + **INACTIVE**: The SNAT rule is unavailable.

* `global_eip_id` - (Optional, String) Specifies the ID of the global EIP associated with SNAT rule.

* `global_eip_address` - (Optional, String) Specifies the IP of the global EIP associated with SNAT rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the SNAT rules.
  The [rules](#nat_snat_rules) structure is documented below.

<a name="nat_snat_rules"></a>
The `rules` block supports:

* `id` - The ID of the SNAT rule.

* `gateway_id` - The ID of the NAT gateway to which the SNAT rule belongs.

* `cidr` - The CIDR block to which the SNAT rule belongs.

* `subnet_id` - The ID of the subnet to which the SNAT rule belongs.

* `source_type` - The source type of the SNAT rule.

* `floating_ip_id` - The IDs of the EIP associated with SNAT rule, multiple EIP IDs separate by commas.
  e.g. **ID1,ID2**.

* `floating_ip_address` - The IPs of the EIP associated with SNAT rule, multiple EIP IPs separate by commas.
  e.g. **IP1,IP2**.

* `description` - The description of the SNAT rule.

* `status` - The status of the SNAT rule.

* `global_eip_id` - The IDs of the global EIP associated with SNAT rule, multiple global EIP IDs separate by commas.
  e.g. **ID1,ID2**.

* `global_eip_address` - The IPs of the global EIP associated with SNAT rule, multiple global EIP IPs separate by commas.
  e.g. **IP1,IP2**.

* `freezed_ip_address` - The IP of the frozen global EIP associated with SNAT rule.

* `created_at` - The creation time of the SNAT rule.
