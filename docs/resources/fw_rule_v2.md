---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fw_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-fw-rule-v2"
description: |-
Manages a VPC Firewall Rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC firewall rule you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/firewall)

# opentelekomcloud_fw_rule_v2

Manages a v2 firewall rule resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_fw_rule_v2" "rule_1" {
  name             = "my_rule"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
  enabled          = "true"
}
```

## Example Ipv6 Usage
```hcl
resource "opentelekomcloud_fw_rule_v2" "rule_1" {
  name        = "rule_1"
  description = "Ipv6 deny"
  protocol    = "tcp"
  ip_version  = 6
  enabled     = true
  action      = "deny"

  destination_ip_address = "2001:db8::"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) A unique name for the firewall rule. Changing this
  updates the `name` of an existing firewall rule.

* `description` - (Optional) A description for the firewall rule. Changing this
  updates the `description` of an existing firewall rule.

* `protocol` - (Required) The protocol type on which the firewall rule operates.
  Valid values are: `tcp`, `udp`, `icmp`, and `any`. Changing this updates the
  `protocol` of an existing firewall rule.

* `action` - (Required) Action to be taken ( must be "allow" or "deny") when the
  firewall rule matches. Changing this updates the `action` of an existing
  firewall rule.

* `ip_version` - (Optional) IP version, either 4 (default) or 6. Changing this
  updates the `ip_version` of an existing firewall rule.

* `source_ip_address` - (Optional) The source IP address on which the firewall
  rule operates. Changing this updates the `source_ip_address` of an existing
  firewall rule.

* `destination_ip_address` - (Optional) The destination IP address on which the
  firewall rule operates. Changing this updates the `destination_ip_address`
  of an existing firewall rule.

* `source_port` - (Optional) The source port on which the firewall
  rule operates. Changing this updates the `source_port` of an existing
  firewall rule.

* `destination_port` - (Optional) The destination port on which the firewall
  rule operates. Changing this updates the `destination_port` of an existing
  firewall rule.

* `enabled` - (Optional) Enabled status for the firewall rule (must be "true"
  or "false" if provided - defaults to "true"). Changing this updates the
  `enabled` status of an existing firewall rule.

* `tenant_id` - (Optional) The owner of the firewall rule. Required if admin
  wants to create a firewall rule for another tenant. Changing this creates a
  new firewall rule.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `protocol` - See Argument Reference above.

* `action` - See Argument Reference above.

* `ip_version` - See Argument Reference above.

* `source_ip_address` - See Argument Reference above.

* `destination_ip_address` - See Argument Reference above.

* `source_port` - See Argument Reference above.

* `destination_port` - See Argument Reference above.

* `enabled` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

## Import

Firewall Rules can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_fw_rule_v2.rule_1 8dbc0c28-e49c-463f-b712-5c5d1bbac327
```
