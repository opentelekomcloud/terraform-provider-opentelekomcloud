---
subcategory: "VPC Endpoint (VPCEP)"
---

Up-to-date reference of API arguments for VPCEP service you can get at
`https://docs.otc.t-systems.com/vpc-endpoint/api-ref/apis/apis_for_managing_vpc_endpoint_services`.

# opentelekomcloud_vpcep_service_v1

Manages a VPC Endpoint Service v1 resource within OpenTelekomCloud.

## Example Usage

```hcl
variable os_subnet_id {}
variable vpc_id {}
variable domain_id_1 {}
variable domain_id_2 {}

resource "opentelekomcloud_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = var.os_subnet_id
}

resource "opentelekomcloud_vpcep_service_v1" "service" {
  name        = "service_1"
  port_id     = opentelekomcloud_lb_loadbalancer_v2.lb_1.vip_port_id
  vpc_id      = var.vpc_id
  server_type = "LB"

  port {
    client_port = 80
    server_port = 8080
  }

  whitelist = [var.domain_id_1, var.domain_id_2]

  tags = {
    "key" : "value",
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specifies the name of the VPC endpoint service.
  The value contains a maximum of 16 characters, including letters, digits, underscores (_), and hyphens (-).
  * If you do not specify this parameter, the VPC endpoint service name is in the format: `regionName.serviceId`.
  * If you specify this parameter, the VPC endpoint service name is in the format: `regionName.serviceName.serviceId`.

* `port_id` - (Required) Specifies the ID for identifying the backend resource of the VPC endpoint service.
  The value is as follows:
  * If the backend service is an enhanced load balancer, the value is the ID of the port bound to the private IP address of the load balancer.
  * the backend resource is an ECS, the value is the NIC ID of the ECS where the VPC endpoint service is deployed.
  * the backend resource is a virtual IP address, the value is the NIC ID of the physical server where virtual resources are created.

* `pool_id` - (Optional) Specifies the ID of the cluster associated with the target VPCEP resource.

* `vip_port_id` - (Optional) Specifies the ID of the virtual NIC to which the virtual IP address is bound.

* `vpc_id` - (Optional) Specifies the ID of the VPC (router) to which the backend resource of the VPC endpoint service belongs.

* `approval_enabled` - (Optional) Specifies whether connection approval is required.

  * `false`: indicates that connection approval is not required.
    The created VPC endpoint is in the `accepted` state.
  * `true`: indicates that connection approval is required.
    The created VPC endpoint is in the `pendingAcceptance` state until the owner of the associated VPC endpoint
    service approves the connection.

  The default value is `true`.

* `service_type` - (Optional) Specifies the type of the VPC endpoint service.
  Only your private services can be configured into interface VPC endpoint services.

  There are two types of VPC endpoint services: `interface` and `gateway`.

  * `gateway`: VPC endpoint services of this type are configured by operations people.
    You can use them directly without the need to create one by yourselves.
  * `interface`: VPC endpoint services of this type include cloud services configured by operations people
    and private services created by yourselves. You cannot configure these cloud services, but can use them.

* `server_type` - (Required) Specifies the resource type.
  * `VM`: The backend resource is a server.
  * `VIP`: The backend resource is a virtual IP address that functions as a physical server hosting virtual resources.
  * `LB`: The backend resource is an enhanced load balancer.

* `port` - (Required) Lists the port mappings opened to the VPC endpoint service. See below for the details.

* `whitelist` - (Optional) Lists of domain IDs of target users.

* `tcp_proxy` - (Optional) Specifies whether the client IP address and port number or `marker_id` information is
  transmitted to the server.
  This parameter is available only when the server can parse fields tcp option and tcp payload.

  The values are as follows:

  * `close`: indicates that the TOA and Proxy Protocol methods are neither used.
  * `toa_open`: indicates that the TOA method is used.
  * `proxy_open`: indicates that the Proxy Protocol method is used.
  * `open`: indicates that the TOA and Proxy Protocol methods are both used.

  The default value is `close`.

* `tags` - (Optional) Map of the resource tags.

The `port` block supports:

* `client_port` - (Required) Specifies the port for accessing the VPC endpoint.

* `server_port` - (Required) Specifies the port for accessing the VPC endpoint service.

* `protocol` - (Required) Specifies the protocol used in port mappings. The value can be `TCP` or `UDP`.
  The default value is `TCP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of VPC endpoint service

## Import

VPC endpoint service can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpcep_service_v1.service 71ba78a2-d847-4882-8fd0-42c5854c1cbc
```
