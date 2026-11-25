---
subcategory: "Domain Name Service (DNS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dns_resolver_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-dns-resolver-rule-v2"
description: |-
  Manages a DNS resolver (endpoint) rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DNS endpoint rule you can get at
[documentation portal](https://docs.otc.t-systems.com/domain-name-service/api-ref/apis/endpoint_rule_management)

# opentelekomcloud_dns_resolver_rule_v2

Manages a DNS resolver (endpoint) rule resource within OpenTelekomCloud.

A resolver rule forwards DNS queries for a domain name from the associated VPCs
to the specified DNS servers through an outbound endpoint.

## Example Usage

```hcl
resource "opentelekomcloud_dns_endpoint_v2" "outbound" {
  name      = "outbound-endpoint"
  direction = "outbound"

  ip_addresses {
    subnet_id = var.subnet_id
  }

  ip_addresses {
    subnet_id = var.subnet_id
  }
}

resource "opentelekomcloud_dns_resolver_rule_v2" "rule" {
  name         = "forward-example"
  domain_name  = "example.com"
  endpoint_id  = opentelekomcloud_dns_endpoint_v2.outbound.id
  ip_addresses = ["10.0.0.53", "10.0.1.53"]
  vpc_ids      = [var.vpc_id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The rule name. It can contain `1` to `64` characters.
  Only letters, digits, underscores (`_`), hyphens (`-`) and periods (`.`) are allowed.

* `domain_name` - (Required, String, ForceNew) The domain name to forward queries for.
  The API returns it with a trailing dot. Changing this creates a new rule.

* `endpoint_id` - (Required, String, ForceNew) The ID of the outbound endpoint to which the rule belongs.
  Changing this creates a new rule.

* `ip_addresses` - (Required, Set of Strings) The IP addresses of the DNS servers to which queries are forwarded.

* `vpc_ids` - (Optional, Set of Strings) The IDs of the VPCs associated with the rule.

* `region` - (Optional, String, ForceNew) The region in which to create the rule.
  If omitted, the `region` argument of the provider is used. Changing this creates a new rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID.

* `status` - The rule status.

* `created_at` - The creation time.

* `updated_at` - The last update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

DNS resolver rules can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_dns_resolver_rule_v2.rule 8a36f60a753badb401753bade3400002
```
