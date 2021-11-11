---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_loadbalancer_v3

Use this data source to get the info about an existing ELBv3 load balancer.

## Example Usage

```hcl
data "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  id = var.lb_id
}
```

## Argument Reference

* `id` - (Optional) Specifies the LoadBalancer ID.

* `name` - (Optional) Specifies the LoadBalancer name.

* `router_id` - (Optional) The ID of the router (or VPC) this LoadBalancer belongs.

* `subnet_id` - (Optional) The ID of the subnet to which the LoadBalancer belongs.

* `l7_flavor_id` - (Optional) The ID of the Layer-7 flavor.

* `l4_flavor_id` - (Optional) The ID of the Layer-4 flavor.

* `vip_address` - (Optional) The IP address of the LoadBalancer.

* `vip_port_id` - (Optional) The Port ID of the Load Balancer IP.

## Attributes Reference

In addition, the following attributes are exported:

* `type` - Specifies the flavor type.

* `shared` - Specifies whether the flavor is available to all users.

* `max_connections` - Specifies the maximum concurrent connections.

* `cps` - Specifies the number of new connections per second.

* `qps` - Specifies the number of requests per second at Layer 7.

* `bandwidth` - Specifies the inbound and outbound bandwidth in the unit of Kbit/s.
