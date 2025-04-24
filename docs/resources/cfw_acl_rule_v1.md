---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_acl_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-acl-rule-v1"
description: |-
  Manages a CFW ACL rule resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW ACL rule you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/acl_rule_management/index.html)

# opentelekomcloud_cfw_acl_rule_v1

Manages a CFW ACL rule resource within OpenTelekomCloud.

## Example Usage:
```hcl
variable object_id {}

resource "opentelekomcloud_cfw_acl_rule_v1" "rule_1" {
  object_id = var.object_id
  type      = 0
  name      = "test-acc-tf-acl-rule"
  sequence {
    top = 1
  }
  address_type        = 0
  action_type         = 0
  status              = 1
  long_connect_enable = 0
  direction           = 0
  source {
    type    = 0
    address = "1.1.1.1"
  }
  destination {
    type    = 0
    address = "2.2.2.2"
  }
  service {
    type     = 0
    protocol = -1
  }
}
```

## Argument Reference

The following arguments are supported:

* `object_id` - (Required, String, ForceNew) Protected object ID, which is used to distinguish between Internet border protection and VPC border protection after a cloud firewall is created. If the value of type is 0, the protected object ID belongs to the Internet border. If the value of type is 1, the protected object ID belongs to the VPC border.

* `type` - (Required, Integer) Rule type: `0` (Internet border rule), `1` (inter-VPC rule), or `2` (NAT rule). When type is set to `0`, the source and destination addresses of the rule must be EIPs or domain names of the public network. For an inter-VPC rule, the source and destination addresses must be private IP addresses. For a NAT rule, the source address must be a private IP address, and the destination address must be an EIP or domain name of the public network.

* `name` - (Required, String) Specifies the CFW ACL rule name. The CFW ACL rule name of the same type is unique in the same firewall instance.

* `sequence` - (Required, List, ForceNew) Specifies the request body for changing the rule sequence. The [sequence](#sequence) structure is documented below.

* `address_type` - (Required, Integer) Specifies the Internet protocol type of an address: `0` (IPv4), `1` (IPv6).

* `action_type` - (Required, Integer) Specifies the rule action: `0` (permit), `1` (deny). 

* `status` - (Required, Integer) Specifies the rule status: `0` (disabled), `1` (enabled).

* `applications` - (Optional, List) Specifies the rule application list . Allowed list values: `HTTP`, `HTTPS`, `TLS1`, `DNS`, `SSH`, `MYSQL`, `SMTP`, `RDP`, `RDPS`, `VNC`, `POP3`, `IMAP4`, `SMTPS`, `POP3S`, `FTPS`, `ANY`, or `BGP`.

* `applications_json_string` - (Optional, String) Specifies the JSON string converted from the `applications` field in the application list.

* `long_connect_time` - (Optional, Int) Specifies the persistent connection duration.

* `long_connect_time_hour` - (Optional, Int) Specifies the persistent connection duration (hour).

* `long_connect_time_minute` - (Optional, Int) Specifies the persistent connection duration (minute).

* `long_connect_time_second` - (Optional, Int) Specifies the persistent connection duration (second).

* `long_connect_enable` - (Required, Int) Specifies whether to support persistent connections: `0` (no), `1` (yes).

* `description` - (Optional, String) Specifies the description of the rule.

* `direction` - (Optional, Integer) Specifies the Direction: `0` (inbound) or `1` (outbound). This parameter is **mandatory** **when** `type` is set to `0` (Internet rule) or `2` (NAT rule).
  
* `source` - (Required, List) Specifies the source address Data Transport Object.
  The [source](#dto) structure is documented below.

* `destination` - (Required, List) Specifies the destination address Data Transport Object.
  The [destination](#dto) structure is documented below.

* `service` - (Required, List) Specifies the service object.
  The [service](#service) structure is documented below.

<a name="sequence"></a>
The `sequence` block supports:

* `dest_rule_id` - (Optional, String, ForceNew) Specifies the ID of the target rule.
* `top` - (Optional, Integer, ForceNew) Specifies whether to pin on top: `0` (no), `1` (yes).
* `bottom` - (Optional, Integer, ForceNew) Specifies whether to pin to bottom: `0` (no), `1` (yes).

<a name="dto"></a>
The `source` and `destination` block supports:

* `type` - (Required, Integer) Specifies the Address type: `0` (manual input), `1` (associated IP address group), `2` (domain name), `3` (geographical location), `4` (domain name group) `5` (multiple objects), `6` (domain name group - network), `7` (domain name group - application).

* `address_type` - (Optional, Integer) Specifies theInternet protocol type of an address: `0` (IPv4), `1` (IPv6). If type is 0, this parameter cannot be left blank.

* `address` - (Optional, String) Specifies the IP address information. It cannot be left blank if type is set to 0.

* `address_set_id` - (Optional, String) Specifies the ID of an associated IP address group. This parameter cannot be left blank when type is set to 1. 

* `address_set_name` - (Optional, String) Specifies the name of an associated IP address group. This parameter cannot be left blank when type is set to 1.

* `domain_address_name` - (Optional, String) Specifies the name of a domain name address. This parameter is valid when type is set to 2 (domain name) or 7 (application domain name group).

* `region_list_json` - (Optional, String) Specifies the JSON value of the rule region list.

* `region_list` - (Optional, List) Specifies the rule region list.
    -  `region_id` - (Optional, String) Specifies the region ID.
    -  `region_type` - (Optional, Integer) Specifies the region type: `0` (country), `1` (province), and `2` (continent).

* `domain_set_id` - (Optional, String) Specifies the domain group ID. The value cannot be left blank when type is set to 4 (domain name group) or 7 (domain name group - application). 

* `domain_set_name` - (Optional, String) Specifies the domain group name. The value cannot be left blank when type is set to 4 (domain name group) or 7 (domain name group - application). 

* `ip_address` - (Optional, List) Specifies the IP address list. This parameter cannot be left blank when type is set to 5 (multiple objects).

* `address_set_type` - (Optional, Integer) Specifies the Address group type. It cannot be left blank when type is set to 1 (associated IP address group). It value can be `0` (user-defined address group), `1` (WAF back-to-source IP address group), `2` (DDoS back-to-source IP address group), or `3` (NAT64 address group).

* `predefined_group` - (Optional, List) Specifies the pre-defined address group ID list. This parameter cannot be left blank when type is set to 5 (multiple objects).

* `address_group` - (Optional, List) Specifies the Address group ID list. This parameter cannot be left blank when type is set to 5 (multiple objects).

<a name="service"></a>
The `service` block supports:

* `type` - (Required, Integer) Specifies the service input type: `0` (manual), `1` (automatic).

* `protocol` - (Optional, Integer) Specifies the protocol type: `6` (TCP), `17` (UDP), `1` (ICMP), `58` (ICMPv6), or `-1` (any). It cannot be left blank when type is set to 0 (manual).

* `protocols` - (Optional, List) Specifies the protocol list. Permitted list values: `6` (TCP), `17` (UDP), `1` (ICMP), `58` (ICMPv6), or `-1` (any). It cannot be left blank when type is set to 0 (manual).

* `source_port` - (Optional, String) Specifies the source port.

* `dest_port` - (Optional, String) Specifies the destination port.

* `service_set_id` - (Optional, String) Specifies the Service group ID. This parameter cannot be left blank when type is set to 1 (associated IP address group).

* `service_set_name` - (Optional, String) Specifies the Service group name. This parameter cannot be left blank when type is set to 1 (associated IP address group).

* `custom_service` - (Optional, List) Specifies the custom service.
    -  `protocol` - (Optional, Integer) Specifies the protocol type: `6` (TCP), `17` (UDP), `1` (ICMP), `58` (ICMPv6), or `-1` (any). It cannot be left blank when `type` is set to 0 (manual) in `service` block.
    -  `source_port` - (Optional, String) Specifies the source port.
    -  `dest_port` - (Optional, String) Specifies the destination port.
    -  `description` - (Optional, String) Specifies the service member description.
    -  `name` - (Optional, String) Specifies the service member name.

* `predefined_group` - (Optional, List) Specifies the pre-defined service group ID list.

* `service_group` - (Optional, List) Specifies the service group ID list.

* `service_group_names` - (Optional, List) Specifies the service group name list.
    -  `name` - (Optional, String) Specifies the service group name.
    -  `protocols` - (Optional, List) Specifies the protocols list. Permitted list values: `6` (TCP), `17` (UDP), `1` (ICMP), `58` (ICMPv6), or `-1` (any).
    -  `service_set_type` - (Optional, Integer) Specifies the service group type: `0` (user-defined service group), `1` (predefined service group).
    -  `set_id` - (Optional, String) Specifies the service group ID.

* `service_set_type` - (Optional, Integer) Specifies the service group type: `0` (user-defined service group), `1` (common web service), `2` (common remote login and ping), or `3` (common database).


## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the ACL rule ID.
* `created_date` - Indicates the Rule creation time in YYYY-MM-DD hh:mm:ss format.
* `last_open_time` - Indicates the Last time when the rule was enabled in  YYYY-MM-DD hh:mm:ss format.

## Import

CFW ACL V1 Rule can be imported using the CFW Firewall protection object ID, `object_id` and rule name `name`, e.g.

```shell
terraform import opentelekomcloud_cfw_acl_rule_v1.rule_1 b4cd6aeb0b7445d3bf271457c6941544in09/name
```

## Notes

But due to some attributes missing from the API response, it's required to ignore changes as below:

```hcl
resource "opentelekomcloud_cfw_acl_rule_v1" "rule_1" {
  # ...

  lifecycle {
    ignore_changes = [
      "sequence",
      "applications",
      "applications_json_string",
      "source.0.predefined_group",
      "destination.0.predefined_group",
      "service.0.predefined_group",
    ]
  }
}
```
