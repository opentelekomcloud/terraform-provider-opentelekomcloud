---
subcategory: "Domain Name Service (DNS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dns_resolver_rules_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dns-resolver-rules-v2"
description: |-
  Get the list of DNS resolver (endpoint) rules within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DNS endpoint rule query you can get at
[documentation portal](https://docs.otc.t-systems.com/domain-name-service/api-ref/apis/endpoint_rule_management/querying_endpoint_rules.html)

# opentelekomcloud_dns_resolver_rules_v2

Use this data source to get the list of DNS resolver (endpoint) rules.

## Example Usage

```hcl
data "opentelekomcloud_dns_resolver_rules_v2" "rules" {
  endpoint_id = var.endpoint_id
}

output "forwarded_domains" {
  value = data.opentelekomcloud_dns_resolver_rules_v2.rules.resolver_rules[*].domain_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) The rule name.

* `domain_name` - (Optional, String) The domain name of the rule.

* `endpoint_id` - (Optional, String) The ID of the endpoint to which the rules belong.

* `region` - (Optional, String) The region in which to query the rules.
  If omitted, the `region` argument of the provider is used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resolver_rules` - The list of rules. Each element contains:
  * `id` - The rule ID.
  * `name` - The rule name.
  * `domain_name` - The domain name, with a trailing dot.
  * `endpoint_id` - The ID of the endpoint to which the rule belongs.
  * `status` - The rule status.
  * `ip_address_count` - The number of DNS server IP addresses of the rule.
  * `vpc_ids` - The IDs of the VPCs associated with the rule.
  * `created_at` - The creation time.
  * `updated_at` - The last update time.
