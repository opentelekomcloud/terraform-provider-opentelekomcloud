---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_loadbalancer_v3

Manages a Dedicated loadbalancer resource within OpenTelekomCloud.

## Example Usage

### Basic usage

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  router_id   = var.router_id
  subnet_id   = var.subnet_id
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
  dns_list   = [
    "1.1.1.1",
    "8.8.8.8",
  ]
}

resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  router_id   = opentelekomcloud_vpc_subnet_v1.this.vpc_id
  subnet_id   = opentelekomcloud_vpc_subnet_v1.this.subnet_id
  network_ids = [opentelekomcloud_vpc_subnet_v1.this.network_id]

  availability_zones = [var.az]
}
```

### Public load balancer (with floating IP)

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  name        = "example-loadbalancer"
  router_id   = var.router_id
  subnet_id   = var.subnet_id
  network_ids = [var.network_id]

  availability_zones = [var.az]

  public_ip {
    ip_type              = "5_bgp"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required) ID of the router (or VPC) this LoadBalancer belongs to. Changing
  this creates a new loadbalancer.

* `subnet_id` - (Required) The ID of the subnet to which the LoadBalancer belongs.

* `network_ids` - (Required) Specifies the subnet Network ID.

* `name` - (Optional) Human-readable name for the LoadBalancer.

* `description` - (Optional) Human-readable description for the LoadBalancer.

* `vip_address` - (Optional) The ip address of the LoadBalancer. Changing this creates a new LoadBalancer.

* `admin_state_up` - (Optional) The administrative state of the LoadBalancer. A valid value is only `true` (UP).

* `ip_target_enable` - (Optional) The value can be `true` (enabled) or `false` (disabled).

* `l4_flavor` - (Optional) The ID of the Layer-4 flavor.

* `l7_flavor` - (Optional) The ID of the Layer-7 flavor.

* `availability_zones` - (Required) Specifies the availability zones where the LoadBalancer will be located.
  Changing this will create a new resource.

* `public_ip` - (Optional) The elastic IP address of the instance. The `public_ip` structure
  is described below. Changing this will create a new resource.

The `public_ip` block supports:

* `ip_type` - (Required) Elastic IP type. Possible values are: `5_bgp`, `5_mailbgp`, `5_gray`.

* `bandwidth_size` - (Required) Bandwidth size.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Possible value is `traffic`.

* `bandwidth_share_type` - (Required) Bandwidth sharing type. Possible values are: `PER`, `WHOLE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vip_port_id` - The Port ID of the Load Balancer IP.

* `created_at` - The time the LoadBalancer was created.

* `updated_at` - The time the LoadBalancer was last updated.
