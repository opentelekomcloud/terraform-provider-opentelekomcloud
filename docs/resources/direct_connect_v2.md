---
subcategory: "Direct Connect (DCaaS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_direct_connect_v2"
sidebar_current: "docs-opentelekomcloud-resource-direct-connect-v2"
description: |-
Manages a Direct Connect resource within OpenTelekomCloud.
---

# opentelekomcloud_direct_connect_v2 (Resource)

Up-to-date reference of API arguments for Direct Connect (DCaaS) you can get at
[documentation portal](https://docs.otc.t-systems.com/direct-connect/api-ref/apis/connection/creating_a_connection.html)

Example usage
-----------------
```hcl
resource "opentelekomcloud_direct_connect_v2" "direct_connect" {
  name           = "direct_connect"
  bandwidth      = 100
  location       = "location"
  provider_name  = "provider_name"
  port_type      = "port_type"
  admin_state_up = true
}
```


## Argument Reference

### Required

* `bandwidth` (Number) - Specifies the bandwidth of the connection in Mbit/s.
* `location` (String, ForceNew) - Specifies the connection access location.
* `provider_name` (String, ForceNew) - Specifies the carrier who provides the leased line.
* `port_type` (String, ForceNew) - Specifies the type of the port used by the connection. The value can be 1G, 10G, 40G, or 100G.

### Optional

* `admin_state_up` (Boolean, ForceNew)  - Specifies the administrative status of the connection. The value can be true or false.
* `charge_mode` (String, ForceNew) - Specifies the billing mode. The value can only be port for operations connections.
* `description` (String) - Provides supplementary information about the connection.
* `device_id` (String, ForceNew) - Specifies the gateway device ID of the connection.
* `hosting_id` (String, ForceNew) - Specifies the ID of the operations connection on which the hosted connection is created.
* `interface_name` (String, ForceNew) - Specifies the name of the interface accessed by the connection.
* `name` (String) - Specifies the connection name.
* `order_id` (String, ForceNew) - Specifies the connection order ID, which is used to support duration-based billing and identify user orders.
* `peer_location` (String, ForceNew) - Specifies the physical location of the peer device accessed by the connection, specific to the street or data center name.
* `product_id` (String, ForceNew) - Specifies the product ID corresponding to the connection's order, which is used to custom billing policies such as duration-based packages.
* `provider_status` (String) - Specifies the status of the carrier's leased line. The value can be ACTIVE or DOWN.
* `redundant_id` (String, ForceNew) - Specifies the ID of the redundant connection using the same gateway.
* `status` (String, ForceNew) - Specifies the connection status.
The value can be: `ACTIVE, DOWN, BUILD, ERROR, PENDING_DELETE, DELETED, APPLY, DENY, PENDING_PAY, PAID, ORDERING, ACCEPT, or REJECTED.`
* `tenant_id` (String, ForceNew) - Specifies the project ID.
* `type` (String, ForceNew) - Specifies the connection type. The value can only be `hosted`.
* `vlan` (Number, ForceNew) - Specifies the VLAN ID of the connection.

## Attributes Reference

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
