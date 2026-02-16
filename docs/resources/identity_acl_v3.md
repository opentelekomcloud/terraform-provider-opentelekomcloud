---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_acl_v3"
sidebar_current: "docs-opentelekomcloud-resource-identity-acl-v3"
description: |-
  Manages a IAM ACL resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for IAM agency you can get at
[documentation portal](https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/security_settings/index.html)


# opentelekomcloud_identity_acl_v3

Manages an ACL resource within OpenTelekomCloud IAM service. The ACL allows user access only from specified IP address
ranges and CIDR blocks. The ACL takes effect for IAM users under the Domain account rather than the account itself.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

### ACL through console

```hcl
resource "opentelekomcloud_identity_acl_v3" "acl" {
  type = "console"

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a basic ip range for console access"
  }

  ip_cidrs {
    cidr        = "192.168.0.1/32"
    description = "This is a basic ip address for console access"
  }

  ipv6_ranges {
    range       = "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF"
    description = "This is a basic ipv6 range for console access"
  }

  ipv6_cidrs {
    cidr        = "0000:0000:0000:0000:0000:0000:0000:0000/100"
    description = "This is a basic ipv6 address for console access"
  }
}
```

### ACL through API

```hcl
resource "opentelekomcloud_identity_acl_v3" "acl" {
  type = "api"

  ip_cidrs {
    cidr        = "159.138.39.192/32"
    description = "This is a test ip address"
  }
  ip_ranges {
    range       = "0.0.0.0-255.255.255.0"
    description = "This is a test ip range"
  }
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, ForceNew) Specifies the ACL is created through the Console or API.
  Valid values are **console** and **api**. Changing this parameter will create a new ACL.

* `ip_cidrs` - (Optional, List) Specifies the IPv4 CIDR blocks from which console access or api access is allowed.
  The `ip_cidrs` cannot repeat. The structure is documented below.
    * `cidr` - (Required, String) Specifies the IPv4 CIDR block, for example, **192.168.0.0/24**.
    * `description` - (Optional, String) Specifies a description about an IPv4 CIDR block.

* `ip_ranges` - (Optional, List) Specifies the IP address ranges from which console access or api access is allowed.
  The `ip_ranges` cannot repeat. The structure is documented below.
    * `range` - (Required, String) Specifies the Ip address range, for example, **0.0.0.0-255.255.255.0**.
    * `description` - (Optional, String) Specifies a description about an IP address range.

* `ipv6_cidrs` - (Optional, List) Specifies the IPv6 CIDR blocks from which console access or api access is allowed.
  The `ipv6_cidrs` cannot repeat. The `ipv6_cidrs` can only be used when `type` is `console`. The structure is documented below.
    * `cidr` - (Required, String) Specifies the IPv6 CIDR block, for example, **0000:0000:0000:0000:0000:0000:0000:0000/100**.
    * `description` - (Optional, String) Specifies a description about an IPv6 CIDR block.

* `ipv6_ranges` - (Optional, List) Specifies the IPv6 address ranges from which console access or api access is allowed.
  The `ipv6_ranges` cannot repeat. The `ipv6_ranges` can only be used when `type` is `console`. The structure is documented below.
    * `range` - (Required, String) Specifies the IPv6 address range, for example, **0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF**.
    * `description` - (Optional, String) Specifies a description about an IPv6 address range.

-> **NOTE:** Up to 200 `ip_cidrs`, `ip_ranges`, `ipv6_cidrs`, and `ipv6_ranges` can be created in total for each access method.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of identity ACL.
