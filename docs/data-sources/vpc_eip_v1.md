---
subcategory: "Virtual Private Cloud (VPC)"
---

# opentelekomcloud_vpc_eip_v1

Use this data source to get details about a specific VPC elastic IP.

## Example Usage

```hcl
data "opentelekomcloud_vpc_eip_v1" "eip_v1" {
  id = var.elastic_ip
}

output "eip_vpc_id" {
  value = data.opentelekomcloud_vpc_eip_v1.eip_v1.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
elastic IP in the current tenant. The given filters must match exactly one
elastic IP whose data will be exported as attributes.

* `id` - (Optional) Specifies a resource ID in UUID format.

* `status` - (Optional) The status of the specific elastic IP to retrieve.

* `public_ip_address` - (Optional) The public IP address of the elastic IP.

* `private_ip_address` - (Optional) The private IP address bound to the elastic IP.

* `port_id` - (Optional) The port ID.

-> `private_ip_address` and `port_id` are returned only when a port/private IP address is
associated with the elastic IP.

* `bandwidth_id` - (Optional) The bandwidth ID of specific elastic IP.

## Attributes Reference

All the argument attributes are also exported as result attributes.

* `ip_version` - The IP version of elastic IP.

* `bandwidth_share_type` - Specifies the EIP bandwidth type.

* `bandwidth_size` - Specifies the bandwidth (Mbit/s).

* `create_time` - Specifies the time (UTC) when the elastic IP is assigned.

* `tenant_id` - Specifies the project ID.

* `type` - Specifies the elastic IP type.
