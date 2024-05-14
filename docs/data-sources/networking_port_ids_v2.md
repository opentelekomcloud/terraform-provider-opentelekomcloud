---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_networking_port_ids_v2

Use this data source to get a list of OpenTelekomCloud Port IDs matching the
specified criteria.

## Example Usage

```hcl
data "opentelekomcloud_networking_port_ids_v2" "ports" {
  name = "port"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve port ids. If omitted, the
  `region` argument of the provider is used.

* `project_id` - (Optional, String) The owner of the port.

* `name` - (Optional, String) The name of the port.

* `admin_state_up` - (Optional, Bool) The administrative state of the port.

* `network_id` - (Optional, String) The ID of the network the port belongs to.

* `device_owner` - (Optional, String) The device owner of the port.

* `mac_address` - (Optional, String) The MAC address of the port.

* `device_id` - (Optional, String) The ID of the device the port belongs to.

* `fixed_ip` - (Optional, String) The port IP address filter.

* `status` - (Optional, String) The status of the port.

* `security_group_ids` - (Optional, List) The list of port security group IDs to filter.

* `sort_key` - (Optional) Sort ports based on a certain key. Defaults to none.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.
  Defaults to none.

## Attributes Reference

`ids` is set to the list of OpenTelekomCloud Port IDs.
