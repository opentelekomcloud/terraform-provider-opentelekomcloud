---
subcategory: "APIG"
---

Up-to-date reference of API arguments for Anti-DDoS service you can get at
`https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/index.html`.

# opentelekomcloud_apigw_gateway_v2

API Gateway (APIG) is a high-performance, high-availability, and high-security API hosting service that helps you build,
manage, and deploy APIs at any scale.
With just a few clicks, you can integrate internal systems, and selectively expose capabilities with minimal costs and risks.

## Example Usage

```hcl
resource "opentelekomcloud_apigw_gateway_v2" "gateway" {
  name               = "test-gateway"
  spec_id            = "BASIC"
  vpc_id             = var.vpc_id
  subnet_id          = var.network_id
  security_group_id  = var.default_secgroup.id
  availability_zones = ["eu-de-01", "eu-de-02"]
  description        = "test gateway"
  bandwidth_size     = 5
  maintain_begin     = "22:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies gateway name.

* `spec_id` - (Required, ForceNew, String) Gateway edition. Options:
  This resource provides the following timeouts configuration options:
  - `BASIC`
  - `PROFESSIONAL`
  - `ENTERPRISE`
  - `PLATINUM`

* `vpc_id` - (Required, ForceNew, String) Specifies VPC ID.

* `subnet_id` - (Required, ForceNew, String) Specifies network ID.

* `security_group_id` - (Required, String) Specifies ID of the security group to which the gateway belongs.

* `description` - (Optional, String) Specifies gateway description.

* `availability_zones` - (Optional, List) Specifies gateway description.

* `bandwidth_size` - (Optional, String) Specifies outbound access bandwidth. This parameter is required if public outbound
  access is enabled for the gateway. After you configure the bandwidth for the gateway,
  users can access resources on public networks.

* `ingress_bandwidth_size` - (Optional, String) Specifies public inbound access bandwidth. This parameter is required if public
  inbound access is enabled for the gateway and loadbalancer_provider is set to elb.
  After you bind an EIP to the gateway, users can access APIs in the gateway from public networks using the EIP.

* `loadbalancer_provider` - (Optional, String) Specifies type of the load balancer used by the gateway.
  This resource provides the following timeouts configuration options:
    - `elb`

* `maintain_begin` - (Optional, String) Specifies start time of the maintenance time window.
  It must be in the format "xx:00:00". The value of xx can be 02, 06, 10, 14, 18, or 22.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `maintain_end` - End time of the maintenance time window. It must be in the format "xx:00:00".
  There is a 4-hour difference between the start time and end time.

* `vpc_ingress_address` - VPC ingress address.

* `public_egress_address` - IP address for public outbound access.

* `supported_features` - Supported features.

* `status` - Instance status.

* `project_id` - Instance project id.

* `region` - Instance region.

* `vpcep_service_name` - Name of a VPC endpoint service.

* `private_egress_addresses` - List of private egress addresses.
*
## Import

APIG Gateway can be imported using the `gateway_id`, e.g.

```shell
$ terraform import opentelekomcloud_apigw_gateway_v2.gateway c1881895-cdcb-4d23-96cb-032e6a3ee667
```

Note that the imported state may not be identical to your resource definition, due to `ingress_bandwidth_size` missing from the
API response. It is generally recommended running `terraform plan` after importing a gateway.

```
resource "opentelekomcloud_apigw_gateway_v2" "gateway" {
    ...

  lifecycle {
    ignore_changes = [
      ingress_bandwidth_size
    ]
  }
}
