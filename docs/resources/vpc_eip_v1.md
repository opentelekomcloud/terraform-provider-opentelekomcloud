---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpc_eip_v1"
sidebar_current: "docs-opentelekomcloud-resource-vpc-eip-v1"
description: |-
  Manages a VPC EIP resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPC eip you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/apis/elastic_ip)

# opentelekomcloud_vpc_eip_v1

Manages a V1 EIP resource within OpenTelekomCloud VPC.

## Example Usage

```hcl
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
```

## EIP with name

```hcl
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
    name = "my_eip"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
```

## EIP with shared bandwidth

```hcl
resource "opentelekomcloud_vpc_bandwidth_v2" "shared" {
  name = "shared-bandwidth"
  size = 50
}

resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    id         = opentelekomcloud_vpc_bandwidth_v2.shared.id
    share_type = "WHOLE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the VPC v1 client.
  If omitted, the `region` argument of the provider is used. Changing this
  creates a new service.

* `publicip` - (Required) The elastic IP address object.

* `bandwidth` - (Required) The bandwidth object.

The `publicip` block supports:

* `type` - (Required) The value must be a type supported by [the system](https://docs.otc.t-systems.com/api/eip/eip_api_0001.html#eip_api_0001__en-us_topic_0201534274_table4491214).
  The value can be `5_bgp`, `5_mailbgp` and `5_gray`. Changing this creates a new eip.

* `ip_address` - (Optional) The value must be a valid IP address in the available
  IP address segment. Changing this creates a new eip.

* `ip_version` - (Optional) Specifies the EIP version. Valid values are `4` and `6`.
  Changing this creates a new EIP. IPv6 EIPs are not currently supported.

* `port_id` - (Optional) The port id which this eip will associate with. If the value
  is `""` or this not specified, the eip will be in unbind state.

* `name` - (Required) The ip name, which is a string of 1 to 64 characters.

The `bandwidth` block supports:

* `id` - (Optional) An existing shared bandwidth ID to attach the EIP to. This
  conflicts with `name`.

* `name` - (Optional) The bandwidth name, which is a string of 1 to 64 characters
  that contain letters, digits, underscores (_), and hyphens (-). Required with
  `size` when creating a new bandwidth.

* `size` - (Optional) The bandwidth size. The value ranges from 1 to 300 Mbit/s.
  Required with `name` when creating a new bandwidth.

* `share_type` - (Required) Whether the bandwidth is shared or exclusive. Changing
  this creates a new eip.

* `charge_mode` - (Optional) This is a reserved field. If the system supports charging
  by traffic and this field is specified, then you are charged by traffic for elastic
  IP addresses. Changing this creates a new eip.

* `tags` - (Optional) Tags key/value pairs to associate with the eip.

* `enterprise_project_id` - (Optional) Specifies the enterprise project associated
  with the EIP. If omitted, the provider-level enterprise project or `0` is used.
  Changing this creates a new EIP.

* `value_specs` - (Optional) Legacy additional creation parameters. This compatibility
  argument uses the legacy creation request when set.

* `unbind_port` - (Optional) The value `true` indicates that port will be unassigned from EIP.
  This parameter work only with already allocated resource.

## Attributes Reference

The following attributes are exported:

* `id` - The VPC EIP id.

* `region` - See Argument Reference above.

* `publicip/type` - See Argument Reference above.

* `publicip/ip_address` - See Argument Reference above.

* `publicip/ip_version` - See Argument Reference above.

* `publicip/port_id` - See Argument Reference above.

* `publicip/name` - See Argument Reference above.

* `bandwidth/id` - See Argument Reference above.

* `bandwidth/name` - See Argument Reference above.

* `bandwidth/size` - See Argument Reference above.

* `bandwidth/share_type` - See Argument Reference above.

* `bandwidth/charge_mode` - See Argument Reference above.

* `tags` - See Argument Reference above.

* `enterprise_project_id` - The enterprise project associated with the EIP.

* `public_border_group` - The EIP location, such as `center` or an edge site.

* `allow_share_bandwidth_types` - Shared-bandwidth types to which the EIP can be added.

## Import

EIPs can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpc_eip_v1.eip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
