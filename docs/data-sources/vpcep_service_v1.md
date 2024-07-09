---
subcategory: "VPC Endpoint (VPCEP)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpcep_service_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpcep-service-v1"
description: |-
Get details about a specific VPCEP service from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPCEP service you can get at
[documentation portal](https://docs.otc.t-systems.com/vpc-endpoint/api-ref/apis/apis_for_managing_vpc_endpoint_services/querying_details_about_a_vpc_endpoint_service.html#vpcep-06-0202)

# opentelekomcloud_vpcep_service_v1

Use this data source to get details about a specific VPCEP service.

## Example Usage

```hcl
data "opentelekomcloud_vpcep_service_v1" "service" {
  name = var.service_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Specifies the unique ID of the VPC endpoint service.

* `name` - (Optional) Specifies the name of the VPC endpoint service.
  The value is not case-sensitive and supports fuzzy match.

* `status` - (Optional) Specifies the status of the VPC endpoint service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `port_id` - Specifies the ID for identifying the backend resource of the VPC endpoint service. The ID is in the form of the UUID.

* `vip_port_id` - Specifies the ID of the virtual NIC to which the virtual IP address is bound.
  This parameter is returned only when `port_id` is set to VIP.

* `server_type` - Specifies the resource type.

* `vpc_id` - Specifies the ID of the VPC to which the backend resource of the VPC endpoint service belongs.

* `approval_enabled` - Specifies whether connection approval is required.

* `service_type` - Specifies the type of the VPC endpoint service.

* `created_at` - Specifies the creation time of the VPC endpoint service.

* `updated_at` - Specifies the update time of the VPC endpoint service.

* `project_id` - Specifies the project ID.

* `ports` - Lists the port mappings opened to the VPC endpoint service.

* `tags` - Map of the resource tags.

* `connection_count` - Specifies the number of Creating or Accepted VPC endpoints under the VPC endpoint service.

* `tcp_proxy` - Specifies whether the client IP address and port number or marker_id information is transmitted to the server.

The `port` block supports:

* `client_port` - (Required) Specifies the port for accessing the VPC endpoint.

* `server_port` - (Required) Specifies the port for accessing the VPC endpoint service.
