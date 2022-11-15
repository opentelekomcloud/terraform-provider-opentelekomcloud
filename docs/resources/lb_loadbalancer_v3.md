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
    ip_type              = "5_gray"
    bandwidth_size       = 10
    bandwidth_share_type = "PER"
  }
}
```

#### Already existing

```hcl
resource "opentelekomcloud_lb_loadbalancer_v3" "lb_1" {
  name        = "example-loadbalancer"
  subnet_id   = var.subnet_id
  network_ids = [var.network_id]

  availability_zones = [var.az]

  public_ip {
    id = var.eip_id
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
    ip_type              = "5_gray"
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

* `subnet_id` - (Optional) The ID of the subnet to which the LoadBalancer belongs.

-> `router_id` and `subnet_id` cannot be left blank at the same time.

* `network_ids` - (Required) Specifies the subnet Network ID.

* `name` - (Optional) The LoadBalancer name.

* `description` - (Optional) Provides supplementary information about the load balancer.

* `vip_address` - (Optional) The ip address of the LoadBalancer. Changing this creates a new LoadBalancer.

-> Specify both `subnet_id` and `vip_address` if you want to bind a private IPv4 address to the load balancer.

* `admin_state_up` - (Optional) The administrative state of the LoadBalancer. A valid value is only `true` (UP).

* `ip_target_enable` - (Optional) The value can be `true` (enabled) or `false` (disabled).

* `l4_flavor` - (Optional) The ID of the Layer-4 flavor.

* `l7_flavor` - (Optional) The ID of the Layer-7 flavor.

* `availability_zones` - (Required) Specifies the availability zones where the LoadBalancer will be located.
  Changing this creates a new LoadBalancer.

* `public_ip` - (Optional) The elastic IP address of the instance. The `public_ip` structure
  is described below. Changing this creates a new LoadBalancer.

-> Specify `public_ip` and either `router_id` or `subnet_id` if you want to bind a new IPv4 EIP to the load balancer.

The `public_ip` block supports:

* `id` - (Optional) ID of an existing elastic IP. Required when using existing EIP.

* `ip_type` - (Optional) Elastic IP type. The value can be `5_gray` and `5_mailbgp`.
  Required when creating a new EIP.

->
  In `eu-de` region the value can only be `5_gray`.

* `bandwidth_name` - (Optional) Bandwidth name. Required when creating a new EIP.

* `bandwidth_size` - (Optional) Bandwidth size. Required when creating a new EIP.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Possible value is `traffic`.

* `bandwidth_share_type` - (Optional) Bandwidth sharing type. Possible values are: `PER`, `WHOLE`.
  Required when creating a new EIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vip_port_id` - The Port ID of the Load Balancer IP.

* `created_at` - The time the LoadBalancer was created.

* `updated_at` - The time the LoadBalancer was last updated.

## Import

Loadbalancers can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_lb_loadbalancer_v3.lb_1 7b80e108-1636-44e5-aece-986b0052b7dd
```
