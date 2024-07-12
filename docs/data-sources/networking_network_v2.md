---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_network_v2"
sidebar_current: "docs-opentelekomcloud-datasource-networking-network-v2"
description: |-
  Get the ID of an available network resource from OpenTelekomCloud
---

Up-to-date reference of API arguments for Network you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-cloud/api-ref/native_openstack_neutron_apis_v2.0/network/querying_networks.html#vpc-network-0001)

# opentelekomcloud_networking_network_v2

Use this data source to get the ID of an available OpenTelekomCloud network.

## Example Usage

```hcl
data "opentelekomcloud_networking_network_v2" "network" {
  name = "tf_test_network"
}
```

## Argument Reference

* `network_id` - (Optional) The ID of the network.

* `name` - (Optional) The name of the network.

* `matching_subnet_cidr` - (Optional) The CIDR of a subnet within the network.

* `tenant_id` - (Optional) The owner of the network.

## Attributes Reference

`id` is set to the ID of the found network. In addition, the following attributes are exported:

* `admin_state_up` - The administrative state of the network.

* `name` - See Argument Reference above.

* `shared` - Specifies whether the network resource can be accessed by any tenant or not.
