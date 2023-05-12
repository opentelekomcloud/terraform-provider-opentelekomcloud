---
subcategory: "Virtual Private Cloud (VPC)"
---

Up-to-date reference of API arguments for VPC security group rule you can get at
`https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/security_group`.

# opentelekomcloud_networking_secgroup_rule_v2

Manages a V2 neutron security group rule resource within OpenTelekomCloud.
Unlike Nova security groups, neutron separates the group from the rules
and also allows an admin to target a specific tenant_id.

## Example Usage

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the rule. Changing this creates a new security group rule.

* `direction` - (Required) The direction of the rule, valid values are `ingress`
  or `egress`. Changing this creates a new security group rule.

* `ethertype` - (Required) The layer 3 protocol type, valid values are `IPv4`
  or `IPv6`. Changing this creates a new security group rule.

* `protocol` - (Optional) The layer 4 protocol type, valid values are following. Changing this creates a new security group rule.
  This is required if you want to specify a port range.
  * `tcp`, `udp`, `icmp`, `ah`, `dccp`, `egp`, `esp`, `gre`, `igmp`, `ipv6-encap`,
  `ipv6-frag`, `ipv6-icmp`, `ipv6-nonxt`, `ipv6-opts`, `ipv6-route`, `ospf`,
  `pgm`, `rsvp`, `sctp`, `udplite`, `vrrp`

* `port_range_min` - (Optional) The lower part of the allowed port range, valid
  integer value needs to be between 1 and 65535. Changing this creates a new
  security group rule.

* `port_range_max` - (Optional) The higher part of the allowed port range, valid
  integer value needs to be between 1 and 65535. Changing this creates a new
  security group rule.

* `remote_ip_prefix` - (Optional) The remote CIDR, the value needs to be a valid
  CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule.

* `remote_group_id` - (Optional) The remote group id, the value needs to be an
  OpenTelekomCloud ID of a security group in the same tenant. Changing this creates
  a new security group rule.

* `security_group_id` - (Required) The security group id the rule should belong
  to, the value needs to be an OpenTelekomCloud ID of a security group in the same
  tenant. Changing this creates a new security group rule.

* `tenant_id` - (Optional) The owner of the security group. Required if admin
  wants to create a port for another tenant. Changing this creates a new
  security group rule.

## Attributes Reference

The following attributes are exported:

* `description` - See Argument Reference above.

* `direction` - See Argument Reference above.

* `ethertype` - See Argument Reference above.

* `protocol` - See Argument Reference above.

* `port_range_min` - See Argument Reference above.

* `port_range_max` - See Argument Reference above.

* `remote_ip_prefix` - See Argument Reference above.

* `remote_group_id` - See Argument Reference above.

* `security_group_id` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

## Import

Security Group Rules can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
