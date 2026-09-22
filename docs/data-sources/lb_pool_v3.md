---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_pool_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lb-pool-v3"
description: |-
  Gets details about an ELB v3 backend server group from OpenTelekomCloud.
---

Up-to-date reference of API arguments for DLB backend server groups is available in the
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/backend_server_group).

# opentelekomcloud_lb_pool_v3

Use this data source to get details about an existing ELB v3 backend server group.

## Example Usage

```hcl
data "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = var.loadbalancer_id
  name            = "application-pool"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Specifies the backend server group ID. When specified, the data source retrieves the group
  directly and ignores the other filters.

* `name` - (Optional) Specifies the backend server group name.

* `description` - (Optional) Provides supplementary information about the backend server group.

* `healthmonitor_id` - (Optional) Specifies the ID of the associated health monitor.

* `lb_algorithm` - (Optional) Specifies the load balancing algorithm. Valid values are `ROUND_ROBIN`,
  `LEAST_CONNECTIONS`, `SOURCE_IP`, and `QUIC_CID`.

* `protocol` - (Optional) Specifies the backend server group protocol. Valid values are `TCP`, `UDP`, `HTTP`,
  `HTTPS`, and `QUIC`.

* `loadbalancer_id` - (Optional) Specifies the ID of an associated load balancer.

* `listener_id` - (Optional) Specifies the ID of an associated listener.

* `ip_version` - (Optional) Specifies the supported IP address version.

* `member_address` - (Optional) Specifies the private IP address of a backend server. This field is used only as a
  query filter and is not returned by the API.

* `member_device_id` - (Optional) Specifies the cloud server ID of a backend server. This field is used only as a
  query filter and is not returned by the API.

* `member_instance_id` - (Optional) Specifies the backend server instance ID. This field is used only as a query
  filter and is not returned by the API.

* `member_deletion_protection` - (Optional) Specifies whether removal protection is enabled for backend servers.

* `vpc_id` - (Optional) Specifies the ID of the VPC where the backend server group operates.

* `type` - (Optional) Specifies the backend server group type.

* `protection_status` - (Optional) Specifies the modification protection status. Valid values are `nonProtection`
  and `consoleProtection`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `project_id` - Specifies the project ID of the backend server group.

* `session_persistence` - Specifies the sticky session configuration.
  * `type` - Specifies the sticky session type.
  * `cookie_name` - Specifies the sticky session cookie name.
  * `persistence_timeout` - Specifies the stickiness duration, in minutes.

* `slow_start` - Specifies the slow start configuration.
  * `enable` - Specifies whether slow start is enabled.
  * `duration` - Specifies the slow start duration, in seconds.

* `listener_ids` - Lists the IDs of listeners associated with the backend server group.

* `loadbalancer_ids` - Lists the IDs of load balancers associated with the backend server group.

* `member_ids` - Lists the IDs of backend servers in the backend server group.

* `protection_reason` - Specifies why modification protection is enabled.

* `created_at` - Specifies the time when the backend server group was created.

* `updated_at` - Specifies the time when the backend server group was last updated.
