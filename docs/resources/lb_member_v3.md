---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_member_v3"
sidebar_current: "docs-opentelekomcloud-resource-lb-member-v3"
description: |-
  Manages a LB Member resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DLB member you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/backend_server)

# opentelekomcloud_lb_member_v3

Manages a Dedicated Load Balancer member resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  ip_target_enable = true

  availability_zones = [var.availability_zone]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  name            = "pool_1"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.lb.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "TCP"
}

resource "opentelekomcloud_lb_member_v3" "member" {
  name          = "member-1"
  pool_id       = opentelekomcloud_lb_pool_v3.pool.id
  address       = cidrhost(var.subnet_cidr, 3)
  protocol_port = 8080
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required) Specifies the IP address of the backend server.

  The IP address must be in the subnet specified by `subnet_id`, for example, `192.168.3.11`.

  The IP address can only be the IP address of the primary NIC.

  If `subnet_id` is left blank, cross-VPC backend is enabled. In this case, these servers must use IPv4 addresses.

* `protocol_port` - (Required) Specifies the port used by the backend server to receive requests. The value should be a
  valid port.

* `name` - (Optional) Specifies the backend server name. The value is a string of 0 to 255 characters.

* `project_id` - (Optional) Specifies the project ID.

* `subnet_id` - (Optional) Specifies the ID of the subnet where the backend server works.

  This subnet must be in the same VPC as the subnet of the load balancer with which the backend server is associated.

  Only `IPv4` subnets are supported.

* `weight` - (Optional) Specifies the weight of the backend server.

  Requests are routed to backend servers in the same backend server group based on their weights.

  If the weight is `0`, the backend server will not accept new requests.

  This parameter is invalid when `lb_algorithm` is set to `SOURCE_IP` for the backend server group that contains the
  backend server.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `member_id` - ID of the pool member.

* `operating_status` - Specifies the operating status of the backend server.

  The value can be one of the following:
    * `ONLINE`: The backend server is running normally.
    * `NO_MONITOR`: No health check is configured for the backend server group to which the backend server belongs.
    * `OFFLINE`: The cloud server used as the backend server is stopped or does not exist.

* `ip_version` - Version of IP based on the `address` parameter. The value can be `v4` or `v6`.

## Import

Members can be imported using the `pool_id/member_id`, e.g.

```sh
terraform import opentelekomcloud_lb_member_v3.member 7b80e108-1636-44e5-aece-986b0052b7dd/1bb93b8b-37a4-4b50-92cc-daa4c89d4e4c
```
