---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_port_v2"
sidebar_current: "docs-opentelekomcloud-datasource-networking-port-v2"
description: |-
Get network port ID from OpenTelekomCloud
---

Up-to-date reference of API arguments for Network port you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/port/querying_ports.html#vpc-port02-0001)

# opentelekomcloud_networking_port_v2

Use this data source to get the ID of an available OpenTelekomCloud port.

## Example Usage

```hcl
data "opentelekomcloud_networking_port_v2" "port_1" {
  name = "port_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve port ids. If omitted, the
  `region` argument of the provider is used.

* `project_id` - (Optional) The owner of the port.

* `port_id` - (Optional) The ID of the port.

* `name` - (Optional) The name of the port.

* `admin_state_up` - (Optional) The administrative state of the port.

* `network_id` - (Optional) The ID of the network the port belongs to.

* `device_owner` - (Optional) The device owner of the port.

* `mac_address` - (Optional) The MAC address of the port.

* `device_id` - (Optional) The ID of the device the port belongs to.

* `fixed_ip` - (Optional) The port IP address filter.

* `status` - (Optional) The status of the port.

* `security_group_ids` - (Optional) The list of port security group IDs to filter.

## Attributes Reference

`id` is set to the ID of the found port. In addition, the following attributes are exported:

* `region` - See Argument Reference above.

* `project_id` - See Argument Reference above.

* `port_id` - See Argument Reference above.

* `name` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `network_id` - See Argument Reference above.

* `device_owner` - See Argument Reference above.

* `mac_address` - See Argument Reference above.

* `device_id` - See Argument Reference above.

* `all_fixed_ips` - The collection of Fixed IP addresses on the port in the
  order returned by the Network v2 API.

* `all_security_group_ids` - The set of security group IDs applied on the port.
