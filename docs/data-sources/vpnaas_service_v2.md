---
subcategory: "Virtual Private Network (VPN)"
---

# opentelekomcloud_vpnaas_service_v2

Use this data source to get details about a specific VPN.

## Example Usage

```hcl
variable "vpn_name" { }

data "opentelekomcloud_vpnaas_service_v2" "vpn" {
  name = "${var.vpn_name}"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPNs in the current region. The given filters must match exactly one VPN whose data will be exported as attributes.

* `id` - (Optional) The id of the specific VPN to retrieve.

* `status` - (Optional) The current status of the desired VPN. Can be either CREATING, OK, DOWN, PENDING_UPDATE, PENDING_DELETE, or ERROR.

* `name` - (Optional) A unique name for the VPN. The name must be unique for a tenant. The value is a string of no more than 64 characters and can contain digits, letters,
  underscores (_), and hyphens (-).


## Attributes Reference

The following attributes are exported:

* `id` - ID of the VPC.

* `name` -  See Argument Reference above.

* `status` - See Argument Reference above.
