---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_loadbalancer_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lb-loadbalancer-v3"
description: |-
  Get details about ELBv3 loadbalancer from OpenTelekomCloud
---

Up-to-date reference of API arguments for ELBv3 loadbalancer you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/load_balancer/querying_load_balancers.html#listloadbalancers)

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

* `network_ids` - Specifies the subnet Network ID.

* `description` - Specifies supplementary information about the load balancer.

* `admin_state_up` - The administrative state of the LoadBalancer.

* `ip_target_enable` - The value can be `true` (enabled) or `false` (disabled).

* `availability_zones` - Specifies the availability zones where the LoadBalancer will be located.

* `public_ip` - The elastic IP address of the instance.

  * `id` - Elastic IP ID.

  * `address` - Elastic IP address.

  * `ip_type` - Elastic IP type.

  * `bandwidth_name` - Bandwidth name.

  * `bandwidth_size` - Bandwidth size.

  * `bandwidth_charge_mode` - Bandwidth billing type.

  * `bandwidth_share_type` - Bandwidth sharing type.

* `created_at` - The time the LoadBalancer was created.

* `updated_at` - The time the LoadBalancer was last updated.

* `deletion_protection` - Specifies whether to enable deletion protection for the load balancer.
