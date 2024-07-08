---
subcategory: "Elastic Load Balancer (ELB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_member_v2"
sidebar_current: "docs-opentelekomcloud-resource-lb-member-v2"
description: |-
Manages a ELB Member resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for ELB member you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v2.0/backend_server)

# opentelekomcloud_lb_member_v2

Manages an Enhanced Load Balancer member resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "pool_id" {}
variable "subnet_id" {}

resource "opentelekomcloud_lb_member_v2" "member1" {
  address       = "192.168.199.23"
  protocol_port = 8080
  pool_id       = var.pool_id
  subnet_id     = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `pool_id` - (Required) The id of the pool that this member will be
  assigned to.

* `subnet_id` - (Required) The subnet in which to access the member

* `name` - (Optional) Human-readable name for the member.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the member.  Only administrative users can specify a tenant UUID
  other than their own. Changing this creates a new member.

* `address` - (Required) The IP address of the member to receive traffic from
  the load balancer. Changing this creates a new member.

* `protocol_port` - (Required) The port on which to listen for client traffic.
  Changing this creates a new member.

* `weight` - (Optional)  A positive integer value that indicates the relative
  portion of traffic that this member should receive from the pool. For
  example, a member with a `weight` of `10` receives five times as much traffic
  as a member with a `weight` of `2`. If the value is `0`, the backend server will not accept new requests

* `admin_state_up` - (Optional) The administrative state of the member.
  A valid value is `true` (UP) or `false` (DOWN).

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the member.

* `name` - See Argument Reference above.

* `weight` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `subnet_id` - See Argument Reference above.

* `pool_id` - See Argument Reference above.

* `address` - See Argument Reference above.

* `protocol_port` - See Argument Reference above.

## Import

Members can be imported using the `pool_id/member_id`, e.g.

```sh
terraform import opentelekomcloud_lb_member_v2.member 7b80e108-1636-44e5-aece-986b0052b7dd/1bb93b8b-37a4-4b50-92cc-daa4c89d4e4c
```
