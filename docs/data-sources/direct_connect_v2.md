---
subcategory: "Direct Connect (DCaaS)"
---
# opentelekomcloud_direct_connect_v2 (Data Source)

Use this data source to get details about a specific Direct Connect (DCaaS) connection.

Example usage
-----------------

```hcl
data "opentelekomcloud_direct_connect_v2" "direct_connect" {
  id = "direct_connect_id"
}
```


## Argument Reference

- `id` (String) - Specifies the direct connection ID.

## Attributes Reference
* `bandwidth` (Number) - Specifies the bandwidth of the connection in Mbit/s.
* `location` (String) - Specifies the connection access location.
* `provider_name` (String) - Specifies the carrier who provides the leased line.
* `port_type` (String) - Specifies the type of the port used by the connection. The value can be 1G, 10G, 40G, or 100G.
* `admin_state_up` (Boolean)  - Specifies the administrative status of the connection. The value can be true or false.
* `charge_mode` (String) - Specifies the billing mode. The value can only be port for operations connections.
* `description` (String) - Provides supplementary information about the connection.
* `device_id` (String) - Specifies the gateway device ID of the connection.
* `hosting_id` (String) - Specifies the ID of the operations connection on which the hosted connection is created.
* `interface_name` (String) - Specifies the name of the interface accessed by the connection.
* `name` (String) - Specifies the connection name.
* `order_id` (String) - Specifies the connection order ID, which is used to support duration-based billing and identify user orders.
* `peer_location` (String) - Specifies the physical location of the peer device accessed by the connection, specific to the street or data center name.
* `product_id` (String) - Specifies the product ID corresponding to the connection's order, which is used to custom billing policies such as duration-based packages.
* `provider_status` (String) - Specifies the status of the carrier's leased line. The value can be ACTIVE or DOWN.
* `redundant_id` (String) - Specifies the ID of the redundant connection using the same gateway.
* `status` (String) - Specifies the connection status.
  The value can be ACTIVE, DOWN, BUILD, ERROR, PENDING_DELETE, DELETED, APPLY, DENY, PENDING_PAY, PAID, ORDERING, ACCEPT, or REJECTED.
* `tenant_id` (String) - Specifies the project ID.
* `type` (String) - Specifies the connection type. The value can only be hosted.
* `vlan` (Number) - Specifies the VLAN ID of the connection.
* `applicant` (String) - This is a reserved field, which is not used currently.
* `apply_time` (String) - Specifies the time when the connection was requested.
* `building_line_product_id` (String) - This is a reserved field, which is not used currently.
* `cable_label` (String) - This is a reserved field, which is not used currently.
* `create_time` (String) - Specifies the time when the connection is created.
* `delete_time` (String) - Specifies the time when the connection was deleted.
* `email` (String) - This is a reserved field, which is not used currently.
* `id` (String) - Specifies the connection ID.
* `lag_id` (String) - This is a reserved field, which is not used currently.
* `last_onestop_product_id` (String) - This is a reserved field, which is not used currently.
* `mobile` (String) - This is a reserved field, which is not used currently.
* `onestop_product_id` (String) - This is a reserved field, which is not used currently.
* `peer_port_type` (String) - This is a reserved field, which is not used currently.
* `peer_provider` (String) - This is a reserved field, which is not used currently.
* `period_num` (Number) - This is a reserved field, which is not used currently.
* `period_type` (Number) - This is a reserved field, which is not used currently.
* `reason` (String) - This is a reserved field, which is not used currently.
* `region_id` (String) - Specifies the region ID.
* `service_key` (String) - This is a reserved field, which is not used currently.
* `spec_code` (String) - This is a reserved field, which is not used currently.
* `vgw_type` (String) - Specifies the type of the gateway. Currently, only the default type is supported.
