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

* `description` - (Optional) Provides supplementary information about the load balancer.

* `router_id` - (Optional) The ID of the router (or VPC) this LoadBalancer belongs.

* `subnet_id` - (Optional) The ID of the subnet to which the LoadBalancer belongs.

* `l7_flavor` - (Optional) The ID of the Layer-7 flavor.

* `l4_flavor` - (Optional) The ID of the Layer-4 flavor.

* `vip_address` - (Optional) The IP address of the LoadBalancer.

* `vip_port_id` - (Optional) The Port ID of the Load Balancer IP.

* `provisioning_status` - (Optional) The provisioning status.

* `operating_status` - (Optional) The operating status.

* `guaranteed` - (Optional) Whether the load balancer is dedicated.

* `ipv6_vip_address` - (Optional) The private IPv6 address.

* `ipv6_vip_port_id` - (Optional) The port ID associated with the private IPv6 address.

* `ipv6_vip_subnet_id` - (Optional) The IPv6 subnet ID.

* `eip_id` - (Optional) The ID of an associated EIP.

* `public_ip_id` - (Optional) The ID of an associated public IP.

* `availability_zones` - (Optional) Availability zones of the load balancer.

* `l4_scale_flavor` - (Optional) The Layer-4 elastic flavor ID.

* `l7_scale_flavor` - (Optional) The Layer-7 elastic flavor ID.

* `member_device_id` - (Optional) The device ID of an associated backend server.

* `member_address` - (Optional) The private IP address of an associated backend server.

* `enterprise_project_id` - (Optional) The enterprise project ID.

* `ip_version` - (Optional) The IP version.

* `deletion_protection` - (Optional) Whether deletion protection is enabled.

* `elb_subnet_type` - (Optional) The IP version supported by the load balancer subnet.

* `protection_status` - (Optional) The modification protection status.

## Attributes Reference

In addition, the following attributes are exported:

* `network_ids` - Specifies the subnet Network ID.

* `admin_state_up` - The administrative state of the LoadBalancer.

* `ip_target_enable` - The value can be `true` (enabled) or `false` (disabled).

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

* `provider_name`, `project_id`, `provisioning_status`, and `operating_status` - Provider, project,
  provisioning, and operating metadata.

* `pools` and `listeners` - IDs of associated backend server groups and listeners.

* `ipv6_bandwidth_id` - The shared bandwidth ID used by the public IPv6 address.

* `eips` and `global_eips` - Associated EIP and global EIP details.

* `l4_scale_flavor` and `l7_scale_flavor` - Elastic flavor IDs.

* `frozen_scene`, `billing_info`, and `public_border_group` - Service and billing metadata.

* `waf_failure_action`, `charge_mode`, `protection_status`, and `protection_reason` - WAF,
  billing, and modification protection settings.

* `autoscaling` - Autoscaling status and minimum Layer-7 flavor.

* `loadbalancer_type`, `gateway_flavor_id`, `instance_type`, and `instance_id` - Load balancer
  type and instance metadata.

* `log_group_id` and `log_topic_id` - LTS logging destination.

* `custom_qos_limit` - Custom Layer-4 and Layer-7 connection and CPS limits.

* `service_lb_mode` - The service load balancer mode.

* `proxy_protocol_extensions` - Proxy protocol extension endpoint information.
