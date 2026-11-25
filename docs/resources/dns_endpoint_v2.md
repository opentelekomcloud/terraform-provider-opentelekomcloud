---
subcategory: "Domain Name Service (DNS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dns_endpoint_v2"
sidebar_current: "docs-opentelekomcloud-resource-dns-endpoint-v2"
description: |-
  Manages a DNS resolver endpoint resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DNS endpoint you can get at
[documentation portal](https://docs.otc.t-systems.com/domain-name-service/api-ref/apis/endpoint_management)

# opentelekomcloud_dns_endpoint_v2

Manages a DNS resolver endpoint resource within OpenTelekomCloud.

An inbound endpoint lets networks outside the VPC send DNS queries to the cloud DNS.
An outbound endpoint forwards DNS queries from the VPC to external DNS servers according to endpoint rules.

## Example Usage

```hcl
resource "opentelekomcloud_dns_endpoint_v2" "inbound" {
  name      = "inbound-endpoint"
  direction = "inbound"

  ip_addresses {
    subnet_id = var.subnet_id
  }

  ip_addresses {
    subnet_id = var.subnet_id
    ip        = "192.168.0.100"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The endpoint name. It can contain `1` to `64` characters.
  Only letters, digits, underscores (`_`), hyphens (`-`) and periods (`.`) are allowed.

* `direction` - (Required, String, ForceNew) The endpoint direction. Valid values are `inbound` and `outbound`.
  Changing this creates a new endpoint.

* `ip_addresses` - (Required, List) The IP addresses of the endpoint. At least `2` and at most `6` addresses
  can be specified. The [ip_addresses](#ip_addresses) structure is documented below.

* `region` - (Optional, String, ForceNew) The region in which to create the endpoint.
  If omitted, the `region` argument of the provider is used. Changing this creates a new endpoint.

<a name="ip_addresses"></a>
The `ip_addresses` block supports:

* `subnet_id` - (Required, String) The ID of the VPC subnet (`network_id` of `opentelekomcloud_vpc_subnet_v1`).

* `ip` - (Optional, String) The IP address. It must belong to the subnet.
  If omitted, an address is assigned automatically.

-> Changes of `ip_addresses` bind new and unbind removed addresses one at a time.
  To remove a particular address from the middle of the list, set `ip` explicitly for every element.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The endpoint ID.

* `vpc_id` - The ID of the VPC to which the endpoint belongs.

* `status` - The endpoint status.

* `created_at` - The creation time.

* `updated_at` - The last update time.

* `ip_addresses/ip_address_id` - The ID of the IP address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

DNS endpoints can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_dns_endpoint_v2.inbound ff80808169bf16c70169bf1d02270000
```
