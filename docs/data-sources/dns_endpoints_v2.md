---
subcategory: "Domain Name Service (DNS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dns_endpoints_v2"
sidebar_current: "docs-opentelekomcloud-datasource-dns-endpoints-v2"
description: |-
  Get the list of DNS resolver endpoints within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DNS endpoint query you can get at
[documentation portal](https://docs.otc.t-systems.com/domain-name-service/api-ref/apis/endpoint_management/querying_endpoints.html)

# opentelekomcloud_dns_endpoints_v2

Use this data source to get the list of DNS resolver endpoints.

## Example Usage

```hcl
data "opentelekomcloud_dns_endpoints_v2" "outbound" {
  direction = "outbound"
  vpc_id    = var.vpc_id
}

output "outbound_endpoint_ids" {
  value = data.opentelekomcloud_dns_endpoints_v2.outbound.endpoints[*].id
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, String) The endpoint direction. Valid values are `inbound` and `outbound`.

* `vpc_id` - (Optional, String) The ID of the VPC to which the endpoints belong.

* `name` - (Optional, String) The endpoint name.

* `region` - (Optional, String) The region in which to query the endpoints.
  If omitted, the `region` argument of the provider is used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `endpoints` - The list of endpoints. Each element contains:
  * `id` - The endpoint ID.
  * `name` - The endpoint name.
  * `direction` - The endpoint direction.
  * `status` - The endpoint status.
  * `vpc_id` - The ID of the VPC to which the endpoint belongs.
  * `ip_address_count` - The number of IP addresses of the endpoint.
  * `resolver_rule_count` - The number of resolver rules of the endpoint. Returned only for outbound endpoints.
  * `created_at` - The creation time.
  * `updated_at` - The last update time.
