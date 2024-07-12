---
subcategory: "Elastic Load Balancer (ELB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_loadbalancer_v2"
sidebar_current: "docs-opentelekomcloud-resource-lb-loadbalancer-v2"
description: |-
  Manages a ELB Loadbalancer resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for ELB load balancer you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v2.0/load_balancer)

# opentelekomcloud_lb_loadbalancer_v2

Manages an Enhanced loadbalancer resource within OpenTelekomCloud.

## Example Usage

### Basic usage

```hcl
resource "opentelekomcloud_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  tags = {
    muh = "kuh"
  }
}
```

### Usage with `vpc_subnet_v1`

```hcl
resource "opentelekomcloud_vpc_v1" "main" {
  cidr = "192.168.0.0/16"
  name = "test-vpc-1"
}

resource "opentelekomcloud_vpc_subnet_v1" "private" {
  name       = "${opentelekomcloud_vpc_v1.main.name}-private"
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.main.cidr, 8, 0)
  vpc_id     = opentelekomcloud_vpc_v1.main.id
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.main.cidr, 8, 0), 1)
  dns_list = [
    "1.1.1.1",
    "8.8.8.8",
  ]
}

resource "opentelekomcloud_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = opentelekomcloud_vpc_subnet_v1.private.subnet_id
}
```

### Public load balancer (with floating IP)

```hcl
resource "opentelekomcloud_lb_loadbalancer_v2" "lb_1" {
  name          = "example-loadbalancer"
  vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "associate" {
  floating_ip = var.floating_ip_address
  port_id     = opentelekomcloud_lb_loadbalancer_v2.lb_1.vip_port_id
}
```

## Argument Reference

The following arguments are supported:

* `vip_subnet_id` - (Required) The network on which to allocate the
  loadbalancer's address. A tenant can only create loadalancers on networks
  authorized by policy (e.g. networks that belong to them or networks that
  are shared). Changing this creates a new loadbalancer.

-> When used with `opentelekomcloud_vpc_subnet_v1`, not `id` but
`subnet_id`needs to be used

* `name` - (Optional) Human-readable name for the loadbalancer. Does not have
  to be unique.

* `description` - (Optional) Human-readable description for the loadbalancer.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the Loadbalancer.  Only administrative users can specify a tenant UUID
  other than their own. Changing this creates a new loadbalancer.

* `vip_address` - (Optional) The ip address of the load balancer.
  Changing this creates a new loadbalancer.

* `admin_state_up` - (Optional) The administrative state of the loadbalancer.
  A valid value is only true (UP).

* `loadbalancer_provider` - (Optional) The name of the provider. Changing this
  creates a new loadbalancer.

* `tags` - (Optional) Tags key/value pairs to associate with the loadbalancer.


## Attributes Reference

The following attributes are exported:

* `vip_subnet_id` - See Argument Reference above.

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `vip_address` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `loadbalancer_provider` - See Argument Reference above.

* `vip_port_id` - The Port ID of the Load Balancer IP.

* `tags` - See Argument Reference above.

## Import

Load balancers can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_lb_loadbalancer_v2.lb_1 ec2e6489-8415-4ec0-9934-540f98b0d594
```
