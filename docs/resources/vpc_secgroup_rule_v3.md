---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_secgroup_rule_v3"
sidebar_current: "docs-opentelekomcloud-resource-vpc-secgroup-rule-v3"
description: |-
  Manages a VPC Security Group Rule v3 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC security group rule you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/api_v3/security_group_rule/index.html)

# opentelekomcloud_vpc_secgroup_rule_v3

Manages a VPC security group rule v3 resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "security_group_id" {}

resource "opentelekomcloud_vpc_secgroup_rule_v3" "rule_1" {
  security_group_id = var.security_group_id
  description       = "some basic security rule"
  direction         = "ingress"
  protocol          = "tcp"
  action            = "allow"
  priority          = 1
  multi_port        = "8080"
  remote_ip_prefix  = "10.10.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required, String, ForceNew) Specifies the ID of the security group to which the security group rule belongs.

* `description` - (Optional, String, ForceNew) Provides supplementary information about the security group rule.

* `direction` - (Required, String, ForceNew) Specifies inbound or outbound direction of a security group rule. Supported values: `ingress` (inbound direction), `egress` (outbound direction).

* `ether_type` - (Required, String, ForceNew) Specifies the IP version. Supported values: `IPv4`, `IPv6`. Default: `IPv4 `.

* `protocol` - (Optional, String, ForceNew) Specifies the protocol type. The value can be `icmp`, `tcp`, `udp`, `icmpv6` or an `IP number (0 to 255)`. If the parameter is left blank, **all** protocols are supported. When the protocol is `icmpv6`, IP version should be `IPv6`. When the protocol is `icmp`, IP version should be `IPv4`.

* `multiport` - (Optional, String, ForceNew) Specifies the port or port range. The value can be a single port, e.g. `80`, a port range, e.g. `1-30`, or inconsecutive ports separated by commas, e.g. `22,3389,80`.

* `remote_ip_prefix` - (Optional, String, ForceNew) Specifies the remote IP address. If `direction` is set to `egress`, the parameter specifies the source IP address. If `direction` is set to `ingress`, the parameter specifies the destination IP address. The value is an `IP address` or a `CIDR block`. The parameter is mutually exclusive with parameter `remote_group_id`. If this parameter is left blank, the remote IP address is not limited, and the traffic from all remote IP addresses is allowed or rejected.

* `remote_group_id` - (Optional, String, ForceNew) Specifies the ID of the remote security group, which allows or denies traffic to and from the security group. The value has to be the ID of an existing security group. The parameter is mutually exclusive with parameter `remote_ip_prefix`.

* `action` - (Optional, String, ForceNew) Specifies the action of the security group rule. Supported values: `allow`, `deny`. Default value: `allow`.

* `priority` - (Optional, Integer, ForceNew) Specifies the rule priority in a security group. The value is from 1 to 100. The value 1 indicates the highest priority. Default value: `1`.


## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Security Group Rule ID.

* `project_id` - Indicates the project ID.

* `created_at` - Indicates the time when the security group rule was created. It is a UTC time in yyyy-MM-ddTHH:mm:ssZ format.

* `updated_at` - Indicates the time when the security group rule was updated. It is a UTC time in yyyy-MM-ddTHH:mm:ssZ format.

* `remote_address_group_id` - Indicates the ID of the remote IP address group. The parameter value is mutually exclusive with parameters `remote_ip_prefix` and `remote_group_id`.

## Import

VPC Security Group Rule V3 can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_secgroup_rule_v3.secgroup_rule_1 <id>
```
