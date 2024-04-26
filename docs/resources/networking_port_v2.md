---
subcategory: "Virtual Private Cloud (VPC)"
---

Up-to-date reference of API arguments for VPC port you can get at
`https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/port`.

# opentelekomcloud_networking_port_v2

Manages a V2 port resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) A unique name for the port. Changing this
  updates the `name` of an existing port.

* `network_id` - (Required, String, ForceNew) The ID of the network to attach the port to. Changing
  this creates a new port.

* `admin_state_up` - (Optional, Bool) Administrative up/down status for the port
  (must be "true" or "false" if provided). Changing this updates the
  `admin_state_up` of an existing port.

* `mac_address` - (Optional, String, ForceNew) Specify a specific MAC address for the port. Changing
  this creates a new port.

* `tenant_id` - (Optional, String, ForceNew) The owner of the Port. Required if admin wants
  to create a port for another tenant. Changing this creates a new port.

* `device_owner` - (Optional. String, ForceNew) The device owner of the Port. Changing this creates
  a new port.

* `security_group_ids` - (Optional, List) A list of security group IDs to apply to the
  port. The security groups must be specified by ID and not name (as opposed
  to how they are configured with the Compute Instance).

* `no_security_groups` - (Optional, Bool) If set to `true`, then no security groups
  are applied to the port. If set to `false` and no `security_group_ids` are specified,
  then the port will yield to the default behavior of the Networking service,
  which is to usually apply the `"default"` security group.

* `port_security_enabled` - (Optional, Bool) Whether to explicitly enable or disable
  port security on the port. Port Security is usually enabled by default, so
  omitting argument will usually result in a value of `true`. Setting this
  explicitly to `false` will disable port security. In order to disable port
  security, the port must not have any security groups. Valid values are `true`
  and `false`.

* `device_id` - (Optional, String, ForceNew) The ID of the device attached to the port. Changing this
  creates a new port.

* `fixed_ip` - (Optional, List) An array of desired IPs for this port. The structure is
  described below. A single `fixed_ip` entry is allowed for a port.
  The `fixed_ip` block supports:
  * `subnet_id` - (Required, String) Subnet in which to allocate IP address for
    this port.
  * `ip_address` - (Optional, String) IP address desired in the subnet for this port. If
    you don't specify `ip_address`, an available IP address from the specified
    subnet will be allocated to this port.

* `allowed_address_pairs` - (Optional, Map) An IP/MAC Address pair of additional IP
  addresses that can be active on this port. The structure is described below.
  The `allowed_address_pairs` block supports:
  * `ip_address` - (Required, String) The additional IP address.
  * `mac_address` - (Optional, String) The additional MAC address.

* `extra_dhcp_option` - (Optional, Map) Specifies the extended DHCP option. This is an extended attribute.
  The `extra_dhcp_option` block supports:
  * `name` - (Required, String) Specifies the option name.
  * `value` - (Required, String) Specifies the option value.

* `value_specs` - (Optional, Map) Map of additional options.


## Attributes Reference

The following attributes are exported:

* `admin_state_up` - See Argument Reference above.

* `mac_address` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `device_owner` - See Argument Reference above.

* `security_group_ids` - See Argument Reference above.

* `device_id` - See Argument Reference above.

* `fixed_ip` - See Argument Reference above.

* `all fixed_ips` - The collection of Fixed IP addresses on the port in the order returned by the Network v2 API.

* `port_security_enabled` - See Argument Reference above.

## Import

Ports can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_networking_port_v2.port_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
```

## Notes

### Ports and Instances

There are some notes to consider when connecting Instances to networks using
Ports. Please see the `opentelekomcloud_compute_instance_v2` documentation for further
documentation.
