---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_vpc_subnet_v1

Provides an VPC v1 subnet resource within OpenTelekomCloud.

## Example Usage

### Basic Usage

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_v1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_v1" {
  name   = var.subnet_name
  cidr   = var.subnet_cidr
  vpc_id = opentelekomcloud_vpc_v1.vpc_v1.id

  gateway_ip    = var.subnet_gateway_ip
  ntp_addresses = "10.100.0.33,10.100.0.34"
}
```

### Subnet with tags

```hcl
resource "opentelekomcloud_vpc_subnet_v1" "subnet_with_tags" {
  name   = var.subnet_name
  cidr   = var.subnet_cidr
  vpc_id = opentelekomcloud_vpc_v1.vpc_v1.id

  gateway_ip    = var.subnet_gateway_ip
  ntp_addresses = "10.100.0.33,10.100.0.34"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The subnet name. The value is a string of `1` to `64` characters that can contain letters,
  digits, underscores (`_`), and hyphens (`-`).

* `description` - (Optional) A description of the VPC subnet.

* `cidr` - (Required) Specifies the network segment on which the subnet resides. The value must be in CIDR format.
  The value must be within the CIDR block of the VPC. The subnet mask cannot be greater than `28`.
  Changing this creates a new Subnet.

* `gateway_ip` - (Required) Specifies the gateway of the subnet. The value must be a valid IP address.
  The value must be an IP address in the subnet segment. Changing this creates a new Subnet.

* `vpc_id` - (Required) Specifies the ID of the VPC to which the subnet belongs. Changing this creates a new Subnet.

* `dhcp_enable` - (Optional) Specifies whether the DHCP function is enabled for the subnet. The value can
  be `true` or `false`. If this parameter is left blank, it is set to `true` by default.

* `primary_dns` - (Optional) Specifies the IP address of DNS server 1 on the subnet. The value must be a
  valid IP address. Default is `100.125.4.25`, OpenTelekomCloud internal DNS server.

* `secondary_dns` - (Optional) Specifies the IP address of DNS server 2 on the subnet. The value must be a
  valid IP address. Default is `100.125.129.199`, OpenTelekomCloud secondary internal DNS server.

* `dns_list` - (Optional) Specifies the DNS server address list of a subnet. This field is required if you
  need to use more than two DNS servers. This parameter value is the superset of both DNS server address
  1 and DNS server address 2.

~>
  Please note that primary DNS should be set to OTC-internal for managed services (e.g. CCE, CSS) to work.

* `availability_zone` - (Optional) Identifies the availability zone (AZ) to which the subnet belongs.
  The value must be an existing AZ in the system. Changing this creates a new Subnet.

* `ntp_addresses` - (Optional) Specifies the NTP server address configured for the subnet.

* `tags` - (Optional) The key/value pairs to associate with the subnet.


## Attributes Reference

All the argument attributes are also exported as result attributes:

* `id` - Specifies a resource ID in UUID format. Same as OpenStack network ID (`OS_NETWORK_ID`).

* `status` - Specifies the status of the subnet. The value can be `ACTIVE`, `DOWN`, `UNKNOWN`, or `ERROR`.

* `subnet_id` - Specifies the OpenStack subnet ID.

* `network_id` - Specifies the OpenStack network ID.

## Import

Subnets can be imported using the `subnet id`, e.g.

```shell
terraform import opentelekomcloud_vpc_subnet_v1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
