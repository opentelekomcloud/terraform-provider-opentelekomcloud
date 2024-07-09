---
subcategory: "Virtual Private Network (VPN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpnaas_service_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vpnaas-service-v2"
description: |-
Get details about a specific VPN from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPN service you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-network/api-ref/native_openstack_apis/vpn_service_management/querying_vpn_services.html#en-topic-0093011500)

# opentelekomcloud_vpnaas_service_v2

Use this data source to get details about a specific VPN.

## Example Usage

```hcl
variable "vpn_name" {}

data "opentelekomcloud_vpnaas_service_v2" "vpn" {
  name           = var.vpn_name
  admin_state_up = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain details about the V2 VPN service.

* `name` - (Optional) The name of the service.

* `tenant_id` - (Optional) The owner of the service.

* `description` - (Optional) The human-readable description for the service.

* `admin_state_up` - (Optional) The administrative state of the resource. Can either be true (Up) or false (Down).
  Default is `false`.

* `subnet_id` - (Optional) SubnetID is the ID of the subnet. Default is `null`.

* `router_id` - (Optional) The ID of the router. Default is `null`.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.

* `name` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `router_id` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `subnet_id` - See Argument Reference above.

* `status` - Indicates whether IPsec VPN service is currently operational. Values are `ACTIVE`,
  `DOWN`, `BUILD`, `ERROR`, `PENDING_CREATE`, `PENDING_UPDATE` or `PENDING_DELETE`.

* `external_v6_ip` - The read-only external (public) IPv6 address that is used for the VPN service.

* `external_v4_ip` - The read-only external (public) IPv4 address that is used for the VPN service.

* `description` - See Argument Reference above.

* `value_specs` - See Argument Reference above.
