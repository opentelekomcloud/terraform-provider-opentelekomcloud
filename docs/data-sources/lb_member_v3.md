---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_member_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lb-member-v3"
description: |-
  Gets details about an ELB v3 backend server within OpenTelekomCloud.
---

# opentelekomcloud_lb_member_v3

Gets a single member of an ELB v3 backend server group. See the
[API reference](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/backend_server/querying_backend_servers_in_a_backend_server_group.html).
The query must match exactly one member; no matches or multiple matches produce an error.

## Example Usage

```hcl
data "opentelekomcloud_lb_member_v3" "member" {
  pool_id = var.pool_id
  name    = "application-server"
}
```

## Argument Reference

* `pool_id` - (Required, String) Backend server group ID.

* `member_id` - (Optional, String) Member ID. When configured, fetches this member directly and ignores other filters.

* `name` - (Optional, String) Member name.

* `address` - (Optional, String) Backend IP address.

* `subnet_id` - (Optional, String) Subnet CIDR ID, sent as the API's `subnet_cidr_id` filter.

* `protocol_port` - (Optional, Int) Backend protocol port.

* `weight` - (Optional, Int) Member weight. An explicit `0` is supported.

* `admin_state_up` - (Optional, Bool) Administrative state. An explicit `false` is supported.

* `operating_status` - (Optional, String) Member health status.

* `ip_version` - (Optional, String) IP version.

* `member_type` - (Optional, String) Backend server type.

* `instance_id` - (Optional, String) Associated ECS instance ID.

* `availability_zone` - (Optional, String) Backend server AZ.

## Attributes Reference

All filter attributes are also populated from the matching member. The following additional attributes are exported:

* `id` - Stable composite identifier in the form `pool_id/member_id`.

* `project_id` - Project ID.

* `created_at` - Creation time, if returned by the API.

* `updated_at` - Last update time, if returned by the API.

* `status` - Per-listener health status list. Each entry contains:
  * `listener_id` - Listener ID.
  * `operating_status` - Member health status for this listener.
  * `created_at` - Status creation time.
  * `updated_at` - Status update time.
  * `reason` - Optional health-check failure details:
    * `reason_code` - Failure code.
    * `expected_response` - Expected health-check response code.
    * `healthcheck_response` - Actual health-check response code.
