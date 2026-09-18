---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_loadbalancer_v3"
sidebar_current: "docs-opentelekomcloud-resource-lb-loadbalancer-v3"
description: |-
  Manages a LB Loadbalancer resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DLB load balancer you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/load_balancer)

# opentelekomcloud_lb_loadbalancer_v3

Manages a Dedicated loadbalancer resource within OpenTelekomCloud.

## Example Usage

### Basic usage

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  router_id   = var.router_id
  network_ids = [var.network_id]

  availability_zones = [var.az]

  tags = {
    muh = "kuh"
  }
}
```

### Usage with `vpc_subnet_v1`

```hcl
resource "opentelekomcloud_vpc_v1" "this" {
  name = "test-vpc-1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "this" {
  name       = "${opentelekomcloud_vpc_v1.this.name}-private"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.this.cidr, 8, 0)
  vpc_id     = opentelekomcloud_vpc_v1.this.id
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.this.cidr, 8, 0), 1)
  dns_list = [
    "1.1.1.1",
    "8.8.8.8",
  ]
}

resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  subnet_id   = opentelekomcloud_vpc_subnet_v1.this.subnet_id
  network_ids = [opentelekomcloud_vpc_subnet_v1.this.network_id]

  availability_zones = [var.az]
}
```

### Public load balancer (with floating IP)

#### Newly created

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  name        = "example-loadbalancer"
  subnet_id   = var.subnet_id
  network_ids = [var.network_id]

  availability_zones = [var.az]

  public_ip {
    bandwidth_name       = "lb-bandwidth"
    ip_type              = "5_bgp"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }
}
```

#### Already existing opentelekomcloud_networking_floatingip_v2

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["eu-de-01"]

  public_ip {
    id = opentelekomcloud_networking_floatingip_v2.fip_1.id
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}
```

#### Or opentelekomcloud_vpc_eip_v1

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = opentelekomcloud_vpc_subnet_v1.this.vpc_id
  network_ids = [opentelekomcloud_vpc_subnet_v1.this.network_id]

  availability_zones = ["eu-de-01"]

  public_ip {
    id = opentelekomcloud_vpc_eip_v1.fip_1.id
  }
}

resource "opentelekomcloud_vpc_eip_v1" "fip_1" {
  bandwidth {
    charge_mode = "traffic"
    name        = "eip"
    share_type  = "PER"
    size        = 100
  }

  publicip {
    type = "5_bgp"
  }
}
```

#### Assign new bandwidth to EIP without recreating

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_ids = [data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id]

  availability_zones = ["eu-de-01"]

  public_ip {
    ip_type              = "5_bgp"
    bandwidth_name       = "lb_band"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

resource "opentelekomcloud_vpc_bandwidth_v2" "bw" {
  name = "lb_band"
  size = 20
}

resource "opentelekomcloud_vpc_bandwidth_associate_v2" "associate" {
  bandwidth    = opentelekomcloud_vpc_bandwidth_v2.bw.id
  floating_ips = [opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.public_ip.0.id]
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Optional) ID of the router (or VPC) this LoadBalancer belongs to. Changing
  this creates a new LoadBalancer.

* `subnet_id` - (Optional) The ID of the subnet to which the LoadBalancer belongs. Required when using `vip_address`.

-> `router_id` and `subnet_id` cannot be left blank at the same time.

* `network_ids` - (Required) Specifies the subnet Network ID.

* `name` - (Optional) The LoadBalancer name.

* `description` - (Optional) Provides supplementary information about the load balancer.

* `vip_address` - (Optional) The ip address of the LoadBalancer. Changing this creates a new LoadBalancer.

-> Specify both `subnet_id` and `vip_address` if you want to bind a private IPv4 address to the load balancer.

* `admin_state_up` - (Optional) The administrative state of the LoadBalancer. A valid value is only `true` (UP).

* `ip_target_enable` - (Optional) The value can be `true` (enabled) or `false` (disabled).

-> If both `l4_flavor` and `l7_flavor` is empty, both ALB and NLB will be attached to the load balancer with the default flavor. It is advisable to specify one of them, unless your intention is to associate both flavors with the default setting.

* `l4_flavor` - (Optional) The ID of the Layer-4 flavor.

* `l7_flavor` - (Optional) The ID of the Layer-7 flavor.

* `availability_zones` - (Required) Specifies the availability zones where the LoadBalancer will be located.
  Changing this creates a new LoadBalancer.

* `public_ip` - (Optional) The elastic IP address of the instance. The `public_ip` structure
  is described below. Changing this creates a new LoadBalancer.

-> Specify `public_ip` and either `router_id` or `subnet_id` if you want to bind a new IPv4 EIP to the load balancer.

The `public_ip` block supports:

* `id` - (Optional) ID of an existing elastic IP. Required when using existing EIP.

* `ip_type` - (Optional) Elastic IP type. The value can be `5_bgp` or `5_mailbgp`.
  Required when creating a new EIP.

->
  In `eu-de` region the value can be `5_bgp` or `5_mailbgp`.
  In `eu-nl` region the value can only be `5_bgp` and `5_mailbgp`.
  In `eu-ch2` region the value can only be `5_bgp`.

* `bandwidth_name` - (Optional) Bandwidth name. Required when creating a new EIP.

* `bandwidth_size` - (Optional) Bandwidth size. Required when creating a new EIP.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Possible value is `traffic`.

* `bandwidth_share_type` - (Optional) Bandwidth sharing type. Possible values are: `PER`, `WHOLE`.
  Required when creating a new EIP.

* `deletion_protection` - (Optional) Specifies whether to enable deletion protection for the load balancer.
  `true`: Enable deletion protection.
  `false` (default): Disable deletion protection.

* `ipv6_vip_subnet_id` - (Optional) The ID of the IPv6 subnet where the load balancer resides.

* `ipv6_bandwidth_id` - (Optional) The ID of the shared bandwidth used by the public IPv6 address.

* `guaranteed` - (Optional) Whether the load balancer is a dedicated load balancer.
  Changing this creates a new load balancer.

* `enterprise_project_id` - (Optional) The enterprise project ID. Changing this creates a new load balancer.

* `charge_mode` - (Optional) The load balancer billing mode. Changing this creates a new load balancer.

* `protection_status` - (Optional) The modification protection status.

* `protection_reason` - (Optional) The reason for enabling modification protection.

* `tags` - (Optional, Map) Tags key/value pairs to associate with the load balancer.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vip_port_id` - The Port ID of the Load Balancer IP.

* `created_at` - The time the LoadBalancer was created.

* `updated_at` - The time the LoadBalancer was last updated.

* `provisioning_status` - The provisioning status of the load balancer.

* `operating_status` - The operating status of the load balancer.

* `provider_name` - The load balancer provider.

* `project_id` - The project ID of the load balancer.

* `pools` - IDs of backend server groups associated with the load balancer.

* `listeners` - IDs of listeners associated with the load balancer.

* `ipv6_vip_address` - The private IPv6 address of the load balancer.

* `ipv6_vip_port_id` - The port ID associated with the private IPv6 address.

* `l4_scale_flavor` - The ID of the Layer-4 elastic flavor.

* `l7_scale_flavor` - The ID of the Layer-7 elastic flavor.

* `eips` - EIPs associated with the load balancer, including `eip_id`, `eip_address`, and `ip_version`.

* `global_eips` - Global EIPs associated with the load balancer, including `global_eip_id`,
  `global_eip_address`, and `ip_version`.

* `elb_subnet_type` - The IP version supported by the load balancer subnet.

* `frozen_scene` - The scenario in which the load balancer was frozen.

* `billing_info` - Billing information for the load balancer.

* `autoscaling` - Autoscaling information, including `enable` and `min_l7_flavor_id`.

* `public_border_group` - The public border group of the load balancer.

* `waf_failure_action` - The action taken when WAF is unavailable.

* `loadbalancer_type` - The load balancer type.

* `gateway_flavor_id` - The gateway flavor ID.

* `instance_type` - The load balancer instance type.

* `instance_id` - The load balancer instance ID.

* `log_group_id` - The LTS log group ID.

* `log_topic_id` - The LTS log stream ID.

* `custom_qos_limit` - Custom Layer-4 and Layer-7 connection and CPS limits.

* `service_lb_mode` - The service load balancer mode.

* `proxy_protocol_extensions` - Proxy protocol extension endpoint information.

## Import

Loadbalancers can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_lb_loadbalancer_v3.lb_1 7b80e108-1636-44e5-aece-986b0052b7dd
```
